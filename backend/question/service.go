package question

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateQuestion(question *Question) (*Question, error)
	GetQuestionByID(id uint) (*Question, error)
	GetQuestionsIDOnly(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error)
	GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error)
	UpdateQuestion(question *Question) error
	DeleteQuestionByID(id uint) error
}

type service struct {
	questionRepository Repository
}

func NewService(
	questionRepository Repository,
) Service {
	return &service{
		questionRepository: questionRepository,
	}
}

func (s *service) CreateQuestion(question *Question) (*Question, error) {
	var err error

	res, err := s.questionRepository.CreateQuestion(question)
	if err != nil {
		log.Println("[question][service][CreateQuestion] failed to create question:", err.Error())
		return nil, lib.ErrFailedToCreateQuestion
	}

	return res, err
}

func (s *service) GetQuestionByID(id uint) (*Question, error) {
	res, err := s.questionRepository.GetQuestionByID(id)
	if err != nil {
		log.Println("[question][service][GetQuestionByID] failed to get question by id:", err.Error())
		if errors.Is(err, lib.ErrQuestionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetQuestionByID
	}
	return res, nil
}

func (s *service) GetQuestionsIDOnly(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error) {
	res, err := s.questionRepository.GetQuestionsIDOnly(pagination, filter)
	if err != nil {
		log.Println("[question][service][GetQuestionsIDOnly] failed to get questions:", err.Error())
		return nil, lib.ErrFailedToGetQuestions
	}
	return res, nil
}

func (s *service) GetQuestions(pagination *lib.QueryPagination, filter *GetQuestionsFilter) ([]*Question, error) {
	res, err := s.questionRepository.GetQuestions(pagination, filter)
	if err != nil {
		log.Println("[question][service][GetQuestions] failed to get questions:", err.Error())
		return nil, lib.ErrFailedToGetQuestions
	}
	return res, nil
}

func (s *service) UpdateQuestion(question *Question) error {
	err := s.questionRepository.UpdateQuestionDataByID(question)
	if err != nil {
		log.Println("[question][service][UpdateQuestion] failed to update question:", err.Error())
		return lib.ErrFailedToUpdateQuestion
	}
	return nil
}

func (s *service) DeleteQuestionByID(id uint) error {
	err := s.questionRepository.DeleteQuestionByID(id)
	if err != nil {
		log.Println("[question][service][DeleteQuestionByID] failed to delete question:", err.Error())
		if errors.Is(err, lib.ErrQuestionNotFound) {
			return err
		}
		return lib.ErrFailedToDeleteQuestion
	}
	return nil
}
