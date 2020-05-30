package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Set the default values on the config from the tags. This is one of the tag handlers
// specifically implemented to learn the reflect library. Note that a Config pointer
// needs to be passed to this function so that the config is editable and f.CanSet()
// returns true (e.g. it's addressable) and the value can be parsed and set.
func defaults(c *Config) error {
	// A tagHandlerFunc for handling the default tag t on the field ft with value f.
	defaultHandler := func(t string, f reflect.Value, ft reflect.StructField) error {
		if f.CanSet() {
			// Parse the value from the tag and set it on the struct
			dv, err := parse(ft.Type, t)
			if err != nil {
				return fmt.Errorf("cannot set %s (%s) to default %q: %s", ft.Name, ft.Type.Name(), t, err)
			}
			f.Set(reflect.ValueOf(dv))
			return nil
		}
		return fmt.Errorf("cannot set %s (%s) to default %q: field cannot be set", ft.Name, ft.Type.Name(), t)
	}

	// Execute the handle tags recursive function with the default handler for default tags.
	return handleTags(reflect.ValueOf(c), "default", defaultHandler)
}

// Set the environment variable values on the config from the tags. Note that a Config
// pointer needs to be passed so that the config is editable (see defaults and
// handleTags for more on this).
func environs(c *Config) error {
	// A tagHandlerFunc for handling the env tag t on the field ft with value f.
	environHandler := func(t string, f reflect.Value, ft reflect.StructField) error {
		// Get the value from the environment, if empty then continue
		ev := os.Getenv(t)
		if ev == "" {
			return nil
		}

		// Check if the field can be set without a panic
		if f.CanSet() {
			dv, err := parse(ft.Type, ev)
			if err != nil {
				return fmt.Errorf("cannot set %s (%s) to env var %q: %s", ft.Name, ft.Type.Name(), ev, err)
			}
			f.Set(reflect.ValueOf(dv))
			return nil
		}
		return fmt.Errorf("cannot set %s (%s) to env var %q: field cannot be set", ft.Name, ft.Type.Name(), ev)
	}

	// Execute the handle tags recursive function with environ handler for env tags
	return handleTags(reflect.ValueOf(c), "env", environHandler)
}

// A function to handle the value of a tag t on a field f of the specified type, ft.
// This type is used in defaults and environs to set the config from fields.
type tagHandlerFunc func(t string, f reflect.Value, ft reflect.StructField) error

// handleTags walks the specified struct as a reflect.Value passing any fields with the
// tag name to the handler function. Additionally handleTags recursively dives into
// nested structs so that all tags inline with a specific struct are handled. This
// function was to dive into the reflect library - which is an interesting kind of
// beast. Note that if you want to be able to modify the struct being passed in, then
// you have to pass a value that is a pointer, otherwise the struct will be immutable.
// This function attempts to avoid panics where possible and return errors so that any
// panics can be raised at a higher level, however it is important to note that many of
// the reflect calls do lead to panics and if so, checks need to be added to error
// instead of panic.
func handleTags(v reflect.Value, name string, handle tagHandlerFunc) error {
	// Extract the element if v is an interface or is nil
	if v.Kind() == reflect.Interface && !v.IsNil() {
		elm := v.Elem()
		if v.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			v = elm
		}
	}

	// If v is a pointer (for editing), get the value the pointer directs to.
	// Note that an interface pointer will not be able to be unwrapped, which is why
	// this function accepts a reflect.Value rather than an interface{}.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get the type of the object for field level processing.
	t := v.Type()

	// Iterate over the fields of the struct
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)

		// If field is a struct then recursively call defaults
		if ft.Type.Kind() == reflect.Struct {
			handleTags(f, name, handle)
			continue
		}

		// Get the tag by name and handle it if it's available, returning any errors.
		tag := ft.Tag.Get(name)
		if tag != "" {
			if err := handle(tag, f, ft); err != nil {
				return err
			}
		}
	}

	return nil
}

// Time types for comparison
var (
	timeType     = reflect.TypeOf(time.Time{})
	durationType = reflect.TypeOf(time.Minute)
)

// Parse a value into a specific type from the value string. Used to parse environment
// variables (which are all strings) as well as defaults inside of tags. This function
// is supposed to be a general purpose function, callers can choose to panic or handle
// the error depending on what kind of parsing is happening.
func parse(t reflect.Type, v string) (interface{}, error) {
	// Handle complex types first
	switch t {
	case timeType:
		return time.Parse(time.RFC3339, v)
	case durationType:
		return time.ParseDuration(v)
	}

	// Handle simpler types
	switch t.Kind() {
	case reflect.Bool:
		return strconv.ParseBool(v)
	case reflect.String:
		return v, nil
	case reflect.Int:
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	case reflect.Int8:
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	case reflect.Int16:
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	case reflect.Int32:
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	case reflect.Int64:
		i, err := strconv.ParseInt(v, 10, 64)
		return int64(i), err
	case reflect.Uint:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	case reflect.Uint8:
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	case reflect.Uint16:
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	case reflect.Uint32:
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	case reflect.Uint64:
		i, err := strconv.ParseUint(v, 10, 64)
		return uint64(i), err
	case reflect.Float64:
		return strconv.ParseFloat(v, 64)
	case reflect.Float32:
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	default:
		return nil, fmt.Errorf("could not parse %q: unknown kind %v", v, t.Kind())
	}
}
