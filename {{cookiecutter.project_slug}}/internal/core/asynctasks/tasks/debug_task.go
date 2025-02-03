package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

// Define task type
const TypeDebugTask = "debug:task"

// Define payload
type DebugTaskPayload struct {
	UserID  int
	Message string
}

// Define a function to handle this task
func HandleDebugTask(ctx context.Context, t *asynq.Task) error {
	var p DebugTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Running debug task: user_id=%d, message=%s", p.UserID, p.Message)
	return nil
}
