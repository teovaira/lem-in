package simulator

import "lemin/internal/graph"

// Simulate assigns ants to paths and produces the full turn sequence.
// paths must be sorted shortest-first (as returned by FindPaths).
// Outer slice = turns. Inner slice = moves in that turn, sorted by AntID ascending.
func Simulate(c *graph.Colony, paths []graph.Path) [][]graph.Move {
	// antPath[id] = which path index ant id was assigned to (1-indexed ant IDs).
	antPath := make([]int, c.Ants+1)
	// antOffset[id] = how many turns ant id waits before entering (0 = enters on turn 1).
	antOffset := make([]int, c.Ants+1)
	// load[i] = how many ants are already queued on paths[i].
	load := make([]int, len(paths))

	// Greedy assignment: for each ant pick the path that gives the earliest finish turn.
	// Finish turn for ant on paths[i] = len(paths[i]) - 1 + load[i].
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

	// Compute the total number of turns needed.
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

	// Generate moves turn by turn.
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
