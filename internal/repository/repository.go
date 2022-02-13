package repository

import (
	"github.com/whiskey-back/internal/types"
	"gorm.io/gorm"
)

type Common interface {
	GetUserByID(userID string) (*types.User, error)
	GetGifts() (types.Gifts, error)
}

type Verification interface {
	UpdateRequestGift(gift *types.RequestGift) (*types.RequestGift, error)
	UpdateRequestGiftTrue(gift *types.RequestGift) error
	GetUserByID(userID string) (*types.User, error)
	GetGift(giftID int) (*types.Gift, error)
	CreateCheque(check *types.Cheque) error
	CountCheques(check *types.Cheque) (int64, error)
	UpdateBalanceCheck(userID string, balance float64) error
	GetStatCheque() (*types.Statistics, error)
	GetAcceptCheques() (types.Cheques, error)
	CountNotVerifFiles() (int64, error)
	GetAllUsers() (types.UsersClients, error)
	GetUserID(userID string) (*types.Profile, error)
	GetAllCheques() (types.Cheques, error)
	GetCountUsersGift() (*types.CountUsersGift, error)
	GetAllRequestGift() (types.RequestGifts, error)
	AddLogRespSMS(logSMS *types.LogSMS) error
}

type Profile interface {
	Tx() *gorm.DB
	GetProfile(userID string) (*types.Profile, error)
	GetGift(giftID int) (*types.Gift, error)
	UploadFile(file *types.File) error
	CreateCheque(check *types.Cheque) error
	CreateLoggerReqCheck(log *types.LoggerReqCheck) error
	CreateLoggerRespCheck(log *types.LoggerRespCheck) error
	CreateGoodRespCheck(check *types.GoodRespCheck) error
	CountGoodRespCheck(check *types.JSON) (int64, error)
	CountCheques(check *types.Cheque) (int64, error)
	CreatePositionItemsOnCheck(position *types.PositionInCheck) error
	GetAllWhiskey() ([]types.Whiskey, error)
	UpdateMarkerPositionItemsOnCheck(goodRespCheck, name string) error
	UpdateBalanceCheck(userID string, balance float64) error
	CountFiles(name, userID string) (int64, error)
	GetAllShops() ([]types.Shop, error)
	GetUserByID(userID string) (*types.User, error)
	UpdateBalancePrize(tx *gorm.DB, user *types.User) error
	CreateProduct(tx *gorm.DB, prod *types.Product) error
	CreateRequestGift(tx *gorm.DB, reqGift *types.RequestGift) error
	GetFileByID(photoID string) (*types.File, error)
	UpdateCheckFileTrue(photoID string) error
	GetAllNotVerifFiles() (types.Files, error)
	CountNotVerifFiles() (int64, error)
	CountAdminCheques(photoID string) (int64, error)
	GetAllCities() ([]types.Cities, error)
}

type Authorization interface {
	GetUserByPhone(phone string) (*types.User, error)
	CreateUser(user *types.User) error
	PreCreateUser(user *types.CheckUser) error
	DeleteCodePhone(userID string) error
	AddCodePhone(code string, userID string) error
	PreDeleteCodePhone(userID string) error
	PreAddCodePhone(code string, userID string) error
	GetCheckUser(phone, code string) (*types.CheckUser, error)
	DeleteCheckUser(preUser *types.CheckUser) error
	GetAuthUser(phone, code string) (*types.User, error)
	AddLogRespSMS(logSMS *types.LogSMS) error
}

type Repository struct {
	Verification
	Authorization
	Common
	Profile
}

func NewRepositories(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthorizationRepo(db),
		Common:        NewCommonRepo(db),
		Profile:       NewProfileRepo(db),
		Verification:  NewVerifiRepo(db),
	}
}
