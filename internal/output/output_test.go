package output

import (
	"bytes"
	"strings"
	"testing"

	"lemin/internal/graph"
)

// TestPrintRawLines verifies that every line stored in c.RawLines
// appears in the output before the blank separator line.
func TestPrintRawLines(t *testing.T) {
	c := &graph.Colony{
		RawLines: []string{"4", "##start", "0 0 3", "##end", "1 8 3", "0-1"},
	}
	turns := [][]graph.Move{}

	var buf bytes.Buffer
	Print(&buf, c, turns)

	output := buf.String()
	for _, line := range c.RawLines {
		if !strings.Contains(output, line) {
			t.Errorf("expected output to contain %q", line)
		}
	}
}
