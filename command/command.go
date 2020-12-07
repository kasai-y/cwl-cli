package command

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type Command struct {
	session *session.Session
}

type Config struct {
	Session *session.Session
}

func New(c *Config) *Command {
	return &Command{
		session: c.Session,
	}
}
