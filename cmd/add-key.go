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
	"strings"

	"github.com/dcu/onetouch-ssh/ssh"
	"github.com/spf13/cobra"
)

// add-userCmd represents the add-user command
var addKeyCmd = &cobra.Command{
	Use:   "add-key <authy id> <ssh public key>",
	Short: "Adds a public key to a user.",
	Long: `Links an additional public key to an Authy user id.

Example:
    onetouch-ssh add-key 12345678 "ssh-rsa ... user@example-host"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			os.Exit(1)
		}

		usersManager := ssh.NewUsersManager()
		publicKey := strings.Join(args[1:], " ")
		err := usersManager.AddKey(args[0], publicKey)
		if err != nil {
			fmt.Printf("Error adding key to user: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(addKeyCmd)
}
