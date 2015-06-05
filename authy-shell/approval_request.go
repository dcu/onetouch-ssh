package main

import (
	"encoding/json"
	//"github.com/authy/onetouch-ssh"
	"fmt"
	"github.com/dcu/go-authy"
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
	Status string `json:"status"`
	UUID   string `json:"uuid"`

	api *authy.Authy
}

func buildParams() url.Values {
	sshConnection := strings.Split(os.Getenv("SSH_CONNECTION"), " ")
	clientIP := ""
	serverIP := ""
	if len(sshConnection) > 1 {
		clientIP = sshConnection[0]
		if clientIP == "::1" || serverIP == "127.0.0.1" {
			clientIP = "localhost"
		}
	}
	if len(sshConnection) > 3 {
		serverIP = sshConnection[2]
		if serverIP == "::1" || serverIP == "127.0.0.1" {
			serverIP = "localhost"
		}
	}
	hostname := runCommand("hostname")

	params := url.Values{
		"details[Type]":      {"SSH Server"},
		"details[Server IP]": {serverIP},
		"details[User IP]":   {clientIP},
		"details[User]":      {os.Getenv("USER")},
		"logos[][res]":       {"default"},
		"logos[][url]":       {"http://authy-assets-dev.s3.amazonaws.com/authenticator/ipad/logo/high/liberty_bank@2x.png"},
		"message":            {fmt.Sprintf("You are logging to %s", hostname)},
	}
	if command := os.Getenv("SSH_ORIGINAL_COMMAND"); command != "" {
		params.Add("details[Command]", command)
	}

	return params
}

// NewApprovalRequest creates a new approval request.
func NewApprovalRequest(api *authy.Authy, authyID int) (*ApprovalRequest, error) {
	params := buildParams()

	response, err := api.DoRequest("POST", fmt.Sprintf(`/onetouch/json/users/%d/approval_requests`, authyID), params)
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
	return approvalRequest, nil
}

// CheckStatus returns the status of the request.
func (approvalRequest *ApprovalRequest) CheckStatus(timeout time.Duration) string {
	timeWaited := 0 * time.Second
	interval := 2 * time.Second

	status := StatusPending
	for timeWaited < timeout {
		status = approvalRequest.requestStatus()
		if status != StatusPending {
			break
		}

		time.Sleep(interval)
		timeWaited += interval
	}

	return status
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
		Success         bool              `json:"success"`
		ApprovalRequest map[string]string `json:"approval_request"`
	}{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return StatusFailed
	}

	status := jsonResponse.ApprovalRequest["status"]
	return status
}
