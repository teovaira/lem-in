package graph

// Room is a single node in the colony graph.
type Room struct {
	Name string
	X    int
	Y    int
}

// Colony is the fully validated ant farm graph.
type Colony struct {
	Ants     int
	Rooms    map[string]*Room
	Links    map[string][]string
	Start    string
	End      string
	RawLines []string
}

// Path is an ordered sequence of room names from Start to End inclusive.
type Path []string

// Move represents one ant advancing one step on one turn.
type Move struct {
	AntID int
	Room  string
}

// Neighbours returns all rooms directly connected by a tunnel to the given room.
func (c *Colony) Neighbours(room string) []string {
	return c.Links[room]
}

// HasLink returns true if rooms a and b share a direct tunnel.
func (c *Colony) HasLink(a, b string) bool {
	for _, n := range c.Links[a] {
		if n == b {
			return true
		}
	}
	return false
}
