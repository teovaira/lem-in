// Package pathfinder finds the optimal set of vertex-disjoint paths through a colony
// using Edmonds-Karp max-flow with node splitting, minimising total simulation turns.
package pathfinder

import (
	"errors"
	"sort"
	"strings"

	"lemin/internal/graph"
)

// edge is a directed edge in the residual flow graph.
// origCap records the original capacity so tracePath can identify edges that carry flow
// (origCap > 0 && cap < origCap) after maxFlow has run to completion.
type edge struct {
	to, cap, origCap, rev int
}

// addEdge adds a directed edge u→v with the given capacity to g,
// along with its zero-capacity reverse edge v→u for the residual graph.
func addEdge(g [][]edge, u, v, cap int) {
	revOfFwd := len(g[v])
	revOfRev := len(g[u])
	g[u] = append(g[u], edge{v, cap, cap, revOfFwd})
	g[v] = append(g[v], edge{u, 0, 0, revOfRev})
}

// buildFlowGraph constructs a node-split residual graph from c.
// Each room i becomes two nodes: in=2i and out=2i+1, connected by a capacity-1 edge
// (capacity c.Ants for start and end) to enforce the one-ant-per-room constraint.
// Returns the graph g, source node src (Start_in), sink node snk (End_out),
// and names, a slice mapping nodeID/2 to room name.
func buildFlowGraph(c *graph.Colony) (g [][]edge, src, snk int, names []string) {
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
	for _, name := range roomNames {
		i := idx[name]
		in, out := i*2, i*2+1
		cap := 1
		if name == c.Start || name == c.End {
			cap = c.Ants
		}
		addEdge(g, in, out, cap)
	}

	for _, a := range roomNames {
		neighbours := make([]string, len(c.Links[a]))
		copy(neighbours, c.Links[a])
		sort.Strings(neighbours)
		for _, b := range neighbours {
			if idx[a] < idx[b] {
				addEdge(g, idx[a]*2+1, idx[b]*2, 1)
				addEdge(g, idx[b]*2+1, idx[a]*2, 1)
			}
		}
	}

	src = idx[c.Start] * 2
	snk = idx[c.End]*2 + 1
	return
}

// bfs finds one shortest augmenting path from src to snk in residual graph g.
// Returns a parent array where parent[node] encodes the incoming edge as
// edgeIndex<<20 | fromNode, or nil if snk is unreachable.
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
				parent[e.to] = i<<20 | u
				if e.to == snk {
					return parent
				}
				queue = append(queue, e.to)
			}
		}
	}
	return nil
}

// maxFlow runs Edmonds-Karp on g until no augmenting path from src to snk exists,
// saturating all flow. Modifies g in place.
func maxFlow(g [][]edge, src, snk int) {
	for {
		parent := bfs(g, src, snk)
		if parent == nil {
			break
		}
		cur := snk
		for cur != src {
			encoded := parent[cur]
			from := encoded & ((1 << 20) - 1)
			ei := encoded >> 20
			g[from][ei].cap--
			g[cur][g[from][ei].rev].cap++
			cur = from
		}
	}
}

// extractPaths recovers all vertex-disjoint paths from src to snk in g by
// repeatedly calling tracePath until no flow-carrying path remains.
// names maps nodeID/2 to room name and is used to convert node IDs to graph.Path.
// Returns the slice of paths; each path is a []string from Start to End inclusive.
func extractPaths(g [][]edge, src, snk int, names []string) []graph.Path {
	var paths []graph.Path
	for {
		path := tracePath(g, src, snk)
		if path == nil {
			break
		}
		var p graph.Path
		for _, node := range path {
			if node%2 == 1 {
				p = append(p, names[node/2])
			}
		}
		paths = append(paths, p)
	}
	return paths
}

// tracePath finds one path from src to snk in g by following edges that carry flow
// (origCap > 0 && cap < origCap), then cancels that flow so subsequent calls
// yield different paths. Returns the node sequence src…snk, or nil if none exists.
func tracePath(g [][]edge, src, snk int) []int {
	visited := make([]bool, len(g))
	stack := []int{src}
	parent := make([]int, len(g))
	parentEdge := make([]int, len(g))
	for i := range parent {
		parent[i] = -1
	}
	parent[src] = -2
	visited[src] = true

	found := false
	for len(stack) > 0 && !found {
		u := stack[len(stack)-1]
		advanced := false
		for i, e := range g[u] {
			if e.origCap > 0 && e.cap < e.origCap && !visited[e.to] {
				parent[e.to] = u
				parentEdge[e.to] = i
				visited[e.to] = true
				if e.to == snk {
					found = true
					break
				}
				stack = append(stack, e.to)
				advanced = true
				break
			}
		}
		if !advanced && !found {
			stack = stack[:len(stack)-1]
		}
	}

	if !found {
		return nil
	}

	var nodes []int
	cur := snk
	for cur != src {
		nodes = append(nodes, cur)
		from := parent[cur]
		ei := parentEdge[cur]
		g[from][ei].cap++
		g[cur][g[from][ei].rev].cap--
		cur = from
	}
	nodes = append(nodes, src)
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	return nodes
}

// computeTurns returns the minimum number of turns required to move nAnts ants
// across paths using greedy assignment: each ant is placed on the path that
// minimises its personal finish turn (path length - 1 + current load on that path).
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

// FindPaths returns the optimal subset of vertex-disjoint paths through c that
// minimises the total number of simulation turns for c.Ants ants.
// Paths are sorted shortest-first. Returns an error if no path exists from Start to End.
func FindPaths(c *graph.Colony) ([]graph.Path, error) {
	g, src, snk, names := buildFlowGraph(c)
	maxFlow(g, src, snk)
	paths := extractPaths(g, src, snk, names)
	if len(paths) == 0 {
		return nil, errors.New("no path between start and end")
	}

	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j])
		}
		return strings.Join(paths[i], ",") < strings.Join(paths[j], ",")
	})

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
