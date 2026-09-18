package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	if os.Getenv("RUN_CRASH_TEST") == "1" {
		main()
		// Should not reach here, ensure exit with 1 if it does
		os.Exit(1)
	}
	output := execExitTest(t, "TestMain", false)

	assert := assert.New(t)

	assert.Contains(output, "netrouse")
	assert.Contains(output, "Usage:")
	assert.Contains(output, "Available Commands:")
}

func execExitTest(t *testing.T, test string, exitsError bool) string {
	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(os.Environ(), "RUN_CRASH_TEST=1")
	out, err := cmd.Output()
	if exitsError && err == nil {
		t.Fatal("Process exited without error")
	} else if !exitsError && err == nil {
		return string(out)
	}
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return string(out)
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
	return ""
}
