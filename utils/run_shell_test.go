package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetachCommand(t *testing.T) {
	c := require.New(t)

	err := detachCommand("/bin/echo", "-n", "echo command")
	c.Nil(err)
}
