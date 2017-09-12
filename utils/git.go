package utils

// ParseGitCommand rturns the type of the operation and the repo given a git command
func ParseGitCommand(command string) (typ string, repo string) {
	result := gitCmdRx.FindStringSubmatch(command)

	if len(result) == 3 {
		if result[1] == "git-receive-pack" {
			typ = "push"
		} else {
			typ = "fetch"
		}
		return typ, result[2]
	}

	return "", ""
}
