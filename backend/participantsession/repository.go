package participantsession

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Repository interface {
	CreateParticipantSession(participantSession *ParticipantSession) (*ParticipantSession, error)
	GetParticipantSessionBySerial(serial string) (*ParticipantSession, error)
	GetLatestAuthorizedParticipantSessionByParticipantID(participantID uint) (*ParticipantSession, error)
	AuthorizeSession(serial string, durationMinutes uint) error
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

func (r *repository) CreateParticipantSession(participantSession *ParticipantSession) (*ParticipantSession, error) {
	err := r.db.Create(participantSession).Error
	return participantSession, err
}

func (r *repository) GetParticipantSessionBySerial(serial string) (*ParticipantSession, error) {
	var participantSession ParticipantSession

	cacheKey := r.GetParticipantSessionBySerialCacheKey(serial)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &participantSession)
		return &participantSession, nil
	}

	err = r.db.Where("serial = ? AND not_archived", serial).First(&participantSession).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrParticipantSessionNotFound
		}
		return nil, err
	}
	res, _ := json.Marshal(participantSession)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &participantSession, nil
}

func (r *repository) GetLatestAuthorizedParticipantSessionByParticipantID(participantID uint) (*ParticipantSession, error) {
	var participantSession ParticipantSession

	cacheKey := r.GetLatestAuthorizedParticipantSessionByParticipantIDCacheKey(participantID)
	val, err := r.cache.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &participantSession)
		return &participantSession, nil
	}

	err = r.db.Where("participant_id = ? AND is_authorized AND not_archived", participantID).Order("updated_at DESC").First(&participantSession).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, lib.ErrParticipantSessionNotFound
		}
		return nil, err
	}
	res, _ := json.Marshal(participantSession)
	r.cache.Set(context.Background(), cacheKey, res, r.cfg.CacheTTL)
	return &participantSession, nil
}

func (r *repository) AuthorizeSession(serial string, durationMinutes uint) error {
	var startExam bool
	currentData, err := r.GetParticipantSessionBySerial(serial)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.ErrFailedToGetParticipantSession
		} else {
			return err
		}
	}

	r.cache.Del(context.Background(), r.GetParticipantSessionBySerialCacheKey(currentData.Serial))
	r.cache.Del(context.Background(), r.GetLatestAuthorizedParticipantSessionByParticipantIDCacheKey(currentData.ParticipantID))

	// if there is no authorized session previously, then this function should update participant's start time and duration
	_, err = r.GetLatestAuthorizedParticipantSessionByParticipantID(currentData.ParticipantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.ErrFailedToGetParticipantSession
		} else {
			return err
		}
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ParticipantSession{}).Update("is_authorized", true).Error
		if err != nil {
			return err
		}
		if startExam {
			return tx.Table("participants").Where("id = ?", currentData.ParticipantID).Updates(
				map[string]interface{}{
					"started_at":               time.Now(),
					"allowed_duration_minutes": durationMinutes,
				}).Error
		}
		return nil
	})
}

func (r *repository) GetParticipantSessionBySerialCacheKey(serial string) string {
	return fmt.Sprintf("participantSession:serial:%s", serial)
}

func (r *repository) GetLatestAuthorizedParticipantSessionByParticipantIDCacheKey(participantID uint) string {
	return fmt.Sprintf("participantSession:latestAuthorized:participantID:%d", participantID)
}