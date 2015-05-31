package ssh

import (
	"testing"
)

func Test_PutAndGet(t *testing.T) {
	user := &User{
		Username: "david",
	}
	config := NewConfig("/tmp/db")

	err := config.Put("key1", user)
	if err != nil {
		panic(err)
	}

	user2 := User{}
	err = config.Get("key1", &user2)
	if err != nil {
		panic(err)
	}

	if user.Username != user2.Username {
		t.Error("Username doesn't match.")
	}
}
