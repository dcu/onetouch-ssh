// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/dcu/onetouch-ssh/ssh"
	"github.com/spf13/cobra"
)

// removeUserCmd represents the removeUser command
var removeUserCmd = &cobra.Command{
	Use:   "remove-user <authy-id>",
	Short: "Removes the user with the given authy id",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			os.Exit(1)
		}

		usersManager := ssh.NewUsersManager()
		err := usersManager.RemoveUser(args[0])
		if err != nil {
			fmt.Printf("Error removing user: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(removeUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeUserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeUserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
