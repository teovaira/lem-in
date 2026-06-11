package simulator

import (
	"testing"

	"lemin/internal/graph"
	"lemin/internal/parser"
	"lemin/internal/pathfinder"
)

func TestSimulateOneAntDirectLink(t *testing.T) {
	c := &graph.Colony{
		Ants:  1,
		Start: "start",
		End:   "end",
		Rooms: map[string]*graph.Room{
			"start": {Name: "start"},
			"end":   {Name: "end"},
		},
		Links: map[string][]string{},
	}
	paths := []graph.Path{{"start", "end"}}
	turns := Simulate(c, paths)
	if len(turns) != 1 {
		t.Fatalf("want 1 turn, got %d", len(turns))
	}
	if len(turns[0]) != 1 {
		t.Fatalf("turn 1: want 1 move, got %d", len(turns[0]))
	}
	if turns[0][0].AntID != 1 || turns[0][0].Room != "end" {
		t.Errorf("turn 1: want L1-end, got L%d-%s", turns[0][0].AntID, turns[0][0].Room)
	}
}

func TestSimulateOneAntOnePath(t *testing.T) {
	c := &graph.Colony{
		Ants:  1,
		Start: "start",
		End:   "end",
		Rooms: map[string]*graph.Room{
			"start": {Name: "start", X: 0, Y: 0},
			"a":     {Name: "a", X: 1, Y: 0},
			"b":     {Name: "b", X: 2, Y: 0},
			"end":   {Name: "end", X: 3, Y: 0},
		},
		Links: map[string][]string{
			"start": {"a"},
			"a":     {"start", "b"},
			"b":     {"a", "end"},
			"end":   {"b"},
		},
	}

	paths := []graph.Path{
		{"start", "a", "b", "end"},
	}

	turns := Simulate(c, paths)

	if len(turns) != 3 {
		t.Fatalf("want 3 turns, got %d", len(turns))
	}

	want := []graph.Move{
		{AntID: 1, Room: "a"},
		{AntID: 1, Room: "b"},
		{AntID: 1, Room: "end"},
	}
	for i, move := range want {
		if len(turns[i]) != 1 {
			t.Fatalf("turn %d: want 1 move, got %d", i+1, len(turns[i]))
		}
		got := turns[i][0]
		if got.AntID != move.AntID || got.Room != move.Room {
			t.Errorf("turn %d: want L%d-%s, got L%d-%s", i+1, move.AntID, move.Room, got.AntID, got.Room)
		}
	}
}

// verifyNoCollision checks no two ants share the same intermediate room on the same turn.
// Start and end are excluded — they may hold multiple ants simultaneously.
func verifyNoCollision(t *testing.T, turns [][]graph.Move, start, end string) {
	t.Helper()
	for ti, turn := range turns {
		seen := make(map[string]int)
		for _, m := range turn {
			if m.Room == start || m.Room == end {
				continue
			}
			if prev, ok := seen[m.Room]; ok {
				t.Errorf("turn %d: ants %d and %d both in room %q", ti+1, prev, m.AntID, m.Room)
			}
			seen[m.Room] = m.AntID
		}
	}
}

// verifyTunnelOnce checks no tunnel is used more than once per turn.
// It reconstructs which tunnel each ant used by tracking each ant's previous room.
func verifyTunnelOnce(t *testing.T, c *graph.Colony, turns [][]graph.Move) {
	t.Helper()
	prevRoom := make(map[int]string)
	for id := 1; id <= c.Ants; id++ {
		prevRoom[id] = c.Start
	}
	for ti, turn := range turns {
		type tunnel struct{ a, b string }
		used := make(map[tunnel]int)
		for _, m := range turn {
			from := prevRoom[m.AntID]
			a, b := from, m.Room
			if a > b {
				a, b = b, a
			}
			key := tunnel{a, b}
			if prev, ok := used[key]; ok {
				t.Errorf("turn %d: tunnel %s-%s used by ants %d and %d", ti+1, a, b, prev, m.AntID)
			}
			used[key] = m.AntID
			prevRoom[m.AntID] = m.Room
		}
	}
}

