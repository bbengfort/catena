/*
Package migrations manages the state of the Catena database. SQL files should be added
to this directory that implement DDL commands that update the schema of the database
that Catena is connected to. The SQL files are then compiled into the binary using the
catena makemigrations command along with go generate. The catena server and command can
compare the state of the database with its expected state and run any migrations that
are required.
*/
package migrations

//go:generate go run ../cmd/makemigrations

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// weak sauce helper for debugging
var debug = false

// contains all of the available migrations generated by go generate and added in the
// init function of the file migrations.go (the file generated by go generate).
var migrations []Migration

// External API

// Migrate the database to the specified revision, if the revision is negative,
// then apply all unapplied migrations to the database. If the revision is less than
// the current revision, then the database is rolled back to that state. This function
// cannot drop the migrations table, use the Delete() function to completely rollback
// all migrations and delete the migrations table. Returns the total number of
// migrations that were executed against the database.
func Migrate(r int64, conn *sql.DB) (n int, err error) {
	var tx *sql.Tx
	if tx, err = conn.Begin(); err != nil {
		return 0, fmt.Errorf("could not begin migration transaction: %s", err)
	}
	defer tx.Rollback()

	// Refresh the database and also apply the 0 migration (initialization).
	if err = refreshTx(tx); err != nil {
		return 0, err
	}

	// There are no migrations except the initialization migration.
	if len(migrations) == 1 {
		return 0, nil
	}

	// Migration 0 has already been applied, so skip it and apply all others.
	for _, m := range migrations[1:] {
		if r < 0 || m.Revision <= r {
			if !m.Active {
				if err = m.upTx(tx); err != nil {
					return 0, fmt.Errorf("could not apply revision %d: %s", m.Revision, err)
				}
				n++
			}
		} else if r >= 0 && m.Revision > r {
			if m.Active {
				if err = m.downTx(tx); err != nil {
					return 0, fmt.Errorf("could not rollback revision %d: %s", m.Revision, err)
				}
				n++
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("could not commit %d migrations: %s", n, err)
	}

	return n, nil
}

// Num returns the number of migrations in the package
func Num() int {
	return len(migrations)
}

// Revision returns the migration for the specified revision.
func Revision(r int64, conn *sql.DB) (Migration, error) {
	if conn != nil {
		if err := Refresh(conn); err != nil {
			return Migration{}, err
		}
	}

	i := sort.Search(len(migrations), func(i int) bool {
		return r <= migrations[i].Revision
	})

	if i < len(migrations) && migrations[i].Revision == r {
		return migrations[i], nil
	}
	return Migration{}, fmt.Errorf("no migration found for revision %d", r)
}

// Current returns the most recently applied revision.
func Current(conn *sql.DB) (Migration, error) {
	if conn != nil {
		if err := Refresh(conn); err != nil {
			return Migration{}, err
		}
	}

	if !migrations[0].dbsync {
		return Migration{}, errors.New("migrations have not been synchronized with the database")
	}

	var current Migration
	for _, m := range migrations {
		if !m.Active {
			break
		}
		current = m
	}

	return current, nil
}

// Refresh the state of the migrations from the database.
func Refresh(conn *sql.DB) (err error) {
	var tx *sql.Tx
	if tx, err = conn.Begin(); err != nil {
		return fmt.Errorf("could not begin refresh transaction: %s", err)
	}
	defer tx.Rollback()

	if err = refreshTx(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func refreshTx(tx *sql.Tx) (err error) {
	// Apply migration 0 which initializes the migration schema
	if err = migrations[0].upTx(tx); err != nil {
		return err
	}

	var rows *sql.Rows
	if rows, err = tx.Query("SELECT revision, name, active, applied, created FROM migrations ORDER BY revision"); err != nil {
		return fmt.Errorf("could not fetch migrations: %s", err)
	}
	defer rows.Close()

	var i int
	for rows.Next() {
		mr := new(Migration)
		if err = rows.Scan(&mr.Revision, &mr.Name, &mr.Active, &mr.Applied, &mr.Created); err != nil {
			return fmt.Errorf("could not scan migration: %s", err)
		}

		if i < len(migrations) {
			// update the migration at position i with information from the database
			if migrations[i].Revision == mr.Revision {
				migrations[i].Active = mr.Active
				migrations[i].Applied = mr.Applied
				migrations[i].Created = mr.Created
				migrations[i].dbsync = true
			} else {
				return fmt.Errorf("unknown revision %d %q in database does not match local revision %d %q", mr.Revision, mr.Name, migrations[i].Revision, migrations[i].Name)
			}

		} else {
			// the database has a migration we're unaware of, which is bad
			return fmt.Errorf("unknown revision %d %q stored in the database but not locally", mr.Revision, mr.Name)
		}

		i++ // keep track of which migrations we have synchronized
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error while reading migrations: %s", err)
	}

	// Check if we need to insert migrations into the database
	if i < len(migrations) {
		var stmt *sql.Stmt
		if stmt, err = tx.Prepare("INSERT INTO migrations (revision, name, created) VALUES ($1, $2, $3)"); err != nil {
			return fmt.Errorf("could not prepare migrations insert statement: %s", err)
		}
		for j := i; j < len(migrations); j++ {
			migrations[j].Created = time.Now().UTC()
			if _, err = stmt.Exec(migrations[j].Revision, migrations[j].Name, migrations[j].Created); err != nil {
				return fmt.Errorf("could not insert revision %d %q", migrations[j].Revision, migrations[j].Name)
			}
			migrations[j].dbsync = true
		}
	}

	return nil
}

var sqltempl = template.Must(template.New("").Parse(`-- Revision {{ .Revision }} generated on {{ .Timestamp }}
-- migrate: up
-- insert up migration sql here

-- migrate: down
-- insert down migration sql here
`))

// New creates a new migration file from a template for the next revision and checks to
// make sure that it is valid. Specify the migrations directory for verification.
func New(name, dir string) (path string, err error) {
	if name == "" {
		name = fmt.Sprintf("auto_%s", time.Now().Format("200601021504"))
	}

	r := migrations[len(migrations)-1].Revision + 1

	var matches []string
	if matches, err = filepath.Glob(filepath.Join(dir, fmt.Sprintf("%04d_*.sql", r))); err != nil {
		return "", fmt.Errorf("could not check for duplicate revisions: %s", err)
	} else if len(matches) > 0 {
		return "", fmt.Errorf("a migration with revision %d already exists: %s (did you run go generate?)", r, matches[0])
	}

	name = strings.Replace(name, " ", "_", -1)
	path = filepath.Join(dir, fmt.Sprintf("%04d_%s.sql", r, name))

	builder := &bytes.Buffer{}
	ctx := struct {
		Revision  int64
		Timestamp string
	}{Revision: r, Timestamp: time.Now().Format("2006-01-02 15:04")}
	if err = sqltempl.Execute(builder, ctx); err != nil {
		return "", fmt.Errorf("could not execute sql migration template: %s", err)
	}

	var f *os.File
	if f, err = os.Create(path); err != nil {
		return "", fmt.Errorf("could not create %s: %s", path, err)
	}
	defer f.Close()

	if _, err = f.Write(builder.Bytes()); err != nil {
		return "", fmt.Errorf("could not write sql migration template: %s", err)
	}

	return path, nil
}

// Migration combines the information about the state of the database and how it has
// been migrated from the migrations table alongside the migration code stored in SQL
// files and compiled into the binary using go generate.
type Migration struct {
	Revision int64     // the unique id of the migration, prefix from the migration file
	Name     string    // the human readable name of the migration, suffix of migration file
	Active   bool      // if the migration has been applied or not
	Applied  time.Time // the timestamp the migration was applied
	Created  time.Time // the timestamp the migration was created in the database
	filename string    // the filename of the associated migration file
	up       string    // the sql query to apply the migration (read from -- migrate: up)
	down     string    // the sql query to rollback the migration (read from -- migrate: down)
	dbsync   bool      // if the migration has been synchronized to the database
}

// Up applies the migration to the database.
func (m *Migration) Up(conn *sql.DB) (err error) {
	var tx *sql.Tx
	if tx, err = conn.Begin(); err != nil {
		return fmt.Errorf("could not begin transaction to apply revision %d: %s", m.Revision, err)
	}
	defer tx.Rollback()

	if err = m.upTx(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Migration) upTx(tx *sql.Tx) (err error) {
	if _, err = tx.Exec(m.up); err != nil {
		return fmt.Errorf("could not exec apply revision %d: %s", m.Revision, err)
	}

	// If this is migration 0, we have a special sql query so we don't keep updating the applied timestamp
	sql := "UPDATE migrations SET active=$1, applied=$2 WHERE revision=$3"
	if m.Revision == 0 {
		sql = "UPDATE migrations SET active=$1, applied=$2 WHERE revision=$3 AND active='f' AND applied is NULL"
	}

	if _, err = tx.Exec(sql, true, time.Now().UTC(), m.Revision); err != nil {
		return fmt.Errorf("could not update migration status: %s", err)
	}

	return nil
}

// Down rolls back the migration from the database.
func (m *Migration) Down(conn *sql.DB) (err error) {
	var tx *sql.Tx
	if tx, err = conn.Begin(); err != nil {
		return fmt.Errorf("could not begin transaction to rollback revision %d: %s", m.Revision, err)
	}
	defer tx.Rollback()

	if err = m.downTx(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func (m *Migration) downTx(tx *sql.Tx) (err error) {
	if _, err = tx.Exec(m.down); err != nil {
		return fmt.Errorf("could not exec rollback revision %d: %s", m.Revision, err)
	}

	if _, err = tx.Exec("UPDATE migrations SET active=$1, applied=NULL WHERE revision=$3", false, m.Revision); err != nil {
		return fmt.Errorf("could not update migration status: %s", err)
	}

	return nil
}

func (m *Migration) String() string {
	builder := &strings.Builder{}
	fmt.Fprintf(builder, "revision: %d\nname: %q\n", m.Revision, m.Name)
	if m.dbsync {
		fmt.Fprintf(
			builder,
			"active: %t\napplied: %s\ncreated: %s\n",
			m.Active,
			m.Applied.Format("Jan 02, 2006 at 15:04:05 MST"),
			m.Created.Format("Jan 02, 2006 at 15:04:05 MST"),
		)
	}

	fmt.Fprintf(builder, "filename: %s\n", m.filename)
	fmt.Fprintf(builder, "predecessors: %d\n", m.Predecessors())
	fmt.Fprintf(builder, "successors: %d\n", m.Successors())

	if debug {
		if m.up != "" {
			fmt.Fprint(builder, "\nup\n--\n")
			fmt.Fprintln(builder, m.up)
		}

		if m.down != "" {
			fmt.Fprint(builder, "\ndown\n----\n")
			fmt.Fprintln(builder, m.down)
		}
	}

	return builder.String()
}

// Filename returns the original filename (before generation)
func (m *Migration) Filename() string {
	return m.filename
}

// UpSQL returns the query that will be executed when Up() is run
func (m *Migration) UpSQL() string {
	return m.up
}

// DownSQL returns the query that will be executed when Down() is run
func (m *Migration) DownSQL() string {
	return m.down
}

// DBSync returns true if the migration is in the database.
func (m *Migration) DBSync() bool {
	return m.dbsync
}

// Predecessors returns the number of migrations before this migration.
func (m *Migration) Predecessors() (n int) {
	for _, o := range migrations {
		if m.Revision == o.Revision {
			break
		}
		n++
	}
	return n
}

// Successors returns the number of migrations after this migration.
func (m *Migration) Successors() (n int) {
	i := sort.Search(len(migrations), func(i int) bool {
		return m.Revision <= migrations[i].Revision
	})

	if i < len(migrations) && migrations[i].Revision == m.Revision {
		if i+1 == len(migrations) {
			return 0
		}
		return len(migrations[i+1:])
	}
	panic(fmt.Errorf("revisions %d was not in the package migrations", m.Revision))
}

// Parse a migration file into an unsynchronized migration struct. This function is only
// used by go generate and though it can help users diagnose migration parsing issues,
// is generally not useful outside of the package.
func Parse(filename string) (m *Migration, err error) {
	m = &Migration{
		filename: filename,
		dbsync:   false,
	}

	if !strings.HasSuffix(filename, ".sql") {
		return nil, errors.New("migration filenames must end in .sql extension")
	}

	parts := strings.Split(strings.TrimSuffix(filename, ".sql"), "_")
	if len(parts) < 2 {
		return nil, errors.New("must format migration filenames as XXXX_description.sql")
	}

	if m.Revision, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return nil, fmt.Errorf("could not parse revision from %q: %s", filename, err)
	}

	m.Name = strings.Join(parts[1:], " ")

	upb := &strings.Builder{}
	dnb := &strings.Builder{}

	var (
		f       *os.File
		current *strings.Builder
	)

	if f, err = os.Open(filename); err != nil {
		return nil, fmt.Errorf("could not open %q: %s", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "--") {
			line := strings.ToLower(strings.TrimLeft(line, "- \t"))
			if strings.HasPrefix(line, "migrate:") {
				line := strings.TrimLeft(line, "migrate: \t")
				if line == "up" {
					current = upb
				} else if line == "down" {
					current = dnb
				} else {
					return nil, fmt.Errorf("%q is not a valid migrate directive", line)
				}
			}
			continue
		}

		if current == nil {
			return nil, errors.New("did not encounter a 'migrate:' directive")
		}

		current.WriteString(line + " ")
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	m.up = upb.String()
	m.down = dnb.String()

	// append the migration to the migrations
	return m, nil
}

// ByRevision implements sort.Interface for []Migration based on the Revision field.
type ByRevision []Migration

func (a ByRevision) Len() int           { return len(a) }
func (a ByRevision) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRevision) Less(i, j int) bool { return a[i].Revision < a[j].Revision }

// Internal API

// Add a migration to the local migrations slice, panic if things go wrong.
func local(revision int64, name, filename string, up []byte, down []byte) {
	m := &Migration{
		Revision: revision,
		Name:     name,
		filename: filename,
		dbsync:   false,
	}

	if len(up) > 0 {
		m.up = string(up)
	}

	if len(down) > 0 {
		m.down = string(down)
	}

	if len(migrations) > 0 {
		prev := migrations[len(migrations)-1]
		if m.Revision <= prev.Revision {
			panic(fmt.Errorf("cannot insert revision %d after revision %d", m.Revision, prev.Revision))
		}
	}

	migrations = append(migrations, *m)
}
