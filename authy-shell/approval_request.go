package main

import (
	"encoding/json"
	//"github.com/dcu/onetouch-ssh"
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/dcu/go-authy"
	"github.com/mgutz/ansi"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	// StatusPending is set when the request is pending.
	StatusPending = "pending"

	// StatusApproved is set when the request is approved.
	StatusApproved = "approved"

	// StatusDenied is set when the request is denied.
	StatusDenied = "denied"

	// StatusFailed is set when the request is failed.
	StatusFailed = "failed"
)

// ApprovalRequest is the approval request response.
type ApprovalRequest struct {
	Status   string `json:"status"`
	UUID     string `json:"uuid"`
	Notified bool   `json:"notified"`

	api     *authy.Authy
	bar     *pb.ProgressBar
	showBar bool
}

func buildParams() url.Values {
	hostname := runCommand("hostname")
	sshConnection := strings.Split(os.Getenv("SSH_CONNECTION"), " ")
	clientIP := ""
	serverIP := ""

	if len(sshConnection) > 1 {
		clientIP = formatIPAndLocation(sshConnection[0])
	}

	if len(sshConnection) > 2 {
		serverIP = formatIPAndLocation(sshConnection[2])
	}

	params := url.Values{
		"details[Type]":      {"SSH Server"},
		"details[Server IP]": {serverIP},
		"details[User IP]":   {clientIP},
		"details[User]":      {os.Getenv("USER")},
		"logos[][res]":       {"default"},
		"logos[][url]":       {"http://authy-assets-dev.s3.amazonaws.com/authenticator/ipad/logo/high/liberty_bank@2x.png"},
	}
	if command := os.Getenv("SSH_ORIGINAL_COMMAND"); command != "" {
		typ, repo := parseGitCommand(command)
		if typ != "" {
			params.Add("message", fmt.Sprintf("git %s on %s", typ, hostname))
			params.Add("details[Repository]", repo)
		} else {
			params.Add("message", fmt.Sprintf("You are executing command on %s", hostname))
			params.Add("details[Command]", command)
		}
	} else {
		params.Add("message", fmt.Sprintf("You are login to %s", hostname))
	}

	return params
}

// NewApprovalRequest creates a new approval request.
func NewApprovalRequest(api *authy.Authy, authyID int) (*ApprovalRequest, error) {
	params := buildParams()
	path := fmt.Sprintf(`/onetouch/json/users/%d/approval_requests`, authyID)

	response, err := api.DoRequest("POST", path, params)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	jsonResponse := struct {
		Success         bool             `json:"success"`
		ApprovalRequest *ApprovalRequest `json:"approval_request"`
	}{}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return nil, err
	}

	approvalRequest := jsonResponse.ApprovalRequest
	approvalRequest.api = api
	approvalRequest.showBar = isInteractiveConnection()
	return approvalRequest, nil
}

// CheckStatus returns the status of the request.
func (approvalRequest *ApprovalRequest) CheckStatus(timeout time.Duration) string {
	showBar := isInteractiveConnection()
	timeWaited := 0 * time.Second
	interval := 2 * time.Second

	status := StatusPending
	approvalRequest.bar = pb.New(int(timeout / interval))
	approvalRequest.setStatus(status)

	if showBar {
		approvalRequest.bar.Start()
	}

	for timeWaited < timeout {
		status = approvalRequest.requestStatus()
		approvalRequest.setStatus(status)
		if status != StatusPending {
			break
		}

		time.Sleep(interval)
		timeWaited += interval

		if showBar {
			approvalRequest.bar.Increment()
		}
	}

	if showBar {
		approvalRequest.bar.FinishPrint("")
	}

	return status
}

func (approvalRequest *ApprovalRequest) setStatus(status string) {
	if approvalRequest.showBar {
		approvalRequest.bar.Prefix(formatForStatus(status))
	}
}

func formatForStatus(status string) string {
	switch status {
	case StatusPending:
		{
			status = ansi.Color(status, "yellow+h")
		}
	case StatusApproved:
		{
			status = ansi.Color(status, "green+h")
		}
	case StatusFailed:
		{
			status = ansi.Color(status, "red+h")
		}
	case StatusDenied:
		{
			status = ansi.Color(status, "red+h")
		}
	}
	return fmt.Sprintf("[%s] ", status)
}

func (approvalRequest *ApprovalRequest) requestStatus() string {
	response, err := approvalRequest.api.DoRequest("GET", fmt.Sprintf("/onetouch/json/approval_requests/%s", approvalRequest.UUID), url.Values{})

	if err != nil {
		return StatusFailed
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return StatusFailed
	}

	jsonResponse := struct {
		Success         bool             `json:"success"`
		ApprovalRequest *ApprovalRequest `json:"approval_request"`
	}{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return StatusFailed
	}

	status := jsonResponse.ApprovalRequest.Status
	return status
}