// verifyAllAntsReachEnd checks that every ant ID 1..c.Ants appears in end
// exactly once across all turns, and never moves after reaching end.
func verifyAllAntsReachEnd(t *testing.T, c *graph.Colony, turns [][]graph.Move) {
	t.Helper()
	reached := make(map[int]bool)
	for ti, turn := range turns {
		for _, m := range turn {
			if reached[m.AntID] {
				t.Errorf("turn %d: ant %d moved after reaching end", ti+1, m.AntID)
			}
			if m.Room == c.End {
				reached[m.AntID] = true
			}
		}
	}
	for id := 1; id <= c.Ants; id++ {
		if !reached[id] {
			t.Errorf("ant %d never reached end", id)
		}
	}
}

func TestSimulateFourAntsTwoPaths(t *testing.T) {
	c := &graph.Colony{
		Ants:  4,
		Start: "start",
		End:   "end",
		Rooms: map[string]*graph.Room{
			"start": {Name: "start"},
			"a":     {Name: "a"},
			"b":     {Name: "b"},
			"end":   {Name: "end"},
		},
		Links: map[string][]string{},
	}
	paths := []graph.Path{
		{"start", "a", "end"},
		{"start", "b", "end"},
	}
	turns := Simulate(c, paths)
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateThreeAntsTwoPaths(t *testing.T) {
	c := &graph.Colony{
		Ants:  3,
		Start: "start",
		End:   "end",
		Rooms: map[string]*graph.Room{
			"start": {Name: "start"},
			"a":     {Name: "a"},
			"b":     {Name: "b"},
			"cc":    {Name: "cc"},
			"end":   {Name: "end"},
		},
		Links: map[string][]string{},
	}

	paths := []graph.Path{
		{"start", "a", "end"},
		{"start", "b", "cc", "end"},
	}

	turns := Simulate(c, paths)

	if len(turns) != 3 {
		t.Fatalf("want 3 turns, got %d", len(turns))
	}

	// turn 1: ant1→a, ant3→b
	// turn 2: ant1→end, ant2→a, ant3→cc
	// turn 3: ant2→end, ant3→end
	type move struct {
		id   int
		room string
	}
	wantTurns := [][]move{
		{{1, "a"}, {3, "b"}},
		{{1, "end"}, {2, "a"}, {3, "cc"}},
		{{2, "end"}, {3, "end"}},
	}

	for ti, wantMoves := range wantTurns {
		got := turns[ti]
		if len(got) != len(wantMoves) {
			t.Fatalf("turn %d: want %d moves, got %d: %v", ti+1, len(wantMoves), len(got), got)
		}
		gotMap := make(map[int]string)
		for _, m := range got {
			gotMap[m.AntID] = m.Room
		}
		for _, w := range wantMoves {
			if gotMap[w.id] != w.room {
				t.Errorf("turn %d: ant %d want room %q, got %q", ti+1, w.id, w.room, gotMap[w.id])
			}
		}
	}
}

func exampleTurns(t *testing.T, file string) (*graph.Colony, [][]graph.Move) {
	t.Helper()
	c, err := parser.Parse(file)
	if err != nil {
		t.Fatalf("%s: parse: %v", file, err)
	}
	paths, err := pathfinder.FindPaths(c)
	if err != nil {
		t.Fatalf("%s: pathfinder: %v", file, err)
	}
	return c, Simulate(c, paths)
}

func TestSimulateExample00(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example00.txt")
	if len(turns) > 6 {
		t.Errorf("want ≤6 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateExample01(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example01.txt")
	if len(turns) > 8 {
		t.Errorf("want ≤8 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateExample02(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example02.txt")
	if len(turns) > 11 {
		t.Errorf("want ≤11 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateExample03(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example03.txt")
	if len(turns) > 6 {
		t.Errorf("want ≤6 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateExample04(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example04.txt")
	if len(turns) > 6 {
		t.Errorf("want ≤6 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}

func TestSimulateExample05(t *testing.T) {
	c, turns := exampleTurns(t, "../../testdata/example05.txt")
	if len(turns) > 8 {
		t.Errorf("want ≤8 turns, got %d", len(turns))
	}
	verifyAllAntsReachEnd(t, c, turns)
	verifyNoCollision(t, turns, c.Start, c.End)
	verifyTunnelOnce(t, c, turns)
}
