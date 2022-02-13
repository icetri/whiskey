package repository

import (
	"github.com/whiskey-back/internal/types"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	db *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{
		db: db,
	}
}

func (p *ProfileRepo) GetProfile(userID string) (*types.Profile, error) {

	profile := new(types.Profile)
	if err := p.db.Debug().Table("users").Where("id = ?", userID).
		Preload("Cheques").
		Preload("Products").
		Take(profile).
		Error; err != nil {
		return nil, err
	}

	return profile, nil
}

func (p *ProfileRepo) GetGift(giftID int) (*types.Gift, error) {

	gift := new(types.Gift)
	if err := p.db.Debug().Table("gifts").Where("id = ?", giftID).Take(gift).Error; err != nil {
		return nil, err
	}

	return gift, nil
}

func (p *ProfileRepo) UploadFile(file *types.File) error {

	if err := p.db.Debug().Table("files").Create(file).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) Tx() *gorm.DB {
	return p.db.Begin()
}

func (p *ProfileRepo) CountFiles(name, userID string) (int64, error) {

	var count int64
	if err := p.db.Debug().Table("files").Where("name = ? AND user_id = ?", name, userID).Count(&count).Error; err != nil {
		return 10, err
	}

	return count, nil
}

//func (p *ProfileRepo) UpdateStatus() error {
//
//	if err := p.db.Debug().Table("files").Create(file).Error; err != nil {
//		return err
//	}
//
//	return nil
//}

func (p *ProfileRepo) CreateLoggerReqCheck(log *types.LoggerReqCheck) error {

	if err := p.db.Debug().Table("logger_req_check").Create(log).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) CreateLoggerRespCheck(log *types.LoggerRespCheck) error {

	if err := p.db.Debug().Table("logger_resp_check").Create(log).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) CreateCheque(check *types.Cheque) error {

	if err := p.db.Debug().Table("cheques").Create(check).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) CreateGoodRespCheck(check *types.GoodRespCheck) error {

	if err := p.db.Debug().Table("good_resp_check").Create(check).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) CountGoodRespCheck(check *types.JSON) (int64, error) {

	var count int64
	if err := p.db.Debug().Table("good_resp_check").
		Where("fiscal_sign = ? AND fiscal_driver_number = ? AND fiscal_document_number = ?", check.Fiscalsign, check.Fiscaldrivenumber, check.Fiscaldocumentnumber).
		Count(&count).Error; err != nil {
		return 10, err
	}

	return count, nil
}

func (p *ProfileRepo) CountCheques(check *types.Cheque) (int64, error) {

	var count int64
	if err := p.db.Debug().Table("cheques").
		Where("fn = ? AND fd = ? AND fp = ?", check.FN, check.FD, check.FP).
		Count(&count).Error; err != nil {
		return 10, err
	}

	return count, nil
}

func (p *ProfileRepo) CreatePositionItemsOnCheck(position *types.PositionInCheck) error {

	if err := p.db.Debug().Table("position_in_check").Create(position).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) GetAllWhiskey() ([]types.Whiskey, error) {

	whiskey := make([]types.Whiskey, 0)
	if err := p.db.Debug().Table("whiskey").Select("*").Find(&whiskey).Error; err != nil {
		return nil, err
	}

	return whiskey, nil
}

func (p *ProfileRepo) UpdateMarkerPositionItemsOnCheck(goodRespCheck, name string) error {

	if err := p.db.Debug().Table("position_in_check").Where("name = ? AND good_resp_check_id = ?", name, goodRespCheck).Update("marker", "FIND").Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) UpdateBalanceCheck(userID string, balance float64) error {

	if err := p.db.Debug().Table("users").Select("balance").Where("id = ?", userID).Updates(&types.Profile{Balance: balance}).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) UpdateBalancePrize(tx *gorm.DB, user *types.User) error {

	if err := tx.Debug().Table("users").Select("balance", "action").Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) GetAllShops() ([]types.Shop, error) {

	shops := make([]types.Shop, 0)
	if err := p.db.Debug().Table("shops").Select("*").Find(&shops).Error; err != nil {
		return nil, err
	}

	return shops, nil
}

func (p *ProfileRepo) GetUserByID(userID string) (*types.User, error) {

	user := new(types.User)
	if err := p.db.Debug().Table("users").Where("id = ?", userID).Take(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (p *ProfileRepo) CreateProduct(tx *gorm.DB, prod *types.Product) error {

	if err := tx.Debug().Table("products").Create(prod).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) CreateRequestGift(tx *gorm.DB, reqGift *types.RequestGift) error {

	if err := tx.Debug().Table("request_gift").Create(reqGift).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) GetFileByID(photoID string) (*types.File, error) {

	file := new(types.File)
	if err := p.db.Debug().Table("files").Where("id = ?", photoID).Take(file).Error; err != nil {
		return nil, err
	}

	return file, nil
}

func (p *ProfileRepo) UpdateCheckFileTrue(photoID string) error {

	if err := p.db.Debug().Table("files").Where("id = ?", photoID).Update("review", true).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProfileRepo) GetAllNotVerifFiles() (types.Files, error) {

	files := make(types.Files, 0)
	if err := p.db.Debug().Table("files").Where("bucket = ? AND review = ?", "verification", false).Find(&files).Error; err != nil {
		return nil, err
	}

	return files, nil
}

func (p *ProfileRepo) CountNotVerifFiles() (int64, error) {

	var count int64
	if err := p.db.Debug().Table("files").Where("bucket = ? AND review = ?", "verification", false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (p *ProfileRepo) CountAdminCheques(photoID string) (int64, error) {

	var count int64
	if err := p.db.Debug().Table("cheques").Where("photo_id = ?", photoID).Count(&count).Error; err != nil {
		return 10, err
	}

	return count, nil
}

func (p *ProfileRepo) GetAllCities() ([]types.Cities, error) {

	cities := make([]types.Cities, 0)
	if err := p.db.Debug().Table("cities").Select("*").Find(&cities).Error; err != nil {
		return nil, err
	}

	return cities, nil
}
