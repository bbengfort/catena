/*
Package logs implements simple hierarchical logging functionality for debugging and
logging. The package can write to any configured logger, but generally writes to stdout
for use with system logging. The log level specifies the verbosity of the output. For
example if the level is set to Info, then Debug and Trace messages will become no-ops
and ignored by the logger.

This package also provides a Caution log level - caution messages are only printed if
a specific threshold of messages has been reached. This helps to reduce the number of
repeated messages (e.g. connection down) that occur in logging while still giving
effective debugging and systems administration feedback to the user.
*/
package logs

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// Logging levels for specify the verbosity of log output. The higher the level, the
// less verbose the output is, e.g. will log messages >= to the specified level.
const (
	LevelTrace uint8 = iota
	LevelDebug
	LevelInfo
	LevelCaution
	LevelStatus
	LevelWarn
	LevelSilent
)

var levelStrings = [...]string{
	"TRACE", "DEBUG", "INFO", "CAUTION", "STATUS", "WARN", "SILENT",
}

// Default values for new loggers that are created by the package
const (
	DefaultLogLevel         = LevelInfo
	DefaultCautionThreshold = 80
	DefaultUseColor         = true
)

// New returns a new logger that can be configured to write its output.
func New(prefix string) *Logger {
	l := &Logger{
		prefix:         prefix,
		timestamp:      " [2006-01-02 15:04:05] ",
		colorize:       NoColor(),
		logLevel:       LevelInfo,
		logger:         log.New(os.Stdout, "", 0),
		cautionCounter: new(counter),
	}
	l.cautionCounter.init()
	return l
}

// Logger wraps a log.Logger to provide advanced functionality such as logging levels,
// terminal colorization, and caution thresholding. Different loggers can be configured
// to have different output formats and layouts.
type Logger struct {
	prefix         string      // the prefix to add to output (usually module name)
	timestamp      string      // the timestamp format layout
	colorize       bool        // whether or not to colorize the output
	logLevel       uint8       // the level that the logger outputs at
	logger         *log.Logger // the underlying logger being wrapped
	cautionCounter *counter    // the caution counter used to reduce verbose input
}

//===========================================================================
// Interact with debug output
//===========================================================================

// ParseLevel returns the log level from the level string.
func ParseLevel(level string) (uint8, error) {
	// convert into all uppercase with no whitespace
	level = strings.ToUpper(strings.TrimSpace(level))
	for i, name := range levelStrings {
		if level == name {
			return uint8(i), nil
		}
	}
	return 0, fmt.Errorf("could not parse %q as a log level", level)
}

// LogLevel returns a string representation of the current level.
func (l *Logger) LogLevel() string {
	return levelStrings[l.logLevel]
}

// SetLogLevel modifies the log level for messages at runtime. Ensures that
// the highest level that can be set is the trace level. This function is
// often called from outside of the package in an init() function to define
// how logging is handled in the console.
func (l *Logger) SetLogLevel(level uint8) {
	if level > LevelSilent {
		level = LevelSilent
	}
	l.logLevel = level
}

// SetBackend updates the logger to use a different writer than os.Stdout
func (l *Logger) SetBackend(out io.Writer) {
	l.logger = log.New(out, "", 0)
}

// SetTimestamp updates the logger to use a different timestamp layout function
func (l *Logger) SetTimestamp(layout string) {
	l.timestamp = layout
}

// SetCautionThreshold to the specified number of messages before print.
func (l *Logger) SetCautionThreshold(threshold uint) {
	l.cautionCounter.threshold = threshold
}

// Colorize returns whether or not the output is being colorized
func (l *Logger) Colorize() bool {
	return l.colorize
}

// EnableColors turns terminal colorization on (overrides is terminal check).
func (l *Logger) EnableColors() {
	l.colorize = true
}

// DisableColors turns terminal colorization off (overrides is terminal check).
func (l *Logger) DisableColors() {
	l.colorize = false
}

//===========================================================================
// Debugging output functions
//===========================================================================

// Print to the standard logger at the specified level. Arguments are handled
// in the manner of log.Printf, but a newline is appended.
func (l *Logger) print(level uint8, msg string, a ...interface{}) {
	if l.logLevel <= level {
		// Build a string buffer with the estimated length
		var buf strings.Builder
		estlen := len(msg) + len(l.prefix) + (2 * len(a)) + len(l.timestamp)
		if l.colorize {
			estlen += 8
		}
		buf.Grow(estlen)

		// Colorize the output
		if l.colorize {
			buf.WriteString(colors[level])
		}

		// Write the prefix
		buf.WriteString(l.prefix)

		// Write the timestamp
		if l.timestamp != "" {
			buf.WriteString(time.Now().Format(l.timestamp))
		}

		// Execute string format if arguments are provided
		fmt.Fprintf(&buf, msg, a...)

		// Reset the colorization
		if l.colorize {
			buf.WriteString(colorReset)
		}

		// Log the message
		l.logger.Print(buf.String())
	}
}

// Warn prints to the standard logger if level is warn or greater; arguments
// are handled in the manner of log.Printf, but a newline is appended.
func (l *Logger) Warn(msg string, a ...interface{}) {
	l.print(LevelWarn, msg, a...)
}

// Warne is a helper function to simply warn about an error received.
func (l *Logger) Warne(err error) {
	l.Warn(err.Error())
}

// Status prints to the standard logger if level is status or greater;
// arguments are handled in the manner of log.Printf, but a newline is appended.
func (l *Logger) Status(msg string, a ...interface{}) {
	l.print(LevelStatus, msg, a...)
}

// Caution prints to the standard logger if the level is caution or greater and if the
// number of times caution has been called with the same message has reached the
// threshold. This reduces the number of repeated log output messages while still
// allowing the system to report valuable information.
func (l *Logger) Caution(msg string, a ...interface{}) {
	if l.logLevel > LevelCaution {
		// Don't waste memory if the log level is set above caution.
		return
	}

	msg = fmt.Sprintf(msg, a...)
	if l.cautionCounter.log(msg) {
		l.print(LevelCaution, msg)
	}
}

// Info prints to the standard logger if level is info or greater; arguments
// are handled in the manner of log.Printf, but a newline is appended.
func (l *Logger) Info(msg string, a ...interface{}) {
	l.print(LevelInfo, msg, a...)
}

// Debug prints to the standard logger if level is debug or greater;
// arguments are handled in the manner of log.Printf, but a newline is appended.
func (l *Logger) Debug(msg string, a ...interface{}) {
	l.print(LevelDebug, msg, a...)
}

// Trace prints to the standard logger if level is trace or greater;
// arguments are handled in the manner of log.Printf, but a newline is appended.
func (l *Logger) Trace(msg string, a ...interface{}) {
	l.print(LevelTrace, msg, a...)
}
