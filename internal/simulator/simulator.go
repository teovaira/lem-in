// Package simulator assigns ants to paths and generates the turn-by-turn move sequence.
package simulator

import "lemin/internal/graph"

// Simulate assigns the ants in c to paths using greedy scheduling and returns
// the complete move sequence as a slice of turns.
// paths must be sorted shortest-first (as returned by FindPaths).
// Each inner slice contains the moves for one turn, ordered by AntID ascending.
// Moves to the start room are never recorded; an ant that reaches End stops moving.
func Simulate(c *graph.Colony, paths []graph.Path) [][]graph.Move {
	antPath := make([]int, c.Ants+1)
	antOffset := make([]int, c.Ants+1)
	load := make([]int, len(paths))

	// Greedy assignment: assign each ant to the path with the earliest finish turn.
	for id := 1; id <= c.Ants; id++ {
		best := 0
		bestTurns := len(paths[0]) - 1 + load[0]
		for i := 1; i < len(paths); i++ {
			t := len(paths[i]) - 1 + load[i]
			if t < bestTurns {
				bestTurns = t
				best = i
			}
		}
		antPath[id] = best
		antOffset[id] = load[best]
		load[best]++
	}

	totalTurns := 0
	for id := 1; id <= c.Ants; id++ {
		finish := antOffset[id] + len(paths[antPath[id]]) - 1
		if finish > totalTurns {
			totalTurns = finish
		}
	}

	turns := make([][]graph.Move, totalTurns)
	for t := range turns {
		turns[t] = []graph.Move{}
	}

	// Ant id on path p with offset o:
	//   - enters path on turn o+1 (moves to p[1])
	//   - reaches end on turn o+len(p)-1 (moves to p[len(p)-1])
	//   - on turn T it is at p[T-o], for o < T <= o+len(p)-1
	for id := 1; id <= c.Ants; id++ {
		p := paths[antPath[id]]
		o := antOffset[id]
		for T := o + 1; T <= o+len(p)-1; T++ {
			room := p[T-o]
			turns[T-1] = append(turns[T-1], graph.Move{AntID: id, Room: room})
		}
	}

	return turns
}
