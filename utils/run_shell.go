package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syscall"
)

var (
	gitCmdRx = regexp.MustCompile(`^(git-receive-pack|git-upload-pack) '(.*)'$`)
)

// RunShell runs the shell indicated in the $SHELL env var
func RunShell() {
	shell := os.Getenv("SHELL")

	var err error
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shell)
		err = detachCommand(shell, shellName, "-c", sshCommand)
	} else if shell != "" {
		err = detachCommand(shell)
	}

	if err != nil {
		panic(err)
	}
}

// RunShellFromPath runs a given shell
func RunShellFromPath(shellPath string) {
	var err error
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shellPath)
		err = detachCommand(shellPath, shellName, "-c", sshCommand)
	} else {
		err = detachCommand(shellPath)
	}

	if err != nil {
		panic(err)
	}
}

// RunShellFromPathWithArgs runs a shell given the given arguments
func RunShellFromPathWithArgs(shellPath string, shellArgs []string) {
	var err error
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shellPath)
		err = detachCommand(shellPath, shellName, "-c", sshCommand)
	} else {
		shellCommand := make([]string, 0)
		shellCommand = append(shellCommand, shellPath)
		shellCommand = append(shellCommand, shellArgs...)
		err = detachCommand(shellCommand...)
	}

	if err != nil {
		panic(err)
	}
}
func waitPid(pid int) error {
	for {
		var ws syscall.WaitStatus
		_, err := syscall.Wait4(pid, &ws, 0, nil)

		if err != nil {
			return err
		}

		if ws.Exited() {
			break
		}
	}

	return nil
}

func closePorts() {
	_ = os.Stdin.Close()
	_ = os.Stdout.Close()
	_ = os.Stderr.Close()
}

func detachCommand(command ...string) error {
	sys := syscall.SysProcAttr{}
	files := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	attr := syscall.ProcAttr{
		Env:   os.Environ(),
		Files: files,
		Sys:   &sys,
	}

	pid, err := syscall.ForkExec(command[0], command[1:], &attr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return err
	}

	closePorts()
	return waitPid(pid)
}

// RunCommand runns the given command and arguments
func RunCommand(command ...string) string {
	output, err := exec.Command(command[0], command[1:]...).Output()
	if err != nil {
		return err.Error()
	}

	return string(output)
}
