package users

import (
	"context"

	"github.com/hibiken/asynq"
)

const TaskSendEmail = "task:send_email"

type PayloadSendEmail struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Subject   string `json:"subject"`
	BodyHtml  string `json:"bodyHtml"`
	BodyPlain string `json:"bodyPlain"`
}

type UserRedisTaskDistributor interface {
	DistributeTaskSendEmail(ctx context.Context, payload *PayloadSendEmail, opts ...asynq.Option) error
}

type UserRedisTaskProcessor interface {
	ProcessTaskSendEmail(ctx context.Context, task *asynq.Task) error
}
