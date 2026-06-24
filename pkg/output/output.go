package output

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Success = color.New(color.FgGreen).SprintfFunc()
	Error   = color.New(color.FgRed).SprintfFunc()
	Info    = color.New(color.FgCyan).SprintfFunc()
	Warn    = color.New(color.FgYellow).SprintfFunc()
)

// FileCreated prints a "created" line for a file.
func FileCreated(path string) {
	fmt.Printf("  %s  %s\n", Success("CREATE"), path)
}

// FileSkipped prints a "skipped" line for an existing file.
func FileSkipped(path string) {
	fmt.Printf("  %s   %s\n", Warn("SKIP"), path)
}

// FileOverwritten prints an "overwritten" line.
func FileOverwritten(path string) {
	fmt.Printf("  %s  %s\n", Warn("OVERWRITE"), path)
}

// DirCreated prints a "created directory" line.
func DirCreated(path string) {
	fmt.Printf("  %s  %s\n", Success("DIR"), path)
}
