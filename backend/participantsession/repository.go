package participantsession

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"github.com/prajnapras19/project-form-exam-sman2/backend/constants"
	"github.com/prajnapras19/project-form-exam-sman2/backend/lib"
	"github.com/prajnapras19/project-form-exam-sman2/backend/participant"
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
	cfg                   *config.Config
	db                    *gorm.DB
	cache                 *redis.Client
	participantRepository participant.Repository
}

func NewRepository(
	cfg *config.Config,
	db *gorm.DB,
	cache *redis.Client,
	participantRepository participant.Repository,
) Repository {
	return &repository{
		cfg:                   cfg,
		db:                    db,
		cache:                 cache,
		participantRepository: participantRepository,
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
		if string(val) != constants.None {
			json.Unmarshal([]byte(val), &participantSession)
			return &participantSession, nil
		} else {
			return nil, lib.ErrParticipantSessionNotFound
		}
	}

	err = r.db.Where("participant_id = ? AND is_authorized AND not_archived", participantID).Order("updated_at DESC").First(&participantSession).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.cache.Set(context.Background(), cacheKey, []byte(constants.None), r.cfg.CacheTTL)
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
	currentParticipant, err := r.participantRepository.GetParticipantByID(currentData.ParticipantID)

	// if there is no authorized session previously, then this function should update participant's start time and duration
	_, err = r.GetLatestAuthorizedParticipantSessionByParticipantID(currentData.ParticipantID)
	if err != nil {
		if errors.Is(err, lib.ErrParticipantSessionNotFound) {
			startExam = true
		} else {
			return err
		}
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ParticipantSession{}).Where("serial = ?", serial).Update("is_authorized", true).Error
		if err != nil {
			return err
		}
		if startExam {
			err = tx.Table("participants").Where("id = ?", currentData.ParticipantID).Updates(
				map[string]interface{}{
					"started_at":               time.Now(),
					"allowed_duration_minutes": durationMinutes,
				}).Error
		}
		return nil
	})

	if err == nil {
		r.cache.Del(context.Background(), r.GetParticipantSessionBySerialCacheKey(currentData.Serial))
		r.cache.Del(context.Background(), r.GetLatestAuthorizedParticipantSessionByParticipantIDCacheKey(currentData.ParticipantID))
		r.cache.Del(context.Background(), r.participantRepository.GetParticipantByIDCacheKey(currentParticipant.ID))
		r.cache.Del(context.Background(), r.participantRepository.GetParticipantByExamIDAndNameCacheKey(currentParticipant.ExamID, currentParticipant.Name))
	}
	return err
}

func (r *repository) GetParticipantSessionBySerialCacheKey(serial string) string {
	return fmt.Sprintf("participantSession:serial:%s", serial)
}

func (r *repository) GetLatestAuthorizedParticipantSessionByParticipantIDCacheKey(participantID uint) string {
	return fmt.Sprintf("participantSession:latestAuthorized:participantID:%d", participantID)
}
