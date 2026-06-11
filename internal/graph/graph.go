// Package graph defines the shared data types used across all lem-in packages.
package graph

// Room is a single node in the colony graph.
type Room struct {
	Name string // unique identifier; never starts with L or #
	X    int    // horizontal coordinate (signed)
	Y    int    // vertical coordinate (signed)
}

// Colony is the fully validated ant farm graph produced by the parser.
type Colony struct {
	Ants     int                 // number of ants to move (>= 1)
	Rooms    map[string]*Room    // keyed by room name
	Links    map[string][]string // adjacency list; both directions stored per undirected edge
	Start    string              // name of the ##start room
	End      string              // name of the ##end room
	RawLines []string            // every non-empty input line in original order, for verbatim output
}

// Path is an ordered sequence of room names from Start to End inclusive.
type Path []string

// Move represents one ant advancing one step on one turn.
type Move struct {
	AntID int    // 1-indexed ant identifier
	Room  string // destination room name
}

// Neighbours returns all room names directly connected by a tunnel to room.
// Returns nil if room is not present in the colony.
func (c *Colony) Neighbours(room string) []string {
	return c.Links[room]
}

// HasLink reports whether rooms a and b share a direct tunnel.
// Returns false if either room is unknown.
func (c *Colony) HasLink(a, b string) bool {
	for _, n := range c.Links[a] {
		if n == b {
			return true
		}
	}
	return false
}
