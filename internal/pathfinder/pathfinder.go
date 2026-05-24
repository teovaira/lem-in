// Package pathfinder finds the optimal set of non-overlapping paths through a colony.
package pathfinder

import (
	"errors"
	"sort"
	"strings"

	"lemin/internal/graph"
)

type edge struct {
	to, cap, rev int
}

func addEdge(g [][]edge, u, v, cap int) {
	g[u] = append(g[u], edge{v, cap, len(g[v])})
	g[v] = append(g[v], edge{u, 0, len(g[u]) - 1})
}

// buildFlowGraph builds the node-split residual graph.
// Returns graph, source ID, sink ID, and room-name slice indexed by nodeID/2.
func buildFlowGraph(c *graph.Colony) (g [][]edge, src, snk int, names []string) {
	// assign a stable index to each room
	roomNames := make([]string, 0, len(c.Rooms))
	for name := range c.Rooms {
		roomNames = append(roomNames, name)
	}
	sort.Strings(roomNames)

	idx := make(map[string]int, len(roomNames))
	for i, name := range roomNames {
		idx[name] = i
	}
	names = roomNames

	n := len(roomNames)
	g = make([][]edge, n*2)
	for i := range g {
		g[i] = []edge{}
	}

	// node-split edges: room_in -> room_out
	for _, name := range roomNames {
		i := idx[name]
		in, out := i*2, i*2+1
		cap := 1
		if name == c.Start || name == c.End {
			cap = c.Ants
		}
		addEdge(g, in, out, cap)
	}

	// tunnel edges: a_out -> b_in (both directions, undirected)
	sortedKeys := make([]string, 0, len(c.Links))
	for k := range c.Links {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, a := range sortedKeys {
		neighbours := make([]string, len(c.Links[a]))
		copy(neighbours, c.Links[a])
		sort.Strings(neighbours)
		for _, b := range neighbours {
			if idx[a] < idx[b] { // add each undirected edge once
				addEdge(g, idx[a]*2+1, idx[b]*2, 1)
				addEdge(g, idx[b]*2+1, idx[a]*2, 1)
			}
		}
	}

	src = idx[c.Start] * 2
	snk = idx[c.End]*2 + 1
	return
}

// bfs finds one augmenting path from src to snk.
// Returns parent array (parent[node] = edge index used), or nil if no path.
func bfs(g [][]edge, src, snk int) []int {
	parent := make([]int, len(g))
	for i := range parent {
		parent[i] = -1
	}
	parent[src] = -2
	queue := []int{src}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for i, e := range g[u] {
			if e.cap > 0 && parent[e.to] == -1 {
				parent[e.to] = i<<20 | u // encode: edge index + from node
				if e.to == snk {
					return parent
				}
				queue = append(queue, e.to)
			}
		}
	}
	return nil
}

// extractPaths runs Edmonds-Karp to completion and returns all augmenting paths as room-name slices.
func extractPaths(c *graph.Colony, g [][]edge, src, snk int, names []string) []graph.Path {
	var paths []graph.Path
	for {
		parent := bfs(g, src, snk)
		if parent == nil {
			break
		}
		// push flow=1 and record path
		var nodes []int
		cur := snk
		for cur != src {
			encoded := parent[cur]
			from := encoded & ((1 << 20) - 1)
			ei := encoded >> 20
			g[from][ei].cap--
			g[cur][g[from][ei].rev].cap++
			nodes = append(nodes, cur)
			cur = from
		}
		nodes = append(nodes, src)
		// reverse
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
		// convert node IDs to room names: only _out nodes (odd) represent a room crossing
		var path graph.Path
		for _, node := range nodes {
			if node%2 == 1 { // _out node
				path = append(path, names[node/2])
			}
		}
		paths = append(paths, path)
	}
	return paths
}

// computeTurns returns the number of turns for nAnts ants on the given paths
// using greedy assignment (assign each ant to the path minimising its finish turn).
func computeTurns(paths []graph.Path, nAnts int) int {
	load := make([]int, len(paths))
	maxTurn := 0
	for ant := 1; ant <= nAnts; ant++ {
		best := 0
		bestTurn := len(paths[0]) - 1 + load[0]
		for i := 1; i < len(paths); i++ {
			t := len(paths[i]) - 1 + load[i]
			if t < bestTurn {
				bestTurn = t
				best = i
			}
		}
		load[best]++
		if bestTurn > maxTurn {
			maxTurn = bestTurn
		}
	}
	return maxTurn
}

// FindPaths returns the optimal set of non-overlapping paths minimising total turns for c.Ants ants.
// Returns error if no path exists from Start to End.
func FindPaths(c *graph.Colony) ([]graph.Path, error) {
	g, src, snk, names := buildFlowGraph(c)
	paths := extractPaths(c, g, src, snk, names)
	if len(paths) == 0 {
		return nil, errors.New("no path between start and end")
	}

	// sort paths: shortest first, tie-break lexicographically
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return strings.Join(paths[i], ",") < strings.Join(paths[j], ",")
	})

	// find the prefix [0:k] that minimises turns
	best := computeTurns(paths[:1], c.Ants)
	bestK := 1
	for k := 2; k <= len(paths); k++ {
		t := computeTurns(paths[:k], c.Ants)
		if t < best {
			best = t
			bestK = k
		}
	}
	return paths[:bestK], nil
}
