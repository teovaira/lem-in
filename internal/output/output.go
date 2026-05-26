package output

import (
	"fmt"
	"io"
	"strings"

	"lemin/internal/graph"
)

// Print writes the complete program output to w:
// RawLines verbatim, one blank line, then one line per turn.
func Print(w io.Writer, c *graph.Colony, turns [][]graph.Move) {
	for _, line := range c.RawLines {
		fmt.Fprintln(w, line)
	}
	fmt.Fprintln(w)
	for _, turn := range turns {
		parts := make([]string, len(turn))
		for i, m := range turn {
			parts[i] = fmt.Sprintf("L%d-%s", m.AntID, m.Room)
		}
		fmt.Fprintln(w, strings.Join(parts, " "))
	}
}
