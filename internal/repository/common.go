package repository

import (
	"github.com/whiskey-back/internal/types"
	"gorm.io/gorm"
)

type CommonRepo struct {
	db *gorm.DB
}

func NewCommonRepo(db *gorm.DB) *CommonRepo {
	return &CommonRepo{
		db: db,
	}
}

func (p *CommonRepo) GetUserByID(userID string) (*types.User, error) {

	user := new(types.User)
	if err := p.db.Debug().Table("users").Where("id = ?", userID).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *CommonRepo) GetGifts() (types.Gifts, error) {

	gifts := make(types.Gifts, 0)
	if err := p.db.Debug().Table("gifts").Find(&gifts).Error; err != nil {
		return nil, err
	}

	return gifts, nil
}
