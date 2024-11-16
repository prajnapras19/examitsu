package participant

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateParticipants(participants []*Participant) ([]*Participant, error)
	GetParticipantsByExamID(examID uint) ([]*Participant, error)
	GetParticipantByID(id uint) (*Participant, error)
	UpdateParticipant(participant *Participant) error
	DeleteParticipantByID(id uint) error
}

type service struct {
	cfg                   *config.Config
	participantRepository Repository
}

func NewService(
	cfg *config.Config,
	participantRepository Repository,
) Service {
	return &service{
		cfg:                   cfg,
		participantRepository: participantRepository,
	}
}

func (s *service) CreateParticipants(participants []*Participant) ([]*Participant, error) {
	var err error

	randomPassword, err := lib.GenerateRandomString(s.cfg.ParticipantRandomPasswordLength)
	if err != nil {
		log.Println("[participant][service][CreateParticipants] failed to generate password:", err.Error())
		return nil, lib.ErrFailedToGenerateRandomString
	}
	for i := range participants {
		if participants[i].Password == "" {
			participants[i].Password = randomPassword
		}
	}

	res, err := s.participantRepository.CreateParticipants(participants)
	if err != nil {
		log.Println("[participant][service][CreateParticipants] failed to create participants:", err.Error())
		return nil, lib.ErrFailedToCreateParticipants
	}

	return res, err
}

func (s *service) GetParticipantsByExamID(examID uint) ([]*Participant, error) {
	res, err := s.participantRepository.GetParticipantsByExamID(examID)
	if err != nil {
		log.Println("[mcqoption][service][GetParticipantsByExamID] failed to get participants by exam id:", err.Error())
		if errors.Is(err, lib.ErrParticipantNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetParticipants
	}
	return res, nil
}

func (s *service) GetParticipantByID(id uint) (*Participant, error) {
	res, err := s.participantRepository.GetParticipantByID(id)
	if err != nil {
		log.Println("[participant][service][GetQuestionByID] failed to get participant by id:", err.Error())
		if errors.Is(err, lib.ErrQuestionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetQuestionByID
	}
	return res, nil
}

func (s *service) UpdateParticipant(participant *Participant) error {
	err := s.participantRepository.UpdateParticipant(participant)
	if err != nil {
		log.Println("[participant][service][UpdateParticipant] failed to update participant:", err.Error())
		return lib.ErrFailedToUpdateParticipant
	}
	return nil
}

func (s *service) DeleteParticipantByID(id uint) error {
	err := s.participantRepository.DeleteParticipantByID(id)
	if err != nil {
		log.Println("[participant][service][DeleteParticipantByID] failed to delete participant:", err.Error())
		if errors.Is(err, lib.ErrParticipantNotFound) {
			return err
		}
		return lib.ErrFailedToDeleteParticipant
	}
	return nil
}
