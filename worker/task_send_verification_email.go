package worker

import (
	"context"
	"encoding/json"
	"fmt"
	db "github.com/aalug/blog-go/db/sqlc"
	"github.com/aalug/blog-go/mail"
	"github.com/aalug/blog-go/utils"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerificationEmail = "task:send_verification_email"

type PayloadSendVerificationEmail struct {
	Email string `json:"email"`
}

// DistributeTaskSendVerificationEmail distributes the task of sending a verification email.
func (distributor *RedisTaskDistributor) DistributeTaskSendVerificationEmail(
	ctx context.Context,
	payload *PayloadSendVerificationEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerificationEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")

	return nil
}

// ProcessTaskSendVerificationEmail processes the task of sending a verification email.
func (processor *RedisTaskProcessor) ProcessTaskSendVerificationEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerificationEmail
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// create verify email in the database
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Email:      user.Email,
		SecretCode: utils.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	// send email to user
	verifyUrl := fmt.Sprintf("http://localhost:8080/verify-email?id=%d&code=%s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`
		<h3>Hello %s</h3><br>
		<p class="message">
		Please click the link below to verify your email address:
		</p>
		<a class="button" href="%s">Verify Email</a>
		`, user.Username, verifyUrl)
	err = processor.emailSender.SendEmail(mail.Data{
		To:       []string{user.Email},
		Subject:  "Welcome to Blog Go!",
		Content:  content,
		Template: "verification_email.html",
	})
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("email", user.Email).Msg("processed task")

	return nil
}
