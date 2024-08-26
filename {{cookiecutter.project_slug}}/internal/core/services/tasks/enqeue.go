package tasks

import (
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"{{ cookiecutter.project_slug }}/configs"
)

func EnqueueWalletTransactionsProcessingTask(txIDs []uuid.UUID) {
	logger := configs.GetLogger()
	task, err := NewWalletTransactionsProcessingTask(txIDs)
	if err != nil {
		logger.Error("could not create task:", "err", err)
	}
	info, err := configs.AsynqClient.Enqueue(task, asynq.Retention(24*time.Hour), asynq.Timeout(3*time.Second))
	if err != nil {
		logger.Error("could not enqueue task:", "err", err)
	}
}
