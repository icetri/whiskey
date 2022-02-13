package service

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/logger"
	"github.com/whiskey-back/pkg/minio"
	"github.com/whiskey-back/pkg/sms"
	"github.com/whiskey-back/pkg/yoomoney"
	"mime/multipart"
)

type VerifiService struct {
	repo     repository.Verification
	apiSMS   *sms.ApiSMS
	yooMoney *yoomoney.Yoomoney
	minio    *minio.FileStorage
}

func NewVerifiService(repo repository.Verification, apiSMS *sms.ApiSMS, yooMoney *yoomoney.Yoomoney, minio *minio.FileStorage) *VerifiService {
	return &VerifiService{
		repo:     repo,
		apiSMS:   apiSMS,
		yooMoney: yooMoney,
		minio:    minio,
	}
}

func (s *VerifiService) UploadCSV(file *multipart.FileHeader) error {

	f, err := file.Open()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetGifts in UploadCSV"))
		return infrastruct.ErrorInternalServerError
	}
	defer f.Close()

	gifts := make([]types.ForCSVTable, 0)
	reader := gocsv.DefaultCSVReader(f)
	if err := gocsv.UnmarshalCSVWithoutHeaders(reader, &gifts); err != nil {
		logger.LogError(errors.Wrap(err, "err with gocsv.Unmarshal in UploadCSV"))
		return infrastruct.ErrorInternalServerError
	}

	for _, gift := range gifts {

		reqGift, err := s.repo.UpdateRequestGift(&types.RequestGift{
			Certificate: gift.Certificate,
			ID:          gift.ID,
		})
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateRequestGift in UploadCSV"))
			return infrastruct.ErrorInternalServerError
		}

		user, err := s.repo.GetUserByID(reqGift.UserID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserByID in UploadCSV"))
			return infrastruct.ErrorInternalServerError
		}

		giftProd, err := s.repo.GetGift(reqGift.GiftID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetGift in UploadCSV"))
			return infrastruct.ErrorInternalServerError
		}

		switch reqGift.GiftID {
		case 1:
			if err := s.yooMoney.SendYoomoney(user, fmt.Sprintf("%.2f", giftProd.Sum)); err != nil {
				logger.LogError(errors.Wrap(err, "err with SendYoomoney in UploadCSV"))
				return err
			}

			_, err = s.repo.UpdateRequestGift(&types.RequestGift{
				Certificate: "Зачислен",
				ID:          gift.ID,
			})
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateRequestGift in UploadCSV"))
				return infrastruct.ErrorInternalServerError
			}

		default:
			mess := make([]types.SMSMessages, 0)
			mes := types.SMSMessages{Phone: user.Phone, Sender: s.apiSMS.Sender, ClientID: user.ID, Text: gift.Certificate}
			mess = append(mess, mes)

			resp, err := s.apiSMS.SendSMS(mess)
			if err != nil {
				return err
			}

			if err = s.repo.AddLogRespSMS(&types.LogSMS{
				Resp:   string(resp),
				UserID: user.ID,
			}); err != nil {
				logger.LogError(errors.Wrap(err, "err with AddLogRespSMS in UploadCSV"))
			}
		}

		reqGift.Send = true
		if err := s.repo.UpdateRequestGiftTrue(reqGift); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateRequestGiftTrue in UploadCSV"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}

func (s *VerifiService) Statistics(userID string) (*types.Statistics, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in Statistics"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	statistics, err := s.repo.GetStatCheque()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetStatCheque in Statistics"))
		return nil, infrastruct.ErrorInternalServerError
	}

	count, err := s.repo.CountNotVerifFiles()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CountNotVerifFiles in Statistics"))
		return nil, infrastruct.ErrorInternalServerError
	}

	statistics.TotalChecks = statistics.TotalChecks + count
	statistics.ChecksOnCheck = statistics.ChecksOnCheck + count

	return statistics, nil
}

func (s *VerifiService) RecentCheques(userID string) (types.Cheques, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in RecentCheques"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	cheques, err := s.repo.GetAcceptCheques()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAcceptCheques in RecentCheques"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return cheques, nil
}

func (s *VerifiService) GetAllUsers(userID string) (types.UsersClients, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetAllUsers"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	us, err := s.repo.GetAllUsers()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllUsers in GetAllUsers"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, user := range us {
		if user.Email == "" {
			us[i].Client = "Bot"
		} else {
			us[i].Client = "Site"
		}
	}

	return us, nil
}

func (s *VerifiService) GetUser(userID string, id string) (*types.Profile, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	us, err := s.repo.GetUserID(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserID in GetUser"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i := range us.Products {
		us.Products[i].Gift, err = s.repo.GetGift(us.Products[i].GiftID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetGift in GetUser"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	return us, nil
}

func (s *VerifiService) GetAllCheques(userID string) (types.Cheques, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetAllCheques"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	cheques, err := s.repo.GetAllCheques()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllCheques in GetAllCheques"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return cheques, nil
}

func (s *VerifiService) GetCountUsersGift(userID string) (*types.CountUsersGift, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetCountUsGi"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	countUG, err := s.repo.GetCountUsersGift()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetCountUsersGift in GetCountUsGi"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return countUG, nil
}

func (s *VerifiService) GetAllRequestGift(userID string) (types.RequestGifts, error) {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetAllRequestGift"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	allRequestGift, err := s.repo.GetAllRequestGift()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllRequestGift in GetAllRequestGift"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return allRequestGift, nil
}
