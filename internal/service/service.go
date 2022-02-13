package service

import (
	"mime/multipart"

	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/jwt"
	"github.com/whiskey-back/pkg/minio"
	"github.com/whiskey-back/pkg/sms"
	"github.com/whiskey-back/pkg/yoomoney"
)

type Common interface {
	CheckUsers(userID string) error
	GetGifts() (types.Gifts, error)
	SupportSend(sup *types.Support) error
}

type Verification interface {
	UploadCSV(file *multipart.FileHeader) error
	Statistics(userID string) (*types.Statistics, error)
	RecentCheques(userID string) (types.Cheques, error)
	GetAllUsers(userID string) (types.UsersClients, error)
	GetUser(userID string, id string) (*types.Profile, error)
	GetAllCheques(userID string) (types.Cheques, error)
	GetCountUsersGift(userID string) (*types.CountUsersGift, error)
	GetAllRequestGift(userID string) (types.RequestGifts, error)
}

type Profile interface {
	GetProfile(userID string) (*types.Profile, error)
	UploadCheck(file *multipart.FileHeader, userID string) error
	HandWriteCheck(check *types.Cheque) (*types.Cheque, error)
	PrizeLogic(prize *types.Prize, userID string) (types.Products, error)
	GetAllNotVerifiPhoto(userID string) (types.Files, error)
	VerifiPhoto(cheque *types.Cheque, photoID string, userID string) error
}

type Authorization interface {
	PreRegistrationUser(newUser *types.CheckUser) error
	RegistrationUser(user *types.User) (*types.Token, error)
	PreAuthorizationUser(user *types.User) error
	AuthorizationUser(user *types.User) (*types.Token, error)
}

type Service struct {
	Verification
	Authorization
	Common
	Profile
}

func NewServices(
	cfg *config.Config,
	db *repository.Repository,
	tokenManager *jwt.Manager,
	apiSMS *sms.ApiSMS,
	newMinio *minio.FileStorage,
	yooMoney *yoomoney.Yoomoney,
) *Service {
	return &Service{
		Authorization: NewAuthorizationService(cfg, db.Authorization, tokenManager, apiSMS),
		Common:        NewCommonService(db.Common),
		Profile:       NewProfileService(cfg, db.Profile, newMinio),
		Verification:  NewVerifiService(db.Verification, apiSMS, yooMoney, newMinio),
	}
}
