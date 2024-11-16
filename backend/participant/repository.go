package participant

import (
	"errors"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"gorm.io/gorm"
)

type Repository interface {
	CreateParticipants(participants []*Participant) ([]*Participant, error)
	GetParticipantsByExamID(examID uint) ([]*Participant, error)
	GetParticipantByID(id uint) (*Participant, error)
	// TODO: get by exam id, name, and password (participant side)
	UpdateParticipant(participant *Participant) error
	DeleteParticipantByID(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateParticipants(participants []*Participant) ([]*Participant, error) {
	err := r.db.CreateInBatches(participants, constants.InsertionBatchSize).Error
	return participants, err
}

func (r *repository) GetParticipantsByExamID(examID uint) ([]*Participant, error) {
	var res []*Participant
	err := r.db.Where("exam_id = ?", examID).Find(&res).Error
	return res, err
}

func (r *repository) GetParticipantByID(id uint) (*Participant, error) {
	var question Participant
	err := r.db.Where("id = ?", id).First(&question).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrParticipantNotFound
		}
		return nil, err
	}
	return &question, nil
}

func (r *repository) UpdateParticipant(participant *Participant) error {
	res := r.db.Updates(participant)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][UpdateParticipant] error: %s", res.Error)
		return lib.ErrParticipantNotFound
	}
	return nil
}

func (r *repository) DeleteParticipantByID(id uint) error {
	res := r.db.Model(&Participant{}).Where("id = ?", id).Delete(&Participant{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		log.Printf("[mcqoption][repository][DeleteParticipantByID] error: %s", res.Error)
		return lib.ErrParticipantNotFound
	}
	return nil
}
