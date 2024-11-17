package submission

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	Answer(cacheObject *ExamSessionSubmissionCacheObject) error
}

type service struct {
	submissionRepository Repository
}

func NewService(
	submissionRepository Repository,
) Service {
	return &service{
		submissionRepository: submissionRepository,
	}
}

func (s *service) Answer(cacheObject *ExamSessionSubmissionCacheObject) error {
	err := s.submissionRepository.SaveCacheObject(cacheObject)
	if err != nil {
		log.Println("[submission][service][Answer] failed to save answer:", err.Error())
		return lib.ErrFailedToSaveAnswer
	}
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
