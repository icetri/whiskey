package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "err with Open DB")
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, errors.Wrap(err, "err with ping DB")
	}

	return &Postgres{db: sqlDB}, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) Database() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: p.db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, errors.Wrap(err, "err with Open GORM")
	}

	return db, nil
}
