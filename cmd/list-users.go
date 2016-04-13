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

	"github.com/dcu/onetouch-ssh/ssh"
	"github.com/dcu/onetouch-ssh/utils"
	"github.com/spf13/cobra"
)

// add-userCmd represents the add-user command
var listUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "List details about all users in the database.",
	Long: `Prints a list of users with details about the registration.

Example:
    onetouch-ssh list-users
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			os.Exit(1)
		}

		usersManager := ssh.NewUsersManager()
		users := usersManager.UserIDList()

		fmt.Printf("%16s\t%-24s\t %-48s\n", "Authy ID", "Key Comment", "Public Key")
		for i := 0; i < len(users); i++ {
			user := usersManager.GetUser(users[i])
			for j := 0; j < len(user.PublicKeys); j++ {
				fingerprint, comment := utils.PublicKeyFingerprint(user.PublicKeys[j])
				fmt.Printf(
					"%16s\t%-24s\t%48s\n",
					user.AuthyID,
					comment,
					fingerprint)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listUsersCmd)
}
