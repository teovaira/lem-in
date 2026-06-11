package parser

import (
	"os"
	"strings"
	"testing"

	"lemin/internal/graph"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "lemin_test_*.txt")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("WriteString: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestParseExample00(t *testing.T) {
	c, err := Parse("../../testdata/example00.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Ants != 4 {
		t.Errorf("want Ants=4, got %d", c.Ants)
	}
	if c.Start != "0" {
		t.Errorf("want Start=0, got %q", c.Start)
	}
	if c.End != "1" {
		t.Errorf("want End=1, got %q", c.End)
	}
	if len(c.Rooms) != 4 {
		t.Errorf("want 4 rooms, got %d", len(c.Rooms))
	}
	if !c.HasLink("0", "2") || !c.HasLink("2", "3") || !c.HasLink("3", "1") {
		t.Error("missing expected links")
	}
}

func TestParseExample03(t *testing.T) {
	c, err := Parse("../../testdata/example03.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Start == "" {
		t.Error("Start should not be empty")
	}
	if c.End == "" {
		t.Error("End should not be empty")
	}
	if len(c.Rooms) == 0 {
		t.Error("Rooms should not be empty")
	}
}

func TestParseNumericRoomNames(t *testing.T) {
	path := writeTempFile(t, "1\n##start\n0 1 4\n##end\n1 5 4\n0-1\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := c.Rooms["0"]
	if !ok {
		t.Fatal("room named '0' not found")
	}
	if r.X != 1 || r.Y != 4 {
		t.Errorf("room '0': want X=1 Y=4, got X=%d Y=%d", r.X, r.Y)
	}
	if c.Start != "0" {
		t.Errorf("want Start=0, got %q", c.Start)
	}
}

func TestParseRoomNamedStart(t *testing.T) {
	path := writeTempFile(t, "1\n##start\nstart 1 6\n##end\nend 5 6\nstart-end\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := c.Rooms["start"]; !ok {
		t.Error("room named 'start' not found")
	}
	if c.Start != "start" {
		t.Errorf("want Start='start', got %q", c.Start)
	}
}

func TestParseCommentsInRawLines(t *testing.T) {
	path := writeTempFile(t, "2\n#a comment\n##start\nroom1 0 0\n##end\nroom2 1 0\nroom1-room2\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := strings.Join(c.RawLines, "\n")
	if !strings.Contains(raw, "#a comment") {
		t.Error("single-# comment must be in RawLines")
	}
	if !strings.Contains(raw, "##start") {
		t.Error("##start must be in RawLines")
	}
	if !strings.Contains(raw, "##end") {
		t.Error("##end must be in RawLines")
	}
}

func TestParseRawLinesOrder(t *testing.T) {
	path := writeTempFile(t, "3\n##start\na 0 0\n##end\nb 1 0\na-b\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(c.RawLines) == 0 {
		t.Fatal("RawLines is empty")
	}
	if c.RawLines[0] != "3" {
		t.Errorf("RawLines[0]: want ant count '3', got %q", c.RawLines[0])
	}
}

func TestParseMultiPath(t *testing.T) {
	path := writeTempFile(t, "2\n##start\ns 0 0\n##end\ne 4 0\na 1 1\nb 1 -1\ns-a\na-e\ns-b\nb-e\n")
	_, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseExample05CommentsPreserved(t *testing.T) {
	c, err := Parse("../../testdata/example05.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	raw := strings.Join(c.RawLines, "\n")
	if !strings.Contains(raw, "#rooms") {
		t.Error("example05 #rooms comment must be preserved in RawLines")
	}
}

func TestParseSkipsBlankLines(t *testing.T) {
	path := writeTempFile(t, "1\n\n##start\na 0 0\n\n##end\nb 1 0\na-b\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, line := range c.RawLines {
		if line == "" {
			t.Error("blank lines must NOT appear in RawLines")
		}
	}
}

func TestParseStartFollowedByCommentThenRoom(t *testing.T) {
	path := writeTempFile(t, "1\n##start\n#comment here\nmyroom 1 2\n##end\nother 3 4\nmyroom-other\n")
	c, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Start != "myroom" {
		t.Errorf("want Start=myroom, got %q", c.Start)
	}
	raw := strings.Join(c.RawLines, "\n")
	if !strings.Contains(raw, "#comment here") {
		t.Error("comment between ##start and room must be in RawLines")
	}
}

func TestParseErrorNoFile(t *testing.T) {
	_, err := Parse("/nonexistent/path/lemin_no_such_file.txt")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParseErrorBadAntCount(t *testing.T) {
	path := writeTempFile(t, "abc\n##start\na 0 0\n##end\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for non-integer ant count")
	}
}

func TestParseErrorZeroAnts(t *testing.T) {
	path := writeTempFile(t, "0\n##start\na 0 0\n##end\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for zero ants")
	}
}

func TestParseErrorNegativeAnts(t *testing.T) {
	path := writeTempFile(t, "-3\n##start\na 0 0\n##end\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for negative ants")
	}
}

func TestParseErrorNoStart(t *testing.T) {
	path := writeTempFile(t, "1\na 0 0\n##end\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for missing ##start")
	}
	if !strings.Contains(err.Error(), "no start room") {
		t.Errorf("want 'no start room' in error, got %q", err.Error())
	}
}

func TestParseErrorNoEnd(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for missing ##end")
	}
	if !strings.Contains(err.Error(), "no end room") {
		t.Errorf("want 'no end room' in error, got %q", err.Error())
	}
}

