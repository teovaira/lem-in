package output

import (
	"bytes"
	"strings"
	"testing"

	"lemin/internal/graph"
)

func makeColony(rawLines []string) *graph.Colony {
	return &graph.Colony{RawLines: rawLines}
}

func TestPrintRawLines(t *testing.T) {
	c := makeColony([]string{"4", "##start", "0 0 3", "##end", "1 8 3", "0-1"})
	var buf bytes.Buffer
	Print(&buf, c, [][]graph.Move{})
	out := buf.String()
	for _, line := range c.RawLines {
		if !strings.Contains(out, line) {
			t.Errorf("expected output to contain %q", line)
		}
	}
}

func TestPrintBlankLineSeparator(t *testing.T) {
	c := makeColony([]string{"1", "##start", "0 0 0", "##end", "1 1 0", "0-1"})
	turns := [][]graph.Move{{{AntID: 1, Room: "1"}}}
	var buf bytes.Buffer
	Print(&buf, c, turns)
	lines := strings.Split(buf.String(), "\n")
	blankCount := 0
	for _, l := range lines {
		if l == "" {
			blankCount++
		}
	}
	// exactly 2 empty strings: the blank separator and the trailing newline after last turn
	if blankCount != 2 {
		t.Errorf("want exactly 1 blank separator line (2 empty splits), got %d empty strings in: %q", blankCount, buf.String())
	}
}

func TestPrintTurnFormat(t *testing.T) {
	c := makeColony([]string{"2"})
	turns := [][]graph.Move{
		{{AntID: 1, Room: "a"}, {AntID: 2, Room: "b"}},
		{{AntID: 1, Room: "end"}, {AntID: 2, Room: "end"}},
	}
	var buf bytes.Buffer
	Print(&buf, c, turns)
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if lines[2] != "L1-a L2-b" {
		t.Errorf("turn 1: want %q, got %q", "L1-a L2-b", lines[2])
	}
	if lines[3] != "L1-end L2-end" {
		t.Errorf("turn 2: want %q, got %q", "L1-end L2-end", lines[3])
	}
}

func TestPrintNoTrailingBlankLine(t *testing.T) {
	c := makeColony([]string{"1"})
	turns := [][]graph.Move{{{AntID: 1, Room: "end"}}}
	var buf bytes.Buffer
	Print(&buf, c, turns)
	out := buf.String()
	if strings.HasSuffix(out, "\n\n") {
		t.Errorf("output has trailing blank line: %q", out)
	}
}

func TestPrintUsesWriter(t *testing.T) {
	c := makeColony([]string{"sentinel-line"})
	var buf bytes.Buffer
	Print(&buf, c, [][]graph.Move{})
	if !strings.Contains(buf.String(), "sentinel-line") {
		t.Error("Print did not write to the provided writer")
	}
}
