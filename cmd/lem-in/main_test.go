package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmp, err := os.MkdirTemp("", "lemin_e2e_*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create temp dir:", err)
		os.Exit(1)
	}
	binaryPath = filepath.Join(tmp, "lem-in")
	out, err := exec.Command("go", "build", "-o", binaryPath, ".").CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "build failed: %s\n%s\n", err, out)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(tmp)
	os.Exit(code)
}

func run(args ...string) (stdout, stderr string, exitCode int) {
	cmd := exec.Command(binaryPath, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func countTurns(stdout string) int {
	lines := strings.Split(strings.TrimRight(stdout, "\n"), "\n")
	pastBlank := false
	count := 0
	for _, l := range lines {
		if pastBlank {
			if l != "" {
				count++
			}
		} else if l == "" {
			pastBlank = true
		}
	}
	return count
}

func TestE2EExample00(t *testing.T) {
	stdout, _, code := run("../../testdata/example00.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 6 {
		t.Errorf("want ≤6 turns, got %d", n)
	}
	if !strings.HasPrefix(stdout, "4\n") {
		t.Errorf("output must start with colony data (ant count)")
	}
}

func TestE2EExample01(t *testing.T) {
	stdout, _, code := run("../../testdata/example01.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 8 {
		t.Errorf("want ≤8 turns, got %d", n)
	}
}

func TestE2EExample02(t *testing.T) {
	stdout, _, code := run("../../testdata/example02.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 11 {
		t.Errorf("want ≤11 turns, got %d", n)
	}
}

func TestE2EExample03(t *testing.T) {
	stdout, _, code := run("../../testdata/example03.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 6 {
		t.Errorf("want ≤6 turns, got %d", n)
	}
}

func TestE2EExample04(t *testing.T) {
	stdout, _, code := run("../../testdata/example04.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 6 {
		t.Errorf("want ≤6 turns, got %d", n)
	}
}

func TestE2EExample05(t *testing.T) {
	stdout, _, code := run("../../testdata/example05.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if n := countTurns(stdout); n > 8 {
		t.Errorf("want ≤8 turns, got %d", n)
	}
}

func TestE2EExample06(t *testing.T) {
	stdout, _, code := run("../../testdata/example06.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if !strings.HasPrefix(stdout, "100\n") {
		t.Errorf("output must start with the ant count")
	}
}

func TestE2EExample07(t *testing.T) {
	stdout, _, code := run("../../testdata/example07.txt")
	if code != 0 {
		t.Fatalf("want exit 0, got %d", code)
	}
	if !strings.HasPrefix(stdout, "1000\n") {
		t.Errorf("output must start with the ant count")
	}
}

func TestE2EBadExample00(t *testing.T) {
	stdout, stderr, code := run("../../testdata/badexample00.txt")
	if code == 0 {
		t.Error("want non-zero exit code")
	}
	if stdout != "" {
		t.Errorf("want empty stdout on error, got: %q", stdout)
	}
	if !strings.Contains(stderr, "ERROR: invalid data format, invalid number of ants") {
		t.Errorf("want specific error in stderr, got: %q", stderr)
	}
}

func TestE2EBadExample01(t *testing.T) {
	stdout, stderr, code := run("../../testdata/badexample01.txt")
	if code == 0 {
		t.Error("want non-zero exit code")
	}
	if stdout != "" {
		t.Errorf("want empty stdout on error, got: %q", stdout)
	}
	if !strings.Contains(stderr, "ERROR: invalid data format, room links to itself") {
		t.Errorf("want specific error in stderr, got: %q", stderr)
	}
}

func TestE2ENoFile(t *testing.T) {
	stdout, stderr, code := run()
	if code == 0 {
		t.Error("want non-zero exit code")
	}
	if stdout != "" {
		t.Errorf("want empty stdout on error, got: %q", stdout)
	}
	if !strings.Contains(stderr, "ERROR: invalid data format, no input file") {
		t.Errorf("want specific error in stderr, got: %q", stderr)
	}
}

func TestE2EInvalidFile(t *testing.T) {
	stdout, stderr, code := run("../../testdata/nonexistent.txt")
	if code == 0 {
		t.Error("want non-zero exit code")
	}
	if stdout != "" {
		t.Errorf("want empty stdout on error, got: %q", stdout)
	}
	if !strings.Contains(stderr, "ERROR: invalid data format, no input file") {
		t.Errorf("want specific error in stderr, got: %q", stderr)
	}
}
