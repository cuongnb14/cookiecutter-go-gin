package distributed_tasks

import (
	"time"

	"github.com/hibiken/asynq"
	"{{ cookiecutter.project_slug }}/internal/core/distributed_tasks/tasks"
)

func PublishDebugTaskTask() error {
	return PublishTask(
		tasks.TypeDebugTask,
		tasks.DebugTaskPayload{UserID: 1, Message: "hello"},
		asynq.Retention(24*time.Hour), asynq.Timeout(3*time.Second),
	)
}

func PublishSendEmailTask() error {
	return PublishTask(
		tasks.TypeEmailDelivery,
		tasks.EmailDeliveryPayload{UserID: 1, TemplateID: "welcome-email"},
		asynq.Retention(24*time.Hour), asynq.Timeout(3*time.Second),
	)
}
