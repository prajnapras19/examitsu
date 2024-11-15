package exam

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Repository interface {
	CreateExam(exam *Exam) (*Exam, error)
	GetExamBySerial(serial string) (*Exam, error)
	GetAllExams() ([]*Exam, error)
	UpdateExam(exam *Exam) error
	DeleteExamBySerial(serial string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateExam(exam *Exam) (*Exam, error) {
	err := r.db.Create(exam).Error
	return exam, err
}

func (r *repository) GetExamBySerial(serial string) (*Exam, error) {
	var exam Exam
	err := r.db.Where("serial = ?", serial).First(&exam).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrExamNotFound
		}
		return nil, err
	}
	return &exam, nil
}

func (r *repository) GetAllExams() ([]*Exam, error) {
	var exams []*Exam
	err := r.db.Find(&exams).Error
	if err != nil {
		return nil, err
	}
	return exams, nil
}

func (r *repository) UpdateExam(exam *Exam) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Exam{}).
			Where("serial = ?", exam.Serial).
			Update("name", exam.Name).
			Error; err != nil {
			return err
		}

		if err := tx.Model(&Exam{}).
			Where("serial = ?", exam.Serial).
			Update("is_open", exam.Name).
			Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) DeleteExamBySerial(serial string) error {
	res := r.db.Model(&Exam{}).Where("serial = ?", serial).Delete(&Exam{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[exam][repository][DeleteExamBySerial] error: %s", res.Error)
		return lib.ErrExamNotFound
	}
	return nil
}
