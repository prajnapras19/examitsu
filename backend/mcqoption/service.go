package mcqoption

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateMcqOption(mcqOption *McqOption) (*McqOption, error)
	GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error)
	UpdateMcqOption(mcqOption *McqOption) error
	DeleteMcqOptionByID(id uint) error
}

type service struct {
	mcqOptionRepository Repository
}

func NewService(
	mcqOptionRepository Repository,
) Service {
	return &service{
		mcqOptionRepository: mcqOptionRepository,
	}
}

func (s *service) CreateMcqOption(mcqOption *McqOption) (*McqOption, error) {
	var err error

	res, err := s.mcqOptionRepository.CreateMcqOption(mcqOption)
	if err != nil {
		log.Println("[mcqOption][service][CreateMcqOption] failed to create mcqOption:", err.Error())
		return nil, lib.ErrFailedToCreateMcqOption
	}

	return res, err
}

func (s *service) GetMcqOptionsByQuestionID(questionID uint) ([]*McqOption, error) {
	res, err := s.mcqOptionRepository.GetMcqOptionsByQuestionID(questionID)
	if err != nil {
		log.Println("[mcqoption][service][GetMcqOptionByQuestionID] failed to get question by id:", err.Error())
		if errors.Is(err, lib.ErrMcqOptionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetMcqOptions
	}
	return res, nil
}

func (s *service) UpdateMcqOption(mcqOption *McqOption) error {
	err := s.mcqOptionRepository.UpdateMcqOption(mcqOption)
	if err != nil {
		log.Println("[mcqOption][service][UpdateMcqOption] failed to update mcqOption:", err.Error())
		return lib.ErrFailedToUpdateMcqOption
	}
	return nil
}

func (s *service) DeleteMcqOptionByID(id uint) error {
	err := s.mcqOptionRepository.DeleteMcqOptionByID(id)
	if err != nil {
		log.Println("[mcqOption][service][DeleteMcqOptionByID] failed to delete mcqOption:", err.Error())
		if errors.Is(err, lib.ErrMcqOptionNotFound) {
			return err
		}
		return lib.ErrFailedToDeleteMcqOption
	}
	return nil
}
