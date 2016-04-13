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
	"github.com/spf13/cobra"
)

// enableCmd represents the enable command
var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enables all users to use onetouch when using SSH connections",
	Long:  `The enable command adds all the users to "~/.ssh/authorized_keys" which is used to validate the SSH connection.`,
	Run: func(cmd *cobra.Command, args []string) {
		manager := ssh.NewAuthorizedKeysManager()
		err := manager.WriteToDefaultLocation()
		if err != nil {
			fmt.Printf("Error writing authorized keys: %s\n", err)
			os.Exit(1)
		}
		fmt.Printf("Wrote public keys to ssh authorized keys.\n")
	},
}

func init() {
	RootCmd.AddCommand(enableCmd)
}
