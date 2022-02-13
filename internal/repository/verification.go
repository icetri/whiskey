package repository

import (
	"github.com/whiskey-back/internal/types"
	"gorm.io/gorm"
)

type VerifiRepo struct {
	db *gorm.DB
}

func NewVerifiRepo(db *gorm.DB) *VerifiRepo {
	return &VerifiRepo{
		db: db,
	}
}

func (p *VerifiRepo) UpdateRequestGift(gift *types.RequestGift) (*types.RequestGift, error) {

	reqGift := new(types.RequestGift)
	if err := p.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Debug().Table("request_gift").Select("certificate").Updates(gift).Error; err != nil {
			return err
		}

		if err := tx.Debug().Table("request_gift").Where("id = ?", gift.ID).Take(reqGift).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return reqGift, nil
}

func (p *VerifiRepo) UpdateRequestGiftTrue(gift *types.RequestGift) error {

	if err := p.db.Debug().Table("request_gift").Select("sent").Where("id = ?", gift.ID).Updates(gift).Error; err != nil {
		return err
	}

	return nil
}

func (p *VerifiRepo) GetUserByID(userID string) (*types.User, error) {

	user := new(types.User)
	if err := p.db.Debug().Table("users").Where("id = ?", userID).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *VerifiRepo) GetGift(giftID int) (*types.Gift, error) {

	gift := new(types.Gift)
	if err := p.db.Debug().Table("gifts").Where("id = ?", giftID).Take(gift).Error; err != nil {
		return nil, err
	}

	return gift, nil
}

func (p *VerifiRepo) CreateCheque(check *types.Cheque) error {

	if err := p.db.Debug().Table("cheques").Create(check).Error; err != nil {
		return err
	}

	return nil
}

func (p *VerifiRepo) CountCheques(check *types.Cheque) (int64, error) {

	var count int64
	if err := p.db.Debug().Table("cheques").
		Where("fn = ? AND fd = ? AND fp = ?", check.FN, check.FD, check.FP).
		Count(&count).Error; err != nil {
		return 10, err
	}

	return count, nil
}

func (p *VerifiRepo) UpdateBalanceCheck(userID string, balance float64) error {

	if err := p.db.Debug().Table("users").Select("balance").Where("id = ?", userID).Updates(&types.Profile{Balance: balance}).Error; err != nil {
		return err
	}

	return nil
}

func (p *VerifiRepo) CountNotVerifFiles() (int64, error) {

	var count int64
	if err := p.db.Debug().Table("files").Where("bucket = ? AND review = ?", "verification", false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (p *VerifiRepo) GetStatCheque() (*types.Statistics, error) {

	stats := new(types.Statistics)
	if err := p.db.Debug().Raw(`select (select count(*) from cheques) as total_checks,
       (select count(*) from cheques where "check" = 'Принят') as checks_accepted,
       (select count(*) from cheques where "check" = 'Не принят') as checks_rejected,
       (select count(*) from cheques where "check" = 'На проверке') as checks_on_check,
       count(*) as check_orders from request_gift where certificate = '' AND sent = false`).Take(stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

func (p *VerifiRepo) GetAcceptCheques() (types.Cheques, error) {

	cheques := make(types.Cheques, 0)
	if err := p.db.Debug().Table("cheques").Where("check", "Принят").Order("created_at DESC").Limit(50).Find(&cheques).Error; err != nil {
		return nil, err
	}

	return cheques, nil
}

func (p *VerifiRepo) GetAllUsers() (types.UsersClients, error) {

	users := make(types.UsersClients, 0)
	if err := p.db.Debug().Table("users").Select("*").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (p *VerifiRepo) GetUserID(userID string) (*types.Profile, error) {

	user := new(types.Profile)
	if err := p.db.Debug().Table("users").Where("id = ?", userID).
		Preload("Cheques").
		Preload("Products").
		Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *VerifiRepo) GetAllCheques() (types.Cheques, error) {

	cheques := make(types.Cheques, 0)
	if err := p.db.Debug().Table("cheques").Select("*").Find(&cheques).Error; err != nil {
		return nil, err
	}

	return cheques, nil
}

func (p *VerifiRepo) GetCountUsersGift() (*types.CountUsersGift, error) {

	countUG := new(types.CountUsersGift)

	if err := p.db.Debug().Table("users").Select(`
	sum(case when "email" <> '' then 1 else 0 end) count_users_site,
    sum(case when "email" = '' then 1 else 0 end) count_users_bot,
    (select sum(case when "phone" <> '' then 1 else 0 end) from request_gift) as count_gift`).
		Take(countUG).Error; err != nil {
		return nil, err
	}

	return countUG, nil
}

func (p *VerifiRepo) GetAllRequestGift() (types.RequestGifts, error) {

	reqGift := make(types.RequestGifts, 0)

	if err := p.db.Debug().Table("request_gift").Select("*").Find(&reqGift).Error; err != nil {
		return nil, err
	}

	return reqGift, nil
}

func (p *VerifiRepo) AddLogRespSMS(logSMS *types.LogSMS) error {

	if err := p.db.Debug().Table("log_sms").Create(logSMS).Error; err != nil {
		return err
	}

	return nil
}
