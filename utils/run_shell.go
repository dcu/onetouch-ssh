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

func RunShell() {
	shell := os.Getenv("SHELL")

	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shell)
		detachCommand(shell, shellName, "-c", sshCommand)
	} else if shell != "" {
		detachCommand(shell)
	}
}

func RunShellFromPath(shellPath string) {
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shellPath)
		detachCommand(shellPath, shellName, "-c", sshCommand)
	} else {
		detachCommand(shellPath)
	}
}

func RunShellFromPathWithArgs(shellPath string, shellArgs []string) {
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		shellName := filepath.Base(shellPath)
		detachCommand(shellPath, shellName, "-c", sshCommand)
	} else {
		shellCommand := make([]string, 0)
		shellCommand = append(shellCommand, shellPath)
		shellCommand = append(shellCommand, shellArgs...)
		detachCommand(shellCommand...)
	}
}
func waitPid(pid int) {
	for {
		var ws syscall.WaitStatus
		_, err := syscall.Wait4(pid, &ws, 0, nil)

		if err != nil {
			panic(err)
		}

		if ws.Exited() {
			break
		}
	}
}

func closePorts() {
	os.Stdin.Close()
	os.Stdout.Close()
	os.Stderr.Close()
}

func detachCommand(command ...string) {
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
		return
	}

	closePorts()
	waitPid(pid)
}

func RunCommand(command ...string) string {
	output, err := exec.Command(command[0], command[1:]...).Output()
	if err != nil {
		return err.Error()
	}

	return string(output)
}
