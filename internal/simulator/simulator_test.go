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
