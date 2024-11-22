package participant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	CreateParticipants(participants []*Participant) ([]*Participant, error)
	GetParticipantsByExamID(examID uint) ([]*Participant, error)
	GetParticipantByID(id uint) (*Participant, error)
	GetParticipantByExamIDAndName(examID uint, name string) (*Participant, error)
	UpdateParticipant(participant *Participant) error
	DeleteParticipantByID(id uint) error

	GetParticipantTotalPointsByExamID(examID uint) ([]*ParticipantTotalPoint, error)

	GetParticipantByIDCacheKey(id uint) string
	GetParticipantByExamIDAndNameCacheKey(examID uint, name string) string
}

type repository struct {
	cfg   *config.Config
	db    *gorm.DB
	cache *redis.Client
}

func NewRepository(
	cfg *config.Config,
	db *gorm.DB,
	cache *redis.Client,
) Repository {
	return &repository{
		cfg:   cfg,
		db:    db,
		cache: cache,
	}
}

func (r *repository) CreateParticipants(participants []*Participant) ([]*Participant, error) {
	err := r.db.CreateInBatches(participants, constants.InsertionBatchSize).Error
	return participants, err
}

func (r *repository) GetParticipantsByExamID(examID uint) ([]*Participant, error) {
	var res []*Participant
	err := r.db.Where("exam_id = ?", examID).Order("id ASC").Find(&res).Error
	return res, err
}

func (r *repository) GetParticipantByID(id uint) (*Participant, error) {
	var participant Participant

	cacheKey := r.GetParticipantByIDCacheKey(id)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &participant)
		return &participant, nil
	}

	err = r.db.Where("id = ?", id).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrParticipantNotFound
		}
		return nil, err
	}

	res, _ := json.Marshal(participant)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &participant, nil
}

func (r *repository) GetParticipantByExamIDAndName(examID uint, name string) (*Participant, error) {
	var participant Participant

	cacheKey := r.GetParticipantByExamIDAndNameCacheKey(examID, name)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &participant)
		return &participant, nil
	}

	err = r.db.Where("exam_id = ? AND name = ? AND not_archived", examID, name).First(&participant).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrParticipantNotFound
		}
		return nil, err
	}

	res, _ := json.Marshal(participant)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &participant, nil
}

func (r *repository) UpdateParticipant(participant *Participant) error {
	currentData, err := r.GetParticipantByID(participant.ID)
	if err != nil {
		return err
	}

	r.cache.Del(context.Background(), r.GetParticipantByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetParticipantByExamIDAndNameCacheKey(currentData.ExamID, currentData.Name))

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
	currentData, err := r.GetParticipantByID(id)
	if err != nil {
		return err
	}

	r.cache.Del(context.Background(), r.GetParticipantByIDCacheKey(currentData.ID))
	r.cache.Del(context.Background(), r.GetParticipantByExamIDAndNameCacheKey(currentData.ExamID, currentData.Name))

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

func (r *repository) GetParticipantTotalPointsByExamID(examID uint) ([]*ParticipantTotalPoint, error) {
	var res []*ParticipantTotalPoint
	err := r.db.Raw(`
		SELECT
		    p.id AS participant_id,
		    COALESCE(SUM(m.point), 0) AS total_point
		FROM
		    participants p
		LEFT JOIN
		    submissions s
		ON
		    p.id = s.participant_id
		LEFT JOIN
		    mcq_options m
		ON
		    s.mcq_option_id = m.id
		WHERE
		    p.exam_id = ?
			AND p.deleted_at IS NULL
			AND s.deleted_at IS NULL
			AND m.deleted_at IS NULL
		GROUP BY
		    1
		ORDER BY
		    p.id ASC;
	`, examID).Scan(&res).Error
	return res, err
}

func (r *repository) GetParticipantByIDCacheKey(id uint) string {
	return fmt.Sprintf("participant:id:%d", id)
}

func (r *repository) GetParticipantByExamIDAndNameCacheKey(examID uint, name string) string {
	return fmt.Sprintf("participant:examID:%d:name:%s", examID, name)
}
