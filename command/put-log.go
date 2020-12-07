package command

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/pkg/errors"
)

type PutLogCommand struct {
	*Command
}

func (c *Command) NewPutLog() *PutLogCommand {
	return &PutLogCommand{
		Command: c,
	}
}

type PutLogInput struct {
	LogGroupName string
	Prefix       string
	Message      string
}

func (c *PutLogCommand) Run(i *PutLogInput) error {
	logs := cloudwatchlogs.New(c.session)

	streams, err := logs.DescribeLogStreams(&cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        &i.LogGroupName,
		LogStreamNamePrefix: &i.Prefix,
		Descending:          aws.Bool(true),
		Limit:               aws.Int64(1),
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if len(streams.LogStreams) == 0 {
		return errors.New("no LogStream exists.")
	}

	println(streams.LogStreams[0].String())

	result, err := logs.PutLogEvents(&cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &i.LogGroupName,
		LogStreamName: streams.LogStreams[0].LogStreamName,
		SequenceToken: streams.LogStreams[0].UploadSequenceToken,
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   &i.Message,
				Timestamp: aws.Int64(time.Now().UnixNano() / 1e6),
			},
		},
	})
	if err != nil {
		return errors.WithStack(err)
	}

	println(result.String())

	return nil
}
