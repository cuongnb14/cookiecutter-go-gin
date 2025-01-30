package distributed_tasks

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"{{ cookiecutter.project_slug }}/configs"
)

func NewTask[T any](taskType string, data T) (*asynq.Task, error) {
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(taskType, payload), nil
}

func PublishTask[T any](taskType string, payload T, opts ...asynq.Option) error {
	task, err := NewTask(taskType, payload)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	_, err = configs.AsynqClient.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	return nil
}
