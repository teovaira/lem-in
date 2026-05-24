package graph

import "testing"

func TestColonyConstruction(t *testing.T) {
	c := &Colony{
		Ants:  4,
		Start: "start",
		End:   "end",
		Rooms: map[string]*Room{
			"start": {Name: "start", X: 0, Y: 0},
			"a":     {Name: "a", X: 1, Y: 0},
			"end":   {Name: "end", X: 2, Y: 0},
		},
		Links: map[string][]string{
			"start": {"a"},
			"a":     {"start", "end"},
			"end":   {"a"},
		},
		RawLines: []string{"4", "##start", "start 0 0", "##end", "end 2 0", "a 1 0", "start-a", "a-end"},
	}

	if c.Ants != 4 {
		t.Errorf("want Ants=4, got %d", c.Ants)
	}
	if c.Start != "start" {
		t.Errorf("want Start=start, got %s", c.Start)
	}
	if c.End != "end" {
		t.Errorf("want End=end, got %s", c.End)
	}
	if len(c.Rooms) != 3 {
		t.Errorf("want 3 rooms, got %d", len(c.Rooms))
	}
	if c.Rooms["a"].X != 1 || c.Rooms["a"].Y != 0 {
		t.Errorf("room a: want X=1 Y=0, got X=%d Y=%d", c.Rooms["a"].X, c.Rooms["a"].Y)
	}
	if len(c.RawLines) != 8 {
		t.Errorf("want 8 RawLines, got %d", len(c.RawLines))
	}
}

func TestNeighboursDirections(t *testing.T) {
	c := &Colony{
		Links: map[string][]string{
			"start": {"a", "b"},
			"a":     {"start", "end"},
			"b":     {"start"},
			"end":   {"a"},
		},
	}

	tests := []struct {
		room string
		want []string
	}{
		{"start", []string{"a", "b"}},
		{"a", []string{"start", "end"}},
		{"end", []string{"a"}},
		{"b", []string{"start"}},
	}

	for _, tt := range tests {
		got := c.Neighbours(tt.room)
		if len(got) != len(tt.want) {
			t.Errorf("Neighbours(%q): want %v, got %v", tt.room, tt.want, got)
			continue
		}
		gotMap := make(map[string]bool)
		for _, n := range got {
			gotMap[n] = true
		}
		for _, n := range tt.want {
			if !gotMap[n] {
				t.Errorf("Neighbours(%q): missing %q, got %v", tt.room, n, got)
			}
		}
	}
}

func TestNeighboursUnknownRoom(t *testing.T) {
	c := &Colony{
		Links: map[string][]string{
			"a": {"b"},
		},
	}
	got := c.Neighbours("z")
	if got != nil {
		t.Errorf("Neighbours of unknown room: want nil, got %v", got)
	}
}

func TestHasLinkPresent(t *testing.T) {
	c := &Colony{
		Links: map[string][]string{
			"a": {"b", "c"},
			"b": {"a"},
			"c": {"a"},
		},
	}

	if !c.HasLink("a", "b") {
		t.Error("want HasLink(a,b)=true")
	}
	if !c.HasLink("b", "a") {
		t.Error("want HasLink(b,a)=true (undirected)")
	}
	if !c.HasLink("a", "c") {
		t.Error("want HasLink(a,c)=true")
	}
}

func TestHasLinkAbsent(t *testing.T) {
	c := &Colony{
		Links: map[string][]string{
			"a": {"b"},
			"b": {"a"},
			"c": {},
		},
	}

	if c.HasLink("a", "c") {
		t.Error("want HasLink(a,c)=false")
	}
	if c.HasLink("b", "c") {
		t.Error("want HasLink(b,c)=false")
	}
	if c.HasLink("z", "a") {
		t.Error("want HasLink(z,a)=false for unknown room")
	}
}
