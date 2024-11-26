package exam

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
)

type Service interface {
	CreateExam(exam *Exam) (*Exam, error)
	GetExamByID(id uint) (*Exam, error)
	GetExamBySerial(serial string) (*Exam, error)
	GetExams(pagination *lib.QueryPagination, filter *GetExamsFilter) ([]*Exam, error)
	GetAllOpenedExams() ([]*Exam, error)
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

	exam.Serial = uuid.New().String()

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
		return nil, lib.ErrFailedToGetExam
	}
	return res, nil
}

func (s *service) GetExamByID(id uint) (*Exam, error) {
	res, err := s.examRepository.GetExamByID(id)
	if err != nil {
		log.Println("[exam][service][GetExamByID] failed to get exam by id:", err.Error())
		if errors.Is(err, lib.ErrExamNotFound) {
			return nil, err
		}
		return nil, lib.ErrFailedToGetExam
	}
	return res, nil
}

func (s *service) GetExams(pagination *lib.QueryPagination, filter *GetExamsFilter) ([]*Exam, error) {
	res, err := s.examRepository.GetExams(pagination, filter)
	if err != nil {
		log.Println("[exam][service][GetExams] failed to get exams:", err.Error())
		return nil, lib.ErrFailedToGetExams
	}
	return res, nil
}

func (s *service) GetAllOpenedExams() ([]*Exam, error) {
	res, err := s.examRepository.GetAllOpenedExams()
	if err != nil {
		log.Println("[exam][service][GetAllOpenedExams] failed to get exams:", err.Error())
		return nil, lib.ErrFailedToGetExams
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
