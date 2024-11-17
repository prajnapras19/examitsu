package submission

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	rmq "github.com/adjust/rmq/v5"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
)

type Service interface {
	Answer(cacheObject *ExamSessionSubmissionCacheObject) error
	UpsertSubmissionInDB(key string) error
}

type service struct {
	submissionRepository Repository
	redisClient          *redis.Client
	updateAnswerQueue    rmq.Queue
}

func NewService(
	submissionRepository Repository,
	redisClient *redis.Client,
	updateAnswerQueue rmq.Queue,
) Service {
	return &service{
		submissionRepository: submissionRepository,
		redisClient:          redisClient,
		updateAnswerQueue:    updateAnswerQueue,
	}
}

func (s *service) Answer(cacheObject *ExamSessionSubmissionCacheObject) error {
	err := s.submissionRepository.SaveCacheObject(cacheObject)
	if err != nil {
		log.Println("[submission][service][Answer] failed to save answer:", err.Error())
		return lib.ErrFailedToSaveAnswer
	}
	s.updateAnswerQueue.Publish(cacheObject.GetKey())
	return nil
}

func (s *service) GetAnswer(participantID uint, questionID uint) (*Submission, error) {
	res, err := s.submissionRepository.GetSubmissionByParticipantIDAndQuestionID(participantID, questionID)
	if err != nil {
		log.Println("[submission][service][GetAnswer] failed to get answer:", err.Error())
		if errors.Is(err, lib.ErrSubmissionNotFound) {
			return nil, lib.ErrAnswerNotFound
		}
		return nil, lib.ErrFailedToGetAnswer
	}
	return res, nil
}

func (s *service) UpsertSubmissionInDB(key string) error {
	val, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	var cacheObject ExamSessionSubmissionCacheObject
	json.Unmarshal([]byte(val), &cacheObject)
	return s.submissionRepository.UpsertSubmissionInDB(&cacheObject)
}
