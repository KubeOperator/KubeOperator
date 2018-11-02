package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

// A PlaybookTest executes a given Ansible playbook and checks the exit code and
// output contents.
type PlaybookTest struct {
	// inputs
	Path string
	// expected outputs
	ExitCode int
	Output   []string // zero or more strings that should be in the output
}

// Run runs the PlaybookTest.
func (p PlaybookTest) Run(t *testing.T) {
	// A PlaybookTest is intended to be run in parallel with other tests.
	t.Parallel()

	cmd := exec.Command("ansible-playbook", "-e", "testing_skip_some_requirements=1", "-i", "/dev/null", p.Path)
	cmd.Env = append(os.Environ(), "ANSIBLE_FORCE_COLOR=1")
	b, err := cmd.CombinedOutput()

	// Check exit code.
	if (err == nil) && (p.ExitCode != 0) {
		p.checkExitCode(t, 0, p.ExitCode, cmd, b)
	}
	if (err != nil) && (p.ExitCode == 0) {
		got, ok := getExitCode(err)
		if !ok {
			t.Logf("unexpected error (%T): %[1]v", err)
			p.logCmdAndOutput(t, cmd, b)
			t.FailNow()
		}
		p.checkExitCode(t, got, p.ExitCode, cmd, b)
	}

	// Check output contents.
	var missing []string
	for _, s := range p.Output {
		if !bytes.Contains(b, []byte(s)) {
			missing = append(missing, s)
		}
	}
	if len(missing) > 0 {
		t.Logf("missing in output: %q", missing)
		p.logCmdAndOutput(t, cmd, b)
		t.FailNow()
	}
}

// getExitCode returns an exit code and true if the exit code could be taken
// from err, false otherwise.
// The implementation is GOOS-specific, and currently only supports Linux.
func getExitCode(err error) (int, bool) {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return -1, false
	}
	waitStatus, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return -1, false
	}
	return waitStatus.ExitStatus(), true
}

// checkExitCode marks the test as failed when got is different than want.
func (p PlaybookTest) checkExitCode(t *testing.T, got, want int, cmd *exec.Cmd, output []byte) {
	if got == want {
		return
	}
	t.Logf("got exit code %v, want %v", got, want)
	p.logCmdAndOutput(t, cmd, output)
	t.FailNow()
}

// logCmdAndOutput logs how to re-run a command and a summary of the output of
// its last execution for debugging.
func (p PlaybookTest) logCmdAndOutput(t *testing.T, cmd *exec.Cmd, output []byte) {
	const maxLines = 10
	lines := bytes.Split(bytes.TrimRight(output, "\n"), []byte("\n"))
	if len(lines) > maxLines {
		lines = append([][]byte{[]byte("...")}, lines[len(lines)-maxLines:len(lines)]...)
	}
	output = bytes.Join(lines, []byte("\n"))
	dir, err := filepath.Abs(cmd.Dir)
	if err != nil {
		panic(err)
	}
	t.Logf("\n$ (cd %s && %s)\n%s", dir, strings.Join(cmd.Args, " "), output)
}
