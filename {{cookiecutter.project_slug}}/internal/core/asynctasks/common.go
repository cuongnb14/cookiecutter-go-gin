package asynctasks

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"{{ cookiecutter.project_slug }}/configs"
)

func AddAsynqmonRoute(r *gin.Engine) {
	basicAuthMidd := gin.BasicAuth(gin.Accounts{
		configs.Env.BasicAuthUser: configs.Env.BasicAuthPass,
	})

	h := asynqmon.New(asynqmon.Options{
		RootPath:     "/asynqmon",
		RedisConnOpt: configs.GetAsynqRedisOpt(),
	})

	r.Any(h.RootPath()+"/*a", basicAuthMidd, gin.WrapH(h))
}

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
	_, err = configs.GetAsynqClient.Enqueue(task, opts...)
	if err != nil {
		return fmt.Errorf("could not enqueue task: %w", err)
	}
	return nil
}
