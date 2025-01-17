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
	GetMcqOptionByID(id uint) (*McqOption, error)
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
		log.Println("[mcqoption][service][GetMcqOptionByQuestionID] failed to get mcqOption by question id:", err.Error())
		if errors.Is(err, lib.ErrMcqOptionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetMcqOptions
	}
	return res, nil
}

func (s *service) GetMcqOptionByID(id uint) (*McqOption, error) {
	res, err := s.mcqOptionRepository.GetMcqOptionByID(id)
	if err != nil {
		log.Println("[mcqoption][service][GetMcqOptionByID] failed to get mcqOption by id:", err.Error())
		if errors.Is(err, lib.ErrMcqOptionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetMcqOption
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
