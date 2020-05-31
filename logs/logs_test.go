package logs_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/bbengfort/catena/logs"
	"github.com/stretchr/testify/require"
)

func ExampleLogger() {
	logger := logs.New("[test] ")

	logger.SetBackend(os.Stdout)
	logger.SetLogLevel(logs.LevelInfo)
	logger.SetCautionThreshold(2)
	logger.SetTimestamp("")
	logger.DisableColors()

	logger.Trace("routine %s happening", "thing")
	logger.Debug("sending message #%d from %s to %s", 42, "me", "you")
	logger.Info("listening on %s", "127.0.0.1")
	logger.Caution("could not reach %s -- connection is down", "uptime.robot")
	logger.Status("completed %d out of %d tasks", 42, 121)
	logger.Warn("limit of %d queries reached", 21)
	logger.Warne(errors.New("something bad happened"))

	logger.Caution("could not reach %s -- connection is down", "uptime.robot")
	logger.Caution("could not reach %s -- connection is down", "uptime.robot")

	// Output:
	// [test] listening on 127.0.0.1
	// [test] could not reach uptime.robot -- connection is down
	// [test] completed 42 out of 121 tasks
	// [test] limit of 21 queries reached
	// [test] something bad happened
	// [test] could not reach uptime.robot -- connection is down
}

func TestColors(t *testing.T) {
	logger := logs.New("")

	buf := bytes.NewBuffer(make([]byte, 0, 53))
	logger.SetBackend(buf)
	logger.EnableColors()
	logger.SetLogLevel(logs.LevelTrace)
	logger.SetTimestamp("")

	logger.Trace("blue")
	logger.Debug("magenta")
	logger.Info("green")
	logger.Caution("yellow")
	logger.Status("cyan")
	logger.Warn("red")

	expected := "\x1b[34mblue\x1b[0m\n\x1b[35mmagenta\x1b[0m\n\x1b[32mgreen\x1b[0m\n\x1b[33myellow\x1b[0m\n\x1b[36mcyan\x1b[0m\n\x1b[31mred\x1b[0m\n"
	require.Equal(t, expected, string(buf.Bytes()))
}
