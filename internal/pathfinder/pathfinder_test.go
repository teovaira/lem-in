package pathfinder

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"lemin/internal/graph"
)

func loadColony(path string) (*graph.Colony, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &graph.Colony{
		Rooms: make(map[string]*graph.Room),
		Links: make(map[string][]string),
	}
	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return nil, os.ErrInvalid
	}
	ants, err := strconv.Atoi(lines[0])
	if err != nil || ants <= 0 {
		return nil, os.ErrInvalid
	}
	c.Ants = ants

	nextStart, nextEnd := false, false
	for _, line := range lines[1:] {
		switch {
		case line == "##start":
			nextStart = true
		case line == "##end":
			nextEnd = true
		case strings.HasPrefix(line, "#"):
		case strings.Contains(line, "-") && !strings.Contains(line, " "):
			parts := strings.SplitN(line, "-", 2)
			a, b := parts[0], parts[1]
			c.Links[a] = append(c.Links[a], b)
			c.Links[b] = append(c.Links[b], a)
		default:
			fields := strings.Fields(line)
			if len(fields) == 3 {
				x, _ := strconv.Atoi(fields[1])
				y, _ := strconv.Atoi(fields[2])
				r := &graph.Room{Name: fields[0], X: x, Y: y}
				c.Rooms[fields[0]] = r
				if _, ok := c.Links[fields[0]]; !ok {
					c.Links[fields[0]] = []string{}
				}
				if nextStart {
					c.Start = fields[0]
					nextStart = false
				} else if nextEnd {
					c.End = fields[0]
					nextEnd = false
				}
			}
		}
	}
	return c, nil
}

func makeColony(ants int, start, end string, rooms [][3]string, links [][2]string) *graph.Colony {
	c := &graph.Colony{
		Ants:  ants,
		Start: start,
		End:   end,
		Rooms: make(map[string]*graph.Room),
		Links: make(map[string][]string),
	}
	for _, r := range rooms {
		x, _ := strconv.Atoi(r[1])
		y, _ := strconv.Atoi(r[2])
		c.Rooms[r[0]] = &graph.Room{Name: r[0], X: x, Y: y}
		c.Links[r[0]] = []string{}
	}
	for _, l := range links {
		a, b := l[0], l[1]
		c.Links[a] = append(c.Links[a], b)
		c.Links[b] = append(c.Links[b], a)
	}
	return c
}

func TestFindPathsSinglePath(t *testing.T) {
	c := makeColony(1, "start", "end",
		[][3]string{{"start", "0", "0"}, {"a", "1", "0"}, {"b", "2", "0"}, {"end", "3", "0"}},
		[][2]string{{"start", "a"}, {"a", "b"}, {"b", "end"}},
	)
	paths, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 {
		t.Fatalf("want 1 path, got %d", len(paths))
	}
	want := []string{"start", "a", "b", "end"}
	for i, r := range want {
		if paths[0][i] != r {
			t.Errorf("path[0][%d]: want %q got %q", i, r, paths[0][i])
		}
	}
}

func TestFindPathsTwoPaths(t *testing.T) {
	c := makeColony(2, "start", "end",
		[][3]string{{"start", "0", "0"}, {"a", "1", "1"}, {"b", "1", "-1"}, {"end", "2", "0"}},
		[][2]string{{"start", "a"}, {"a", "end"}, {"start", "b"}, {"b", "end"}},
	)
	paths, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("want 2 paths, got %d", len(paths))
	}
	seen := map[string]bool{}
	for _, p := range paths {
		for _, r := range p[1 : len(p)-1] {
			if seen[r] {
				t.Errorf("intermediate room %q shared between paths", r)
			}
			seen[r] = true
		}
	}
}

func TestFindPathsNoPath(t *testing.T) {
	c := makeColony(1, "start", "end",
		[][3]string{{"start", "0", "0"}, {"end", "1", "0"}},
		[][2]string{},
	)
	_, err := FindPaths(c)
	if err == nil {
		t.Fatal("expected error for disconnected graph")
	}
}

func TestFindPathsCycle(t *testing.T) {
	c := makeColony(1, "start", "end",
		[][3]string{{"start", "0", "0"}, {"a", "1", "0"}, {"b", "2", "0"}, {"end", "3", "0"}},
		[][2]string{{"start", "a"}, {"a", "b"}, {"b", "start"}, {"b", "end"}},
	)
	done := make(chan struct{})
	go func() {
		FindPaths(c) //nolint
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("FindPaths did not terminate (possible infinite loop)")
	}
}

func TestFindPathsDirectLink(t *testing.T) {
	c := makeColony(1, "start", "end",
		[][3]string{{"start", "0", "0"}, {"end", "1", "0"}},
		[][2]string{{"start", "end"}},
	)
	paths, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 1 || len(paths[0]) != 2 {
		t.Fatalf("want 1 path of length 2, got %v", paths)
	}
}

