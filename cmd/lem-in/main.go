// Package main is the lem-in command. It reads a colony file, finds the optimal
// paths, simulates ant movement, and writes the result to stdout.
package main

import (
	"fmt"
	"os"

	"lemin/internal/output"
	"lemin/internal/parser"
	"lemin/internal/pathfinder"
	"lemin/internal/simulator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "ERROR: invalid data format, no input file\n")
		os.Exit(1)
	}

	colony, err := parser.Parse(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid data format, %s\n", err)
		os.Exit(1)
	}

	paths, err := pathfinder.FindPaths(colony)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid data format, %s\n", err)
		os.Exit(1)
	}

	turns := simulator.Simulate(colony, paths)
	output.Print(os.Stdout, colony, turns)
}
