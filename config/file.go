package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var exts = [3]string{".json", ".yaml", ".yml"}

// Find searches for configuration files named catena.[json,yaml] in the current
// directory, then in the user config dir named ~/.config/catena/config.[json,yaml] and
// finally in /etc/catena/config.[json,yaml] -- it returns all paths that it finds
// without error but does not guarantee that these configuration files are readable.
// Note that the paths are returned in priority order, so most callers will load the
// files in reverse order.
func Find() (paths []string) {
	paths = make([]string, 0)

	// Look in the current working directory first
	for _, ext := range exts {
		if _, err := os.Stat("catena" + ext); !os.IsNotExist(err) {
			paths = append(paths, "catena"+ext)
		}
	}

	// Look in the user configuration directory next
	cdir, err := os.UserConfigDir()
	if err == nil {
		for _, ext := range exts {
			path := filepath.Join(cdir, "catena", "config"+ext)
			if _, err = os.Stat(path); !os.IsNotExist(err) {
				paths = append(paths, path)
			}
		}
	}

	// Finally look in system configuration locations
	for _, ext := range exts {
		path := filepath.Join("etc", "catena", "config"+ext)
		if _, err = os.Stat(path); !os.IsNotExist(err) {
			paths = append(paths, path)
		}
	}

	return paths
}

// LoadSystem loads the system configurations discovered using Find() in reverse order,
// maintaining the Find() priority. E.g. it first loads from /etc/catena/config.yaml
// then from the user configuration, then the local directory, etc. Like LoadFile() this
// method does not modify the original and instead makes a copy.
func (c Config) LoadSystem() (d Config, err error) {
	// Make a copy of the config
	// NOTE: this expects a shallow copy with no pointer references e.g. maps or slices
	d = c

	paths := Find()
	for i := len(paths) - 1; i >= 0; i-- {
		if d, err = d.LoadFile(paths[i]); err != nil {
			return d, fmt.Errorf("could not load %s: %s", paths[i], err)
		}
	}

	return d, nil
}

// LoadFile updates the configuration from the specified file without modifying
// the original configuration (makes a copy).
func (c Config) LoadFile(path string) (d Config, err error) {
	// Make a copy of the config
	// NOTE: this expects a shallow copy with no pointer references e.g. maps or slices
	d = c

	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return d, err
	}

	switch filepath.Ext(path) {
	case ".json":
		if err = json.Unmarshal(data, &d); err != nil {
			return d, err
		}
	case ".yml", ".yaml":
		if err = yaml.Unmarshal(data, &d); err != nil {
			return d, err
		}
	default:
		return d, fmt.Errorf("unknown file extension %q", filepath.Ext(path))
	}

	return d, nil
}

// DumpFile writes the configuration into the specified file using the extension to
// determine the serialization format.
func (c Config) DumpFile(path string) (err error) {
	var data []byte
	switch filepath.Ext(path) {
	case ".json":
		if data, err = json.MarshalIndent(c, "", "  "); err != nil {
			return err
		}
	case ".yml", ".yaml":
		if data, err = yaml.Marshal(c); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown file extension %q", filepath.Ext(path))
	}

	if err = ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}
