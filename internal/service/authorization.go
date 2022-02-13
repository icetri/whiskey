package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/jwt"
	"github.com/whiskey-back/pkg/logger"
	"github.com/whiskey-back/pkg/sms"
	"gorm.io/gorm"
)

type AuthService struct {
	JWTKey       string
	repo         repository.Authorization
	tokenManager *jwt.Manager
	apiSMS       *sms.ApiSMS
}

func NewAuthorizationService(cfg *config.Config, repo repository.Authorization, tokenManager *jwt.Manager, apiSMS *sms.ApiSMS) *AuthService {
	return &AuthService{
		repo:         repo,
		tokenManager: tokenManager,
		JWTKey:       cfg.JWTKey,
		apiSMS:       apiSMS,
	}
}

func (s *AuthService) makePhoneCode(user *types.User) error {

	if err := s.repo.DeleteCodePhone(user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteCodePhone in makePhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	rand.Seed(time.Now().UnixNano())
	num := 999 + rand.Intn(9000)
	code := strconv.Itoa(num)

	if err := s.repo.AddCodePhone(code, user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCodePhone in makePhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	mess := make([]types.SMSMessages, 0)
	mes := types.SMSMessages{Phone: user.Phone, Sender: s.apiSMS.Sender, ClientID: user.ID, Text: code}
	mess = append(mess, mes)

	resp, err := s.apiSMS.SendSMS(mess)
	if err != nil {
		return infrastruct.ErrorInternalServerError
	}

	if err = s.repo.AddLogRespSMS(&types.LogSMS{
		Resp:   string(resp),
		UserID: user.ID,
	}); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogRespSMS in makePhoneCode"))
	}

	//		logger.LogInfo(code)

	return nil
}

func (s *AuthService) makePreRegistrationPhoneCode(user *types.CheckUser) error {

	if err := s.repo.PreDeleteCodePhone(user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with PreDeleteCodePhone in makePreRegistrationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	rand.Seed(time.Now().UnixNano())
	num := 999 + rand.Intn(9000)
	code := strconv.Itoa(num)

	if err := s.repo.PreAddCodePhone(code, user.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with PreAddCodePhone in makePreRegistrationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	mess := make([]types.SMSMessages, 0)
	mes := types.SMSMessages{Phone: user.Phone, Sender: s.apiSMS.Sender, ClientID: user.ID, Text: code}
	mess = append(mess, mes)

	resp, err := s.apiSMS.SendSMS(mess)
	if err != nil {
		return err
	}

	if err = s.repo.AddLogRespSMS(&types.LogSMS{
		Resp:   string(resp),
		UserID: user.ID,
	}); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogRespSMS in makePreRegistrationPhoneCode"))
	}

	//		logger.LogInfo(code)

	return nil
}

func (s *AuthService) PreRegistrationUser(newUser *types.CheckUser) error {

	checkUser, err := s.repo.GetUserByPhone(newUser.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in PreRegistrationUser"))
		return infrastruct.ErrorInternalServerError
	}

	if checkUser != nil {
		return infrastruct.ErrorPhoneIsExist
	}

	if err = s.repo.PreCreateUser(newUser); err != nil {
		logger.LogError(errors.Wrap(err, "err with PreCreateUser in PreRegistrationUser"))
		return infrastruct.ErrorInternalServerError
	}

	if err = s.makePreRegistrationPhoneCode(newUser); err != nil {
		logger.LogError(errors.Wrap(err, "err with makePreRegistrationPhoneCode in PreRegistrationUser"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *AuthService) RegistrationUser(user *types.User) (*types.Token, error) {

	checkUser, err := s.repo.GetUserByPhone(user.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in RegistrationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if checkUser != nil {
		return nil, infrastruct.ErrorPhoneIsExist
	}

	preUser, err := s.repo.GetCheckUser(user.Phone, user.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.LogError(errors.Wrap(err, "err with GetCheckUser in RegistrationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if preUser == nil {
		return nil, infrastruct.ErrorIncorrectCode
	}

	if user.Code != preUser.Code {
		return nil, infrastruct.ErrorIncorrectCode
	}

	user.Code = ""
	if err = s.repo.CreateUser(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegistrationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err = s.repo.DeleteCheckUser(preUser); err != nil {
		logger.LogError(errors.Wrap(err, "err with DeleteCheckUser in RegistrationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	token, err := s.tokenManager.GenerateJWT(user.ID, s.JWTKey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateJWT in AuthorizationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{
		Token: token,
	}, nil
}

func (s *AuthService) PreAuthorizationUser(authUser *types.User) error {

	user, err := s.repo.GetUserByPhone(authUser.Phone)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.LogError(errors.Wrap(err, "err with GetUserByPhone in PreAuthorizationUser"))
		return infrastruct.ErrorInternalServerError
	}

	if user == nil {
		return infrastruct.ErrorUserNotExist
	}

	if err = s.makePhoneCode(user); err != nil {
		logger.LogError(errors.Wrap(err, "err with makePhoneCode in PreAuthorizationUser"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *AuthService) AuthorizationUser(authUser *types.User) (*types.Token, error) {

	user, err := s.repo.GetAuthUser(authUser.Phone, authUser.Code)
	if err != nil && err != gorm.ErrRecordNotFound {
		logger.LogError(errors.Wrap(err, "err with GetAuthUser in AuthorizationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user == nil {
		return nil, infrastruct.ErrorIncorrectCode
	}

	if user.Code != authUser.Code {
		return nil, infrastruct.ErrorIncorrectCode
	}

	token, err := s.tokenManager.GenerateJWT(user.ID, s.JWTKey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateJWT in AuthorizationUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{
		Token: token,
	}, nil
}