func TestFindPathsOptimalSubset(t *testing.T) {
	// 3 paths: A=[s,a,e] (len=3), B=[s,b,c,e] (len=4), C=[s,d,...,e] (len=12).
	// For 4 ants:
	//   1 path only: turns = (3-1) + (4-1) = 5
	//   2 paths A+B: ant1->A(2), ant2->B(3), ant3->A(3), ant4->B(4) -> 4 turns
	//   3 paths A+B+C: worse because C is long
	// So optimal is 2 paths, turns=4.
	c := makeColony(4, "s", "e",
		[][3]string{
			{"s", "0", "0"}, {"e", "9", "0"},
			{"a", "1", "1"},
			{"b", "1", "2"}, {"c", "2", "2"},
			{"d", "1", "3"}, {"f", "2", "3"}, {"g", "3", "3"}, {"h", "4", "3"},
			{"i", "5", "3"}, {"j", "6", "3"}, {"k", "7", "3"}, {"l", "8", "3"},
		},
		[][2]string{
			{"s", "a"}, {"a", "e"},
			{"s", "b"}, {"b", "c"}, {"c", "e"},
			{"s", "d"}, {"d", "f"}, {"f", "g"}, {"g", "h"}, {"h", "i"}, {"i", "j"}, {"j", "k"}, {"k", "l"}, {"l", "e"},
		},
	)
	paths, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Errorf("expected 2 paths for 4 ants, got %d", len(paths))
	}
	turns := computeTurns(paths, 4)
	if turns > 4 {
		t.Errorf("expected ≤4 turns, got %d", turns)
	}
}

func TestComputeTurns1Ant1Path(t *testing.T) {
	paths := []graph.Path{{"start", "a", "end"}}
	if got := computeTurns(paths, 1); got != 2 {
		t.Errorf("want 2, got %d", got)
	}
}

func TestComputeTurns3Ants2Paths(t *testing.T) {
	// path0 len=2 (1 step), path1 len=4 (3 steps)
	// ant1 -> path0, finish turn 1
	// ant2 -> path0, finish turn 2
	// ant3 -> path1, finish turn 3
	// max = 3
	paths := []graph.Path{
		{"s", "e"},
		{"s", "a", "b", "e"},
	}
	if got := computeTurns(paths, 3); got != 3 {
		t.Errorf("want 3, got %d", got)
	}
}

func testExampleTurns(t *testing.T, file string, maxTurns int) {
	t.Helper()
	c, err := loadColony(file)
	if err != nil {
		t.Fatalf("loadColony: %v", err)
	}
	paths, err := FindPaths(c)
	if err != nil {
		t.Fatalf("FindPaths: %v", err)
	}
	turns := computeTurns(paths, c.Ants)
	if turns > maxTurns {
		t.Errorf("%s: want ≤%d turns, got %d", file, maxTurns, turns)
	}
}

func TestFindPathsExample00(t *testing.T) { testExampleTurns(t, "../../testdata/example00.txt", 6) }
func TestFindPathsExample01(t *testing.T) { testExampleTurns(t, "../../testdata/example01.txt", 8) }
func TestFindPathsExample02(t *testing.T) { testExampleTurns(t, "../../testdata/example02.txt", 11) }
func TestFindPathsExample03(t *testing.T) { testExampleTurns(t, "../../testdata/example03.txt", 6) }
func TestFindPathsExample04(t *testing.T) { testExampleTurns(t, "../../testdata/example04.txt", 6) }
func TestFindPathsExample05(t *testing.T) { testExampleTurns(t, "../../testdata/example05.txt", 8) }

func TestFindPathsPerformance100(t *testing.T) {
	rooms := [][3]string{{"s", "0", "0"}, {"e", "99", "0"}}
	links := [][2]string{}
	for i := 0; i < 10; i++ {
		rooms = append(rooms, [3]string{"a" + strconv.Itoa(i), strconv.Itoa(i + 1), "1"})
	}
	for i := 0; i < 5; i++ {
		rooms = append(rooms, [3]string{"b" + strconv.Itoa(i), strconv.Itoa(i + 1), "2"})
	}
	for i := 0; i < 20; i++ {
		rooms = append(rooms, [3]string{"c" + strconv.Itoa(i), strconv.Itoa(i + 1), "3"})
	}
	links = append(links, [2]string{"s", "a0"})
	for i := 0; i < 9; i++ {
		links = append(links, [2]string{"a" + strconv.Itoa(i), "a" + strconv.Itoa(i+1)})
	}
	links = append(links, [2]string{"a9", "e"})
	links = append(links, [2]string{"s", "b0"})
	for i := 0; i < 4; i++ {
		links = append(links, [2]string{"b" + strconv.Itoa(i), "b" + strconv.Itoa(i+1)})
	}
	links = append(links, [2]string{"b4", "e"})
	links = append(links, [2]string{"s", "c0"})
	for i := 0; i < 19; i++ {
		links = append(links, [2]string{"c" + strconv.Itoa(i), "c" + strconv.Itoa(i+1)})
	}
	links = append(links, [2]string{"c19", "e"})

	c := makeColony(100, "s", "e", rooms, links)
	start := time.Now()
	_, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("took %v, want <5s", elapsed)
	}
}

func TestFindPathsPerformance1000(t *testing.T) {
	rooms := [][3]string{{"s", "0", "0"}, {"e", "99", "0"}}
	links := [][2]string{}
	pathLens := []int{5, 8, 12, 20, 30}
	for pi, plen := range pathLens {
		prefix := string(rune('a' + pi))
		prev := "s"
		for i := 0; i < plen; i++ {
			name := prefix + strconv.Itoa(i)
			rooms = append(rooms, [3]string{name, strconv.Itoa(i + 1), strconv.Itoa(pi + 1)})
			links = append(links, [2]string{prev, name})
			prev = name
		}
		links = append(links, [2]string{prev, "e"})
	}
	c := makeColony(1000, "s", "e", rooms, links)
	start := time.Now()
	_, err := FindPaths(c)
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 30*time.Second {
		t.Errorf("took %v, want <30s", elapsed)
	}
}
