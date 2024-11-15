package mysql

import (
	"fmt"

	"github.com/prajnapras19/project-form-exam-sman2/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service interface {
	InitDB() *gorm.DB
	GetDSN() string
	GetDB() *gorm.DB
}

type service struct {
	cfg config.MySQLConfig
	db  *gorm.DB
}

func NewService(mySQLConfig config.MySQLConfig) Service {
	svc := &service{
		cfg: mySQLConfig,
	}
	svc.db = svc.InitDB()
	return svc
}

func (s *service) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		s.cfg.Username, s.cfg.Password, s.cfg.Hostname, s.cfg.Port, s.cfg.DBName,
		s.cfg.Charset, s.cfg.ParseTime, s.cfg.Loc,
	)
}

func (s *service) InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(s.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(s.cfg.GORMLogLevel),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func (s *service) GetDB() *gorm.DB {
	return s.db
}
