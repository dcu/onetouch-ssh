package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/dcu/go-authy"
	"github.com/dcu/onetouch-ssh"
	"github.com/mgutz/ansi"
)

var (
	cleanCodeRegexp = regexp.MustCompile(`[^\d+]`)
)

// Verification allows to verify a TOTP code.
type Verification struct {
	authyID         string
	api             *authy.Authy
	approvalRequest *ApprovalRequest
}

// NewVerification builds a new TOTP verification.
func NewVerification(authyID string) *Verification {
	config := ssh.NewConfig()
	api := authy.NewAuthyAPI(config.AuthyAPIKey())
	api.BaseURL = "https://api.authy.com"

	return &Verification{
		authyID: authyID,
		api:     api,
	}
}

func (verification *Verification) sendApprovalRequest() *ApprovalRequest {
	approvalRequest, err := NewApprovalRequest(verification.api, verification.authyID)
	if err != nil {
		panic(err)
	}

	return approvalRequest
}

func (verification *Verification) perform() {
	printMessage("Sending approval request to your device... ")
	approvalRequest := verification.sendApprovalRequest()
	printMessage(ansi.Color("[sent]\n", "green+h"))

	status := approvalRequest.CheckStatus(30 * time.Second)
	if status == StatusApproved {
		printMessage("You've been logged in successfully.")
		runShell()
	} else if status == StatusPending && isInteractiveConnection() {
		printMessage("You didn't confirm the request. ")
		printMessage("A text-message was sent to your phone.\n")
		code := verification.askTOTPCode()
		if verification.verifyCode(code) {
			runShell()
		}
	}
}

func (verification *Verification) askTOTPCode() string {
	verification.api.RequestSMS(verification.authyID, url.Values{})

	var code string
	for i := 0; i < 3; i++ {
		fmt.Printf("Enter security code: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			code = scanner.Text()
		}
		code = cleanCodeRegexp.ReplaceAllString(code, "")

		if code != "" {
			break
		}
	}
	return code
}

func (verification *Verification) verifyCode(code string) bool {
	result, err := verification.api.VerifyToken(verification.authyID, code, url.Values{})
	if err != nil {
		return false
	}

	if result.Valid() {
		return true
	}

	return false
}
