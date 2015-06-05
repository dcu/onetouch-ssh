package main

import (
	"bufio"
	"fmt"
	"github.com/authy/onetouch-ssh"
	"github.com/dcu/go-authy"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	cleanCodeRegexp = regexp.MustCompile(`[^\d+]`)
)

// Verification allows to verify a TOTP code.
type Verification struct {
	authyID         int
	api             *authy.Authy
	approvalRequest *ApprovalRequest
}

// NewVerification builds a new TOTP verification.
func NewVerification(authyID string) *Verification {
	authyIDInt, _ := strconv.Atoi(authyID)
	config := ssh.NewConfig()
	api := authy.NewAuthyApi(config.AuthyAPIKey())
	api.ApiUrl = "https://staging-2.authy.com"

	return &Verification{
		authyID: authyIDInt,
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
	printMessage("[sent]\n")

	status := approvalRequest.CheckStatus(30 * time.Second)
	if status == StatusApproved {
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
	verification.api.RequestSms(verification.authyID, url.Values{})

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
