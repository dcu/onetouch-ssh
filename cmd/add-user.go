// Copyright © 2016 David Cuadrado
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dcu/onetouch-ssh/ssh"
	"github.com/spf13/cobra"
)

// add-userCmd represents the add-user command
var addUserCmd = &cobra.Command{
	Use:   "add-user <email> <country code> <phone number> <public ssh key>",
	Short: "Add user to be authenticated with 2FA for the current unix account",
	Long: `Links the current user account with an Authy ID and SSH public key.

Example:
    onetouch-ssh add-user user@example.org 1 5555555555 "ssh-rsa ... user@example-host"

Note: Newly added users will not be enabled automatically -- they will be enabled upon
      running 'onetouch-ssh enable'.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 4 {
			cmd.Help()
			os.Exit(1)
		}

		countryCode, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Printf("Country code %s is invalid: %s", args[1], err)
			os.Exit(1)
		}

		usersManager := ssh.NewUsersManager()
		publicKey := strings.Join(args[3:], " ")
		err = usersManager.AddUser(args[0], countryCode, args[2], publicKey)
		if err != nil {
			fmt.Printf("Error adding user: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addUserCmd)
}
