package output

import (
	"io"

	"lemin/internal/graph"
)

// Print writes the complete program output to w:
// RawLines verbatim, one blank line, then one line per turn.
func Print(w io.Writer, c *graph.Colony, turns [][]graph.Move) {
}
