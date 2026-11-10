//go:build !windows

package quota

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

// ReadCodex runs one short-lived `codex app-server` exchange: initialize,
// initialized, account/rateLimits/read. It starts no thread, turn or inference.
func ReadCodex(ctx context.Context, opts CodexOptions) (*Snapshot, error) {
	home := opts.Home
	if home == "" {
		home = homeDir()
	}
	executable, err := locateCodex(opts.Executable, os.Getenv("PATH"), codexSearchDirs(home))
	if err != nil {
		return nil, err
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = DefaultCodexTimeout
	}
	exchange, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.Command(executable, "app-server")
	// Own process group, so cleanup reaches only this server and its children.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = time.Second
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("%w: start failed", ErrCLINotFound)
	}
	defer func() {
		_ = stdin.Close()
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		_ = cmd.Wait()
	}()
	lines := make(chan []byte)
	go readCodexLines(stdout, lines, exchange.Done())

	await := func(id int) (json.RawMessage, error) {
		for {
			select {
			case <-exchange.Done():
				if ctx.Err() != nil {
					return nil, fmt.Errorf("codex quota read: %w", ctx.Err())
				}
				return nil, ErrTimeout
			case line, ok := <-lines:
				if !ok {
					return nil, fmt.Errorf("%w: codex stream ended", ErrTimeout)
				}
				var message struct {
					ID     *int            `json:"id"`
					Result json.RawMessage `json:"result"`
					Error  *rpcError       `json:"error"`
				}
				if json.Unmarshal(line, &message) != nil || message.ID == nil || *message.ID != id {
					continue
				}
				if message.Error != nil {
					return nil, classifyRPCError(message.Error)
				}
				return message.Result, nil
			}
		}
	}
	send := func(message string) error {
		if _, err := io.WriteString(stdin, message+"\n"); err != nil {
			return fmt.Errorf("%w: codex stdin closed", ErrTimeout)
		}
		return nil
	}

	if err := send(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"clientInfo":{"name":"ply","title":"Ply","version":"1"}}}`); err != nil {
		return nil, err
	}
	if _, err := await(1); err != nil {
		return nil, err
	}
	if err := send(`{"jsonrpc":"2.0","method":"initialized"}`); err != nil {
		return nil, err
	}
	if err := send(`{"jsonrpc":"2.0","id":2,"method":"account/rateLimits/read"}`); err != nil {
		return nil, err
	}
	result, err := await(2)
	if err != nil {
		return nil, err
	}
	return decodeCodexRateLimits(result, time.Now())
}

// readCodexLines skips lines longer than maxCodexLine and stops on EOF or done.
func readCodexLines(r io.Reader, lines chan<- []byte, done <-chan struct{}) {
	defer close(lines)
	reader := bufio.NewReaderSize(r, 64*1024)
	var line []byte
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if err != nil {
			return
		}
		if len(line)+len(chunk) <= maxCodexLine {
			line = append(line, chunk...)
		} else {
			line = line[:0:0]
			for isPrefix && err == nil {
				_, isPrefix, err = reader.ReadLine()
			}
			continue
		}
		if isPrefix {
			continue
		}
		select {
		case lines <- line:
		case <-done:
			return
		}
		line = nil
	}
}
