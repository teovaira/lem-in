// Package parser reads and validates lem-in colony files.
package parser

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"lemin/internal/graph"
)

// Parse reads and fully validates the colony file at the given path.
// Returns a populated *Colony or a descriptive error.
func Parse(filename string) (*graph.Colony, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("no input file")
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, strings.TrimRight(scanner.Text(), "\r"))
	}

	c := &graph.Colony{
		Rooms: make(map[string]*graph.Room),
		Links: make(map[string][]string),
	}

	pos := 0
	for pos < len(lines) {
		line := lines[pos]
		if line == "" {
			pos++
			continue
		}
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			c.RawLines = append(c.RawLines, line)
			pos++
			continue
		}
		ants, err := strconv.Atoi(line)
		if err != nil || ants <= 0 {
			return nil, errors.New("invalid number of ants")
		}
		c.Ants = ants
		c.RawLines = append(c.RawLines, line)
		pos++
		break
	}

	if c.Ants == 0 {
		return nil, errors.New("invalid number of ants")
	}

	nextIsStart := false
	nextIsEnd := false
	startSeen := false
	endSeen := false

	for ; pos < len(lines); pos++ {
		line := lines[pos]

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "##") {
			c.RawLines = append(c.RawLines, line)
			switch line {
			case "##start":
				if startSeen {
					return nil, errors.New("multiple start rooms")
				}
				nextIsStart = true
				startSeen = true
			case "##end":
				if endSeen {
					return nil, errors.New("multiple end rooms")
				}
				nextIsEnd = true
				endSeen = true
			}
			continue
		}

		if strings.HasPrefix(line, "#") {
			c.RawLines = append(c.RawLines, line)
			continue
		}

		if looksLikeRoom(line) {
			room, err := parseRoom(line, c)
			if err != nil {
				return nil, err
			}
			if nextIsStart {
				c.Start = room.Name
				nextIsStart = false
			}
			if nextIsEnd {
				c.End = room.Name
				nextIsEnd = false
			}
			c.RawLines = append(c.RawLines, line)
			continue
		}

		if isLinkLine(line) {
			if nextIsStart {
				return nil, errors.New("no start room found")
			}
			if nextIsEnd {
				return nil, errors.New("no end room found")
			}
			if err := parseLink(line, c); err != nil {
				return nil, err
			}
			c.RawLines = append(c.RawLines, line)
			continue
		}

		return nil, errors.New("invalid data format")
	}

	if c.Start == "" {
		return nil, errors.New("no start room found")
	}
	if c.End == "" {
		return nil, errors.New("no end room found")
	}

	return c, nil
}

// looksLikeRoom returns true if the line has exactly 3 space-separated fields.
// Coordinate validation is left to parseRoom so bad coords produce a proper error.
func looksLikeRoom(line string) bool {
	return len(strings.Fields(line)) == 3
}

func isLinkLine(line string) bool {
	if strings.Contains(line, " ") {
		return false
	}
	parts := strings.SplitN(line, "-", 2)
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}

func parseRoom(line string, c *graph.Colony) (*graph.Room, error) {
	fields := strings.Fields(line)
	x, err1 := strconv.Atoi(fields[1])
	y, err2 := strconv.Atoi(fields[2])
	if err1 != nil || err2 != nil {
		return nil, errors.New("invalid room coordinates")
	}
	name := fields[0]
	if strings.HasPrefix(name, "L") {
		return nil, errors.New("invalid room name")
	}
	if strings.HasPrefix(name, "#") {
		return nil, errors.New("invalid room name")
	}
	if _, exists := c.Rooms[name]; exists {
		return nil, errors.New("duplicate room")
	}
	room := &graph.Room{Name: name, X: x, Y: y}
	c.Rooms[name] = room
	c.Links[name] = []string{}
	return room, nil
}

func parseLink(line string, c *graph.Colony) error {
	parts := strings.SplitN(line, "-", 2)
	a, b := parts[0], parts[1]
	if a == b {
		return errors.New("room links to itself")
	}
	if _, ok := c.Rooms[a]; !ok {
		return errors.New("unknown room in link")
	}
	if _, ok := c.Rooms[b]; !ok {
		return errors.New("unknown room in link")
	}
	if c.HasLink(a, b) {
		return errors.New("duplicate link")
	}
	c.Links[a] = append(c.Links[a], b)
	c.Links[b] = append(c.Links[b], a)
	return nil
}
