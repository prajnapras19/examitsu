package participantsession

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateParticipantSession(participantSession *ParticipantSession) (*ParticipantSession, error)
	GetParticipantSessionBySerial(serial string) (*ParticipantSession, error)
	GetLatestAuthorizedParticipantSessionByParticipantID(participantID uint) (*ParticipantSession, error)
	AuthorizeSession(serial string, durationMinutes uint) error
}

type service struct {
	participantSessionRepository Repository
}

func NewService(
	participantSessionRepository Repository,
) Service {
	return &service{
		participantSessionRepository: participantSessionRepository,
	}
}

func (s *service) CreateParticipantSession(participantSession *ParticipantSession) (*ParticipantSession, error) {
	var err error

	participantSession.Serial = uuid.New().String()

	res, err := s.participantSessionRepository.CreateParticipantSession(participantSession)
	if err != nil {
		log.Println("[participantsession][service][CreateParticipantSession] failed to create participant session:", err.Error())
		return nil, lib.ErrFailedToCreateParticipantSession
	}

	return res, err
}

func (s *service) GetParticipantSessionBySerial(serial string) (*ParticipantSession, error) {
	res, err := s.participantSessionRepository.GetParticipantSessionBySerial(serial)
	if err != nil {
		log.Println("[participantsession][service][GetParticipantSessionBySerial] failed to get participant session by serial:", err.Error())
		if errors.Is(err, lib.ErrParticipantSessionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetParticipantSession
	}
	return res, nil
}

func (s *service) GetLatestAuthorizedParticipantSessionByParticipantID(participantID uint) (*ParticipantSession, error) {
	res, err := s.participantSessionRepository.GetLatestAuthorizedParticipantSessionByParticipantID(participantID)
	if err != nil {
		log.Println("[participantsession][service][GetLatestAuthorizedParticipantSessionByParticipantID] failed to get latest authorized participant session by participant id:", err.Error())
		if errors.Is(err, lib.ErrParticipantSessionNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetParticipantSession
	}
	return res, nil
}

func (s *service) AuthorizeSession(serial string, durationMinutes uint) error {
	err := s.participantSessionRepository.AuthorizeSession(serial, durationMinutes)
	if err != nil {
		log.Println("[participantsession][service][AuthorizeSession] failed to authorize session:", err.Error())
		return lib.ErrFailedToAuthorizeParticipantSession
	}
	return nil
}
