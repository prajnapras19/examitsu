package worker

import (
	"time"

	rmq "github.com/adjust/rmq/v5"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
)

type Service interface {
	InitConsumers()
}

type service struct {
	cfg                       *config.Config
	updateAnswerQueue         rmq.Queue
	updateAnswerQueueConsumer *UpdateAnswerQueueConsumer
}

func NewService(
	cfg *config.Config,
	updateAnswerQueue rmq.Queue,
	updateAnswerQueueConsumer *UpdateAnswerQueueConsumer,
) Service {
	return &service{
		cfg:                       cfg,
		updateAnswerQueue:         updateAnswerQueue,
		updateAnswerQueueConsumer: updateAnswerQueueConsumer,
	}
}

func (s *service) InitConsumers() {
	s.updateAnswerQueue.StartConsuming(s.cfg.UpdateAnswerQueuePrefetchLimit, time.Second)
	s.updateAnswerQueue.AddConsumer(constants.UpdateAnswerConsumerName, s.updateAnswerQueueConsumer)
}
