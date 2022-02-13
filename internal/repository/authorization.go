package repository

import (
	"github.com/whiskey-back/internal/types"
	"gorm.io/gorm"
)

type AuthRepo struct {
	db *gorm.DB
}

func NewAuthorizationRepo(db *gorm.DB) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (p *AuthRepo) GetUserByPhone(phone string) (*types.User, error) {

	user := new(types.User)
	if err := p.db.Debug().Table("users").Where("phone = ?", phone).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *AuthRepo) CreateUser(user *types.User) error {

	if err := p.db.Debug().Table("users").Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) PreCreateUser(user *types.CheckUser) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		if tx.Debug().Table("check_users").Where("phone = ?", user.Phone).Take(user).RowsAffected == 0 {
			if err := tx.Debug().Table("check_users").Create(user).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (p *AuthRepo) DeleteCodePhone(userID string) error {

	var code string
	if err := p.db.Debug().Table("users").Where("id = ?", userID).Update("code", code).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) AddCodePhone(code string, userID string) error {

	if err := p.db.Debug().Table("users").Where("id = ?", userID).Update("code", code).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) PreDeleteCodePhone(userID string) error {

	var code string
	if err := p.db.Debug().Table("check_users").Where("id = ?", userID).Update("code", code).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) PreAddCodePhone(code string, userID string) error {

	if err := p.db.Debug().Table("check_users").Where("id = ?", userID).Update("code", code).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) GetCheckUser(phone, code string) (*types.CheckUser, error) {

	user := new(types.CheckUser)
	if err := p.db.Debug().Table("check_users").Where("phone = ? AND code = ?", phone, code).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *AuthRepo) DeleteCheckUser(preUser *types.CheckUser) error {

	if err := p.db.Debug().Table("check_users").Delete(preUser).Error; err != nil {
		return err
	}

	return nil
}

func (p *AuthRepo) GetAuthUser(phone, code string) (*types.User, error) {

	user := new(types.User)
	if err := p.db.Debug().Table("users").Where("phone = ? AND code = ?", phone, code).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *AuthRepo) AddLogRespSMS(logSMS *types.LogSMS) error {

	if err := p.db.Debug().Table("log_sms").Create(logSMS).Error; err != nil {
		return err
	}

	return nil
}
