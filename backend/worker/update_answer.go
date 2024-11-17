package worker

import (
	"log"

	rmq "github.com/adjust/rmq/v5"
	"github.com/prajnapras19/project-form-exam-sman2/backend/submission"
)

type UpdateAnswerQueueConsumer struct {
	submissionService submission.Service
}

func NewUpdateAnswerQueueConsumer(
	submissionService submission.Service,
) *UpdateAnswerQueueConsumer {
	return &UpdateAnswerQueueConsumer{
		submissionService: submissionService,
	}
}

func (consumer *UpdateAnswerQueueConsumer) Consume(delivery rmq.Delivery) {
	redisKey := delivery.Payload()

	log.Println("[worker][UpdateAnswerQueueConsumer][Consume] redisKey", redisKey)
	err := consumer.submissionService.UpsertSubmissionInDB(redisKey)
	if err != nil {
		// TODO: retry
	}
	if err := delivery.Ack(); err != nil {
		// TODO: handle ack error
	}
}