func TestParseErrorMultipleStart(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##start\nb 1 0\n##end\nc 2 0\na-c\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for multiple ##start")
	}
}

func TestParseErrorMultipleEnd(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\nb 1 0\n##end\nc 2 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for multiple ##end")
	}
}

func TestParseErrorDuplicateRoom(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\na 1 1\n##end\nb 2 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for duplicate room name")
	}
}

func TestParseErrorRoomStartsWithL(t *testing.T) {
	path := writeTempFile(t, "1\n##start\nLair 0 0\n##end\nb 1 0\nLair-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for room name starting with L")
	}
}

func TestParseErrorRoomStartsWithHash(t *testing.T) {
	// Lines starting with # are treated as comments before reaching parseRoom.
	// We test parseRoom directly (white-box) to verify the validation exists.
	c := &graph.Colony{
		Rooms: make(map[string]*graph.Room),
		Links: make(map[string][]string),
	}
	_, err := parseRoom("#x 1 2", c)
	if err == nil {
		t.Fatal("expected error for room name starting with #")
	}
}

func TestParseErrorSelfLink(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\nb 1 0\na-a\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for self-link")
	}
}

func TestParseErrorDuplicateLink(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\nb 1 0\na-b\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for duplicate link")
	}
}

func TestParseErrorUnknownRoom(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\nb 1 0\na-z\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for link referencing undeclared room")
	}
}

func TestParseErrorInvalidCoordinates(t *testing.T) {
	path := writeTempFile(t, "1\n##start\nroom abc def\n##end\nb 1 0\nroom-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for non-integer room coordinates")
	}
	if !strings.Contains(err.Error(), "invalid room coordinates") {
		t.Errorf("want 'invalid room coordinates' in error, got %q", err.Error())
	}
}

func TestParseErrorGarbageLine(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\nb 1 0\na-b\nhello\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for garbage line")
	}
}

func TestParseErrorStartFollowedByLink(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na-b\na 0 0\n##end\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error when ##start is followed by a link line")
	}
	if !strings.Contains(err.Error(), "no start room") {
		t.Errorf("want 'no start room' in error, got %q", err.Error())
	}
}

func TestParseErrorEndFollowedByLink(t *testing.T) {
	path := writeTempFile(t, "1\n##start\na 0 0\n##end\na-b\nb 1 0\na-b\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error when ##end is followed by a link line")
	}
	if !strings.Contains(err.Error(), "no end room") {
		t.Errorf("want 'no end room' in error, got %q", err.Error())
	}
}
