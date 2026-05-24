package simulator

import "lemin/internal/graph"

// Simulate assigns ants to paths and produces the full turn sequence.
// paths must be sorted shortest-first (as returned by FindPaths).
// Outer slice = turns. Inner slice = moves in that turn, sorted by AntID ascending.
func Simulate(c *graph.Colony, paths []graph.Path) [][]graph.Move {
	return nil
}
