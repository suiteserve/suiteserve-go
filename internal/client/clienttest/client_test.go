package clienttest

import (
	"golang.org/x/sync/errgroup"
	"os/exec"
	"testing"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	testCmd := exec.Command("go", "test", "-json", "./testdata")
	clientCmd := exec.Command("go", "run",
		"github.com/suiteserve/go-runner/cmd/suiteservego",
		"https://localhost:8080")
	testOut, err := testCmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	clientCmd.Stdin = testOut
	var eg errgroup.Group
	eg.Go(clientCmd.Run)
	eg.Go(testCmd.Run)
	if err := eg.Wait(); err != nil && err.(*exec.ExitError).ExitCode() != 1 {
		t.Logf("run cmds: %v", err)
	}
}
