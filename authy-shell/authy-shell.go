package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/authy/onetouch-ssh"
	"github.com/dcu/go-authy"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"syscall"
	"time"
)

// ApprovalRequest is the approval request response.
type ApprovalRequest struct {
	Status string `json:"status"`
	UUID   string `json:"uuid"`
}

func main() {
	authyID := os.Args[1]
	sendApprovalRequest(authyID)
	runShell()
}

func sendApprovalRequest(authyID string) {
	params := url.Values{
		"details[type]": {"ssh"},
		"details[ip]":   {"fake"},
		"message":       {"You are logging to ssh server"},
	}

	config := ssh.NewConfig()
	api := authy.NewAuthyApi(config.AuthyAPIKey())
	api.ApiUrl = "https://staging-2.authy.com"

	response, err := api.DoRequest("POST", `/onetouch/json/users/`+authyID+`/approval_requests`, params)
	if err != nil {
		panic(err) // FIXME
	}

	if response.StatusCode != 200 {
		// Send SMS and ask code.
		authyIDInt, _ := strconv.Atoi(authyID)
		api.RequestSms(authyIDInt, url.Values{})

		fmt.Printf("Enter security code: ")
		scanner := bufio.NewScanner(os.Stdin)
		var code string
		if scanner.Scan() {
			code = scanner.Text()
		}
		println(code)
		runShell()
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err) // FIXME
	}

	jsonResponse := struct {
		Success         bool             `json:"success"`
		ApprovalRequest *ApprovalRequest `json:"approval_request"`
	}{}
	println(string(body))

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		panic(err) // FIXME
	}

	approvalRequest := jsonResponse.ApprovalRequest
	if approvalRequest.Status == "pending" {
		time.Sleep(10 * time.Second)
		runShell()
	} else {
		return
	}
}

func runShell() {
	if sshCommand := os.Getenv("SSH_ORIGINAL_COMMAND"); sshCommand != "" {
		detachCommand(sshCommand)
	} else if shell := os.Getenv("SHELL"); shell != "" {
		detachCommand(shell)
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
	println("Launching: " + command[0])

	sys := syscall.SysProcAttr{}
	files := []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}
	attr := syscall.ProcAttr{
		Env:   os.Environ(),
		Files: files,
		Sys:   &sys,
	}

	pid, _ := syscall.ForkExec(command[0], command[1:], &attr)
	closePorts()
	waitPid(pid)
}
