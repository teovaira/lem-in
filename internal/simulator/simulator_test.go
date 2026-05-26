package simulator

import (
	"testing"

	"lemin/internal/graph"
)

// TestSimulateOneAntOnePath is the simplest possible simulator test.
// One ant. One path: start -> a -> b -> end.
// The path has length 4 (4 room names), so the ant makes 3 moves.
// Expected turns:
//
//	turn 1: L1-a
//	turn 2: L1-b
//	turn 3: L1-end
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

// TestSimulateAllAntsReachEnd verifies every ant reaches end exactly once.
func TestSimulateAllAntsReachEnd(t *testing.T) {
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
}

// TestSimulateThreeAntsTwoPaths verifies staggered pipeline assignment.
// 3 ants, 2 paths: path0=[start,a,end] (len 3), path1=[start,b,c,end] (len 4).
// Greedy: ant1→path0 offset0, ant2→path0 offset1, ant3→path1 offset0.
// Total turns = 3.
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
	type move struct{ id int; room string }
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
