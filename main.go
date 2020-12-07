package main

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"

	"github.com/kasai-y/cwl-cli/command"
	"github.com/kasai-y/cwl-cli/log"
)

var version string
var profile string
var debug bool

var putLogFlag struct {
	LogGroupName string
	Prefix       string
	Message      string
}

func main() {

	app := cli.NewApp()
	app.Name = "cwl-cli"
	app.Usage = "CloudWatchLogs CLI."
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "profile",
			Usage:       "profile for AWS",
			EnvVar:      "AWS_PROFILE",
			Destination: &profile,
		},
		cli.BoolFlag{
			Name:        "debug",
			Hidden:      true,
			Destination: &debug,
		},
	}
	app.Commands = cli.Commands{
		{
			Name:   "put-log",
			Usage:  "Put Exclude Event",
			Action: action,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "log-group-name,name,n",
					Usage:       "LogGroupName",
					Required:    true,
					EnvVar:      "AWS_LOG_GROUP_NAME",
					Destination: &putLogFlag.LogGroupName,
				},
				cli.StringFlag{
					Name:        "prefix,p",
					Usage:       "log stream prefix.",
					EnvVar:      "AWS_LOG_STREAM_PREFIX",
					Destination: &putLogFlag.Prefix,
				},
				cli.StringFlag{
					Name:        "message,m",
					Usage:       "Message string.",
					Required:    true,
					Destination: &putLogFlag.Message,
				},
			},
		},
	}
	_ = app.Run(os.Args)
}

func action(c *cli.Context) {
	log.SetConfig(log.Config{Debug: debug})

	if err := run(c); err != nil {
		log.Get().Error("failed", zap.Error(err))
		os.Exit(-1)
	}

	os.Exit(0)
}

func run(c *cli.Context) error {
	s, err := session.NewSessionWithOptions(session.Options{
		Profile: profile,
		Config: aws.Config{
			Region: aws.String(lightsail.RegionNameApNortheast1),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	com := command.New(&command.Config{
		Session: s,
	})

	switch c.Command.Name {
	case "put-log":
		if err := com.NewPutLog().Run(&command.PutLogInput{
			LogGroupName: putLogFlag.LogGroupName,
			Prefix:       putLogFlag.Prefix,
			Message:      putLogFlag.Message,
		}); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
