package exam

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateExam(exam *Exam) (*Exam, error)
	GetExamBySerial(serial string) (*Exam, error)
	GetAllExams() ([]*Exam, error)
	UpdateExam(exam *Exam) error
	DeleteExamBySerial(serial string) error
}

type service struct {
	examRepository Repository
}

func NewService(
	examRepository Repository,
) Service {
	return &service{
		examRepository: examRepository,
	}
}

func (s *service) CreateExam(exam *Exam) (*Exam, error) {
	var err error

	exam.Serial, err = lib.GenerateRandomString(constants.ExamSerialLength)

	res, err := s.examRepository.CreateExam(exam)
	if err != nil {
		log.Println("[exam][service][CreateExam] failed to create exam:", err.Error())
		return nil, lib.ErrFailedToCreateExam
	}

	return res, err
}

func (s *service) GetExamBySerial(serial string) (*Exam, error) {
	res, err := s.examRepository.GetExamBySerial(serial)
	if err != nil {
		log.Println("[exam][service][GetExamBySerial] failed to get exam by serial:", err.Error())
		if errors.Is(err, lib.ErrExamNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetExamBySerial
	}
	return res, nil
}

func (s *service) GetAllExams() ([]*Exam, error) {
	res, err := s.examRepository.GetAllExams()
	if err != nil {
		log.Println("[exam][service][GetAllExams] failed to get all exams:", err.Error())
		return nil, lib.ErrFailedToGetAllExams
	}
	return res, nil
}

func (s *service) UpdateExam(exam *Exam) error {
	err := s.examRepository.UpdateExam(exam)
	if err != nil {
		log.Println("[exam][service][UpdateExam] failed to update exam:", err.Error())
		return lib.ErrFailedToUpdateExam
	}
	return nil
}

func (s *service) DeleteExamBySerial(serial string) error {
	err := s.examRepository.DeleteExamBySerial(serial)
	if err != nil {
		log.Println("[exam][service][DeleteExam] failed to delete exam:", err.Error())
		if errors.Is(err, lib.ErrExamNotFound) {
			return err
		}
		return lib.ErrFailedToDeleteExam
	}
	return nil
}
