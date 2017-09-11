package ssh

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/dcu/go-authy"

	"github.com/dcu/onetouch-ssh/utils"
	"github.com/mgutz/ansi"
)

var (
	ApprovalTimeout        = 45 * time.Second
	cleanCodeRegexp        = regexp.MustCompile(`[^\d+]`)
	MaxAttemptsToReadCode  = 3
	ErrInvalidVerification = errors.New("invalid verification")
)

type Verification struct {
	api     *authy.Authy
	authyID string
}

func NewVerification(authyID string) *Verification {
	api, err := LoadAuthyAPI()
	if err != nil {
		return nil
	}

	return &Verification{
		authyID: authyID,
		api:     api,
	}
}

func (verification *Verification) Run() error {
	utils.PrintMessage("Sending approval request to your device... ")
	request, err := verification.SendOneTouchRequest()
	if err != nil {
		utils.PrintMessage(err.Error())
		return ErrInvalidVerification
	}

	utils.PrintMessage(ansi.Color("[sent]\n", "green+h"))
	status, err := verification.api.WaitForApprovalRequest(request.UUID, ApprovalTimeout, url.Values{})
	if err != nil {
		utils.PrintMessage(err.Error())
		return ErrInvalidVerification
	}

	if status == authy.OneTouchStatusApproved {
		runShell()
		return nil
	}

	if status == authy.OneTouchStatusExpired && utils.IsInteractiveConnection() {
		utils.PrintMessage("You didn't confirm the request. ")

		code := verification.askTOTPCode()
		if verification.verifyCode(code) {
			runShell()
			return nil
		}
	}

	return ErrInvalidVerification
}

func runShell() {
	config, err := LoadConfig()
	// Check if the config could be loaded to get preferred user shell.
	if err == nil {
		utils.PrintMessage("You've been logged in successfully.\n")
		if len(config.ShellPath) != 0 {
			if len(config.ShellArgs) != 0 {
				utils.RunShellFromPathWithArgs(
					config.ShellPath,
					config.ShellArgs)
			} else {
				utils.RunShellFromPath(config.ShellPath)
			}
		} else {
			utils.RunShell()
		}
	} else {
		utils.RunShell()
	}
}

func (verification *Verification) SendOneTouchRequest() (*authy.ApprovalRequest, error) {
	// Better error handling on RunCommand?
	hostname := os.Hostname()
	sshConnection := strings.Split(os.Getenv("SSH_CONNECTION"), " ")
	clientIP := ""
	serverIP := ""

	if len(sshConnection) > 1 {
		clientIP = utils.FormatIPAndLocation(sshConnection[0])
	}

	if len(sshConnection) > 2 {
		serverIP = utils.FormatIPAndLocation(sshConnection[2])
	}

	var message string

	details := authy.Details{
		"Type":      "SSH Server",
		"Server IP": serverIP,
		"User IP":   clientIP,
		"User":      os.Getenv("USER"),
	}

	if command := os.Getenv("SSH_ORIGINAL_COMMAND"); command != "" {
		typ, repo := utils.ParseGitCommand(command)
		if typ != "" {
			message = fmt.Sprintf("git %s on %s", typ, hostname)
			details["Repository"] = repo
		} else {
			message = fmt.Sprintf("Executing command on: %s", hostname)
			details["Command"] = command
		}
	} else {
		message = fmt.Sprintf("Logging in to: %s", hostname)
	}

	return verification.api.SendApprovalRequest(verification.authyID, message, details, url.Values{})
}

func (verification *Verification) sendTOTPCode() {
	_, err := verification.api.RequestSMS(verification.authyID, url.Values{})
	if err != nil {
		return
	}

	// TODO: check if the SMS request was ignored
	utils.PrintMessage("A text-message was sent to your phone.\n")
}

func (verification *Verification) askTOTPCode() string {
	var code string
	for i := 0; i < MaxAttemptsToReadCode; i++ {
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
