package participant

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/exam"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateParticipants(participants []*Participant) ([]*Participant, error)
	GetParticipantsByExamID(examID uint) ([]*Participant, error)
	GetParticipantByID(id uint) (*Participant, error)
	UpdateParticipant(participant *Participant) error
	DeleteParticipantByID(id uint) error
	GetParticipantByExamIDAndName(examID uint, name string) (*Participant, error)

	GenerateToken(examSerial string, participantID uint) string
	VerifyToken(token *jwt.Token) (interface{}, error)
	ValidateToken(tokenString string) (*lib.ExamTokenJWTClaims, error)
}

type service struct {
	cfg                   *config.Config
	participantRepository Repository
	examService           exam.Service
}

func NewService(
	cfg *config.Config,
	participantRepository Repository,
	examService exam.Service,
) Service {
	return &service{
		cfg:                   cfg,
		participantRepository: participantRepository,
		examService:           examService,
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
		log.Println("[participant][service][GetParticipantByID] failed to get participant by id:", err.Error())
		if errors.Is(err, lib.ErrParticipantNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetParticipant
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

func (s *service) GetParticipantByExamIDAndName(examID uint, name string) (*Participant, error) {
	res, err := s.participantRepository.GetParticipantByExamIDAndName(examID, name)
	if err != nil {
		log.Println("[participant][service][GetParticipantByExamIDAndName] failed to get participant:", err.Error())
		if errors.Is(err, lib.ErrParticipantNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetParticipant
	}
	return res, nil
}

func (s *service) GenerateToken(examSerial string, participantID uint) string {
	claims := lib.ExamTokenJWTClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    s.cfg.AuthConfig.ApplicationName,
			ExpiresAt: time.Now().Add(s.cfg.AuthConfig.LoginTokenExpirationDuration).Unix(),
		},
		ParticipantID: participantID,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	signedToken, _ := token.SignedString(s.cfg.AuthConfig.SignatureKey)
	return signedToken
}

func (s *service) VerifyToken(token *jwt.Token) (interface{}, error) {
	if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, lib.ErrSigningMethodInvalid
	} else if method != jwt.SigningMethodHS256 {
		return nil, lib.ErrSigningMethodInvalid
	}
	return s.cfg.AuthConfig.SignatureKey, nil
}

func (s *service) ValidateToken(tokenString string) (*lib.ExamTokenJWTClaims, error) {
	claims := &lib.ExamTokenJWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, s.VerifyToken)
	if err != nil {
		return nil, lib.ErrUnauthorizedRequest
	}
	if !token.Valid {
		return nil, lib.ErrUnauthorizedRequest
	}
	claims, ok := token.Claims.(*lib.ExamTokenJWTClaims)
	if !ok {
		return nil, lib.ErrUnauthorizedRequest
	}

	if participant, err := s.GetParticipantByID(claims.ParticipantID); err != nil {
		return nil, lib.ErrUnauthorizedRequest
	} else {
		exam, err := s.examService.GetExamByID(participant.ExamID)
		if err != nil || !exam.IsOpen {
			return nil, lib.ErrUnauthorizedRequest
		}
	}

	return claims, nil
}
