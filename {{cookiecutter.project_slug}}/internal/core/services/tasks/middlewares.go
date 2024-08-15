package tasks

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"{{ cookiecutter.project_slug }}/configs"
)

func LoggingMiddleware(h asynq.Handler) asynq.Handler {
	logger := configs.GetLogger()
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		start := time.Now()
		//log.Printf("start processing %q", t.Type())
		err := h.ProcessTask(ctx, t)
		if err != nil {
			logger.Errorf("failed to process %q: %v", t.Type(), err)
			return err
		}
		logger.Infof("finished processing %q: elapsed= %v", t.Type(), time.Since(start))
		return nil
	})
}
