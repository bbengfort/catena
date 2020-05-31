package logs

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

// NoColor returns true if stdout is not a terminal.
func NoColor() bool {
	return os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
}

const colorReset = "\033[0m"

// Color constants for defining ANSI terminal output colors
const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

var (
	colors = []string{
		LevelTrace:   colorSeq(colorBlue),
		LevelDebug:   colorSeq(colorMagenta),
		LevelInfo:    colorSeq(colorGreen),
		LevelCaution: colorSeq(colorYellow),
		LevelStatus:  colorSeq(colorCyan),
		LevelWarn:    colorSeq(colorRed),
	}
	boldcolors = []string{
		LevelTrace:   colorSeqBold(colorBlue),
		LevelDebug:   colorSeqBold(colorMagenta),
		LevelInfo:    colorSeqBold(colorGreen),
		LevelCaution: colorSeqBold(colorYellow),
		LevelStatus:  colorSeqBold(colorCyan),
		LevelWarn:    colorSeqBold(colorRed),
	}
	statuscolors = []string{
		0: colorSeqBold(colorRed),
		1: colorSeq(colorCyan),
		2: colorSeq(colorGreen),
		3: colorSeq(colorMagenta),
		4: colorSeq(colorYellow),
		5: colorSeq(colorRed),
	}
)

func colorize(level uint8, bold bool, msg string, a ...interface{}) string {
	var color string
	if bold {
		color = boldcolors[level]
	} else {
		color = colors[level]
	}

	if len(a) > 0 {
		return color + fmt.Sprintf(msg, a...) + colorReset
	}
	return color + msg + colorReset
}

func colorSeq(color int) string {
	return fmt.Sprintf("\033[%dm", color)
}

func colorSeqBold(color int) string {
	return fmt.Sprintf("\033[%d;1m", color)
}
