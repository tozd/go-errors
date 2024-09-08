package errors_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/tozd/go/errors"
)

func TestPanic(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	output := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, "go", "run", "-race", "testdata/panic.go")
	cmd.Stdout = output
	cmd.Stderr = output
	// We have to make a process group and send signals to the whole group.
	// See: https://github.com/golang/go/issues/40467
	cmd.SysProcAttr = &syscall.SysProcAttr{ //nolint:exhaustruct
		Setpgid: true,
	}
	cmd.Cancel = func() error {
		if cmd.Process.Pid < 1 {
			return nil
		}
		// We kill whole process group.
		e := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(e, syscall.ESRCH) {
			return os.ErrProcessDone
		}
		return e //nolint:wrapcheck
	}

	err := cmd.Start()
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	err = syscall.Kill(-cmd.Process.Pid, syscall.SIGINT)
	assert.NoError(t, err) //nolint:testifylint

	err = cmd.Wait()
	var exitError *exec.ExitError
	// TODO: Remove workaround.
	//       Currently "go run" does not return zero exit code when we send INT signal
	//       to the whole process group even if the child process exits with zero exit code.
	//       See: https://github.com/golang/go/issues/40467
	if errors.As(err, &exitError) && exitError.ExitCode() > 0 {
		assert.Equal(t, 1, exitError.ExitCode())
	} else {
		assert.NoError(t, err) //nolint:testifylint
	}

	assert.Regexp(t, `^panic: panic error\n`+
		`\t?key=value\n`+
		`\t?stack trace \(most recent call first\):\n`+
		`\t?main\.main\n`+
		`\t?\t.*/testdata/panic.go:8\n`+
		`\t?runtime\.main\n`+
		`\t?\t.*/src/runtime/proc.go:\d+\n`+
		`\t?runtime\.goexit\n`+
		`\t?\t.*/src/runtime/.*:\d+\n`+
		`\t?\n\n`+
		`goroutine 1 \[running\]:\n`+
		`main\.main\(\)\n`+
		`\t.*/testdata/panic.go:8 \+0xcf\n`+
		`exit status 2\n$`, output.String())
}
