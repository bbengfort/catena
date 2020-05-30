package config_test

import (
	"bytes"
	"crypto/sha256"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/bbengfort/catena/config"
	"github.com/stretchr/testify/require"
)

func TestConfigFiles(t *testing.T) {
	// Create configuration with default values
	conf, err := New()
	require.NoError(t, err)

	// Should not be able to load a file that doesn't exist (doesn't test bad extension)
	_, err = conf.LoadFile("config.foo")
	require.Error(t, err)

	// Should not be able to dump a file with bad extension
	require.Error(t, conf.DumpFile("config.foo"))

	for _, ext := range [3]string{".json", ".yaml", ".yml"} {

		// Check if the test file exists and generate if required
		testpath := filepath.Join("testdata", "config"+ext)
		if _, err = os.Stat(testpath); os.IsNotExist(err) {
			if os.Getenv("CATENA_TEST_FIXTURES") == "" {
				t.Skipf("test fixture %q does not exist, use $CATENA_TEST_FIXTURES to generate", testpath)
			}

			// Generate fixtures if anything is set in the $CATENA_TEST_FIXTURES
			require.NoError(t, conf.DumpFile(testpath))
		}

		// Load the configuration, ensuring that it is not equal to the original (made a copy)
		loaded, err := conf.LoadFile(testpath)
		require.NoError(t, err)
		require.NotEqual(t, conf, loaded)

		// Dump the configuration to a temporary file and ensure it's signature matches
		tmp, err := ioutil.TempFile("", "catena-*-config"+ext)
		require.NoError(t, err)
		defer os.Remove(tmp.Name())

		require.NoError(t, loaded.DumpFile(tmp.Name()))
		tmp.Close()

		siga, err := getFileSignature(testpath)
		require.NoError(t, err)
		sigb, err := getFileSignature(tmp.Name())
		require.NoError(t, err)
		require.True(t, bytes.Equal(siga, sigb))
	}
}

func getFileSignature(path string) ([]byte, error) {
	hasher := sha256.New()
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err = io.Copy(hasher, f); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
