package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/logger"
	"github.com/whiskey-back/pkg/minio"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ProfileService struct {
	repo      repository.Profile
	minio     *minio.FileStorage
	path      string
	check     *config.Check
	dataStart string
}

func NewProfileService(cfg *config.Config, repo repository.Profile, minio *minio.FileStorage) *ProfileService {
	return &ProfileService{
		repo:      repo,
		minio:     minio,
		path:      cfg.PathUrl,
		check:     cfg.Check,
		dataStart: cfg.DateStart,
	}
}

func (s *ProfileService) GetProfile(userID string) (*types.Profile, error) {

	profile, err := s.repo.GetProfile(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetProfile in GetProfile"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, product := range profile.Products {
		profile.Products[i].Gift, err = s.repo.GetGift(product.GiftID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetGift in GetProfile"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	return profile, nil
}

func (s *ProfileService) UploadCheck(file *multipart.FileHeader, userID string) error {

	countFiles, err := s.repo.CountFiles(file.Filename, userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CountFiles in UploadCheck"))
		return infrastruct.ErrorInternalServerError
	}

	if countFiles >= 1 {
		return infrastruct.ErrorDuplicate
	}

	object, err := s.minio.Add(file, s.minio.Bucket, userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Add in UploadCheck"))
		return infrastruct.ErrorInternalServerError
	}

	object.Url = s.getURLFile(s.path, object.Bucket, userID, object.Name)

	if err = s.repo.UploadFile(object); err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFile in UploadCheck"))
		return infrastruct.ErrorInternalServerError
	}

	checkingCheck, code, err := s.checkingCheckPhotoApi(file, userID)
	if err != nil {
		if err == infrastruct.ErrorServiceUnavailable {
			logger.LogError(errors.Wrap(err, "check service unavailable"))
			return err
		}

		logger.LogError(errors.Wrap(err, "err with checkingCheckPhotoApi in UploadCheck"))
		return err
	}

	switch code {
	case 1:

		have, err := s.checkShops(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with checkShops in UploadCheck"))
			return infrastruct.ErrorInternalServerError
		}

		if have == false {

			if err := s.addForChecking(file, userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with addForChecking in UploadCheck"))
				return err
			}

			return nil
		}

		isHave, err := s.checkCities(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with checkCities in UploadCheck"))
			return infrastruct.ErrorInternalServerError
		}

		if isHave == false {

			if err := s.addForChecking(file, userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with addForChecking in UploadCheck"))
				return err
			}

			return nil
		}

		dateCheck, err := time.Parse("2006-01-02T15:04:05", checkingCheck.Data.JSON.Datetime)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with time.Parse dateCheck in UploadCheck"))
			return infrastruct.ErrorInternalServerError
		}

		dateStart, err := time.Parse("2006-01-02T15:04:05", s.dataStart)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with time.Parse dateStart in UploadCheck"))
			return infrastruct.ErrorInternalServerError
		}

		if dateCheck.Unix() < dateStart.Unix() {

			if err := s.addForChecking(file, userID); err != nil {
				logger.LogError(errors.Wrap(err, "err with addForChecking in UploadCheck"))
				return err
			}

			return nil
		}

		count, err := s.repo.CountGoodRespCheck(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CountGoodRespCheck in UploadCheck"))
			return infrastruct.ErrorInternalServerError
		}

		if count < 1 {

			goodRespCheck := &types.GoodRespCheck{
				Code:                    checkingCheck.Data.JSON.Code,
				User:                    checkingCheck.Data.JSON.User,
				FnsUrl:                  checkingCheck.Data.JSON.Fnsurl,
				KktRegID:                checkingCheck.Data.JSON.Kktregid,
				RetailPlace:             checkingCheck.Data.JSON.Retailplace,
				RetailPlaceAddress:      checkingCheck.Data.JSON.Retailplaceaddress,
				UserInn:                 checkingCheck.Data.JSON.Userinn,
				DateTime:                checkingCheck.Data.JSON.Datetime,
				RequestNumber:           checkingCheck.Data.JSON.Requestnumber,
				TotalSum:                checkingCheck.Data.JSON.Totalsum,
				ShiftNumber:             checkingCheck.Data.JSON.Shiftnumber,
				OperationType:           checkingCheck.Data.JSON.Operationtype,
				FiscalDriveNumber:       checkingCheck.Data.JSON.Fiscaldrivenumber,
				FiscalDocumentNumber:    checkingCheck.Data.JSON.Fiscaldocumentnumber,
				FiscalSign:              checkingCheck.Data.JSON.Fiscalsign,
				FiscalDocumentFormatVer: checkingCheck.Data.JSON.Fiscaldocumentformatver,
				Url:                     fmt.Sprintf("https://proverkacheka.com/check/%s-%d-%d", checkingCheck.Data.JSON.Fiscaldrivenumber, checkingCheck.Data.JSON.Fiscaldocumentnumber, checkingCheck.Data.JSON.Fiscalsign),
				UserID:                  userID,
			}

			if err = s.repo.CreateGoodRespCheck(goodRespCheck); err != nil {
				logger.LogError(errors.Wrap(err, "err with CreateGoodRespCheck in UploadCheck"))
				return infrastruct.ErrorInternalServerError
			}

			for i, _ := range checkingCheck.Data.JSON.Items {
				position := &types.PositionInCheck{
					Name:            checkingCheck.Data.JSON.Items[i].Name,
					Price:           checkingCheck.Data.JSON.Items[i].Price,
					Count:           checkingCheck.Data.JSON.Items[i].Quantity,
					Sum:             checkingCheck.Data.JSON.Items[i].Sum,
					GoodRespCheckID: goodRespCheck.ID,
					UserID:          userID,
				}

				if err := s.repo.CreatePositionItemsOnCheck(position); err != nil {
					logger.LogError(errors.Wrap(err, "err with CreatePositionItemsOnCheck in UploadCheck"))
					return infrastruct.ErrorInternalServerError
				}
			}

			totalSum, err := s.findingPositionInCheck(checkingCheck.Data.JSON.Items, goodRespCheck.ID, userID)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with findingPositionInCheck in UploadCheck"))
				return err
			}

			if totalSum != 0 {

				check := new(types.Cheque)
				check.Check = types.Accepted
				check.Winning = totalSum
				check.CheckAmount = strconv.Itoa(checkingCheck.Data.JSON.Totalsum)
				check.Date = checkingCheck.Data.JSON.Datetime
				check.FN = checkingCheck.Data.JSON.Fiscaldrivenumber
				check.FD = strconv.Itoa(checkingCheck.Data.JSON.Fiscaldocumentnumber)
				check.FP = strconv.Itoa(checkingCheck.Data.JSON.Fiscalsign)

				check.UserID = userID
				check.PhotoID = object.ID
				if err = s.repo.CreateCheque(check); err != nil {
					logger.LogError(errors.Wrap(err, "err with CreateCheque in UploadCheck"))
					return infrastruct.ErrorInternalServerError
				}

			} else {

				if err := s.addForChecking(file, userID); err != nil {
					logger.LogError(errors.Wrap(err, "err with addForChecking in UploadCheck"))
					return err
				}

			}

		} else {
			return infrastruct.ErrorDuplicate
		}

	default:
		logger.LogInfo(fmt.Sprintf("url:%s, id:%s, userID:%s", object.Url, object.ID, userID))

		if err := s.addForChecking(file, userID); err != nil {
			logger.LogError(errors.Wrap(err, "err with addForChecking in UploadCheck"))
			return err
		}
	}

	return nil
}

func (s *ProfileService) getURLFile(path, bucket, userID, name string) string {
	return fmt.Sprintf(path, bucket, userID, name)
}

func (s *ProfileService) HandWriteCheck(check *types.Cheque) (*types.Cheque, error) {

	dateCheck, err := time.Parse("2006-01-02T15:04", check.Date)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with time.Parse in HandWriteCheck"))
		return nil, err
	}

	dateStart, err := time.Parse("2006-01-02T15:04:05", s.dataStart)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with time.Parse dateStart in HandWriteCheck"))
		return nil, err
	}

	if dateCheck.Unix() < dateStart.Unix() {
		return nil, infrastruct.ErrorDate
	}

	countCheques, err := s.repo.CountCheques(check)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateGoodRespCheck in HandWriteCheck"))
		return nil, err
	}

	if countCheques >= 1 {
		return nil, infrastruct.ErrorDuplicate
	}

	checkingCheck, code, err := s.checkingCheckApi(check)
	if err != nil {
		if err == infrastruct.ErrorServiceUnavailable {
			logger.LogError(errors.Wrap(err, "check service unavailable"))
			return nil, err
		}

		logger.LogError(errors.Wrap(err, "err with checkingCheckApi in HandWriteCheck"))
		return nil, err
	}

	switch code {
	case 1:

		have, err := s.checkShops(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with checkShops in HandWriteCheck"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if have == false {
			return nil, infrastruct.ErrorShop
		}

		isHave, err := s.checkCities(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with checkCities in HandWriteCheck"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if isHave == false {
			return nil, infrastruct.ErrorCities
		}

		count, err := s.repo.CountGoodRespCheck(&checkingCheck.Data.JSON)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CountGoodRespCheck in HandWriteCheck"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if count < 1 {

			goodRespCheck := &types.GoodRespCheck{
				Code:                    checkingCheck.Data.JSON.Code,
				User:                    checkingCheck.Data.JSON.User,
				FnsUrl:                  checkingCheck.Data.JSON.Fnsurl,
				KktRegID:                checkingCheck.Data.JSON.Kktregid,
				RetailPlace:             checkingCheck.Data.JSON.Retailplace,
				RetailPlaceAddress:      checkingCheck.Data.JSON.Retailplaceaddress,
				UserInn:                 checkingCheck.Data.JSON.Userinn,
				DateTime:                checkingCheck.Data.JSON.Datetime,
				RequestNumber:           checkingCheck.Data.JSON.Requestnumber,
				TotalSum:                checkingCheck.Data.JSON.Totalsum,
				ShiftNumber:             checkingCheck.Data.JSON.Shiftnumber,
				OperationType:           checkingCheck.Data.JSON.Operationtype,
				FiscalDriveNumber:       checkingCheck.Data.JSON.Fiscaldrivenumber,
				FiscalDocumentNumber:    checkingCheck.Data.JSON.Fiscaldocumentnumber,
				FiscalSign:              checkingCheck.Data.JSON.Fiscalsign,
				FiscalDocumentFormatVer: checkingCheck.Data.JSON.Fiscaldocumentformatver,
				Url:                     fmt.Sprintf("https://proverkacheka.com/check/%s-%d-%d", checkingCheck.Data.JSON.Fiscaldrivenumber, checkingCheck.Data.JSON.Fiscaldocumentnumber, checkingCheck.Data.JSON.Fiscalsign),
				UserID:                  check.UserID,
			}

			if err = s.repo.CreateGoodRespCheck(goodRespCheck); err != nil {
				logger.LogError(errors.Wrap(err, "err with CreateGoodRespCheck in HandWriteCheck"))
				return nil, infrastruct.ErrorInternalServerError
			}

			for i, _ := range checkingCheck.Data.JSON.Items {
				position := &types.PositionInCheck{
					Name:            checkingCheck.Data.JSON.Items[i].Name,
					Price:           checkingCheck.Data.JSON.Items[i].Price,
					Count:           checkingCheck.Data.JSON.Items[i].Quantity,
					Sum:             checkingCheck.Data.JSON.Items[i].Sum,
					GoodRespCheckID: goodRespCheck.ID,
					UserID:          check.UserID,
				}

				if err := s.repo.CreatePositionItemsOnCheck(position); err != nil {
					logger.LogError(errors.Wrap(err, "err with CreatePositionItemsOnCheck in HandWriteCheck"))
					return nil, infrastruct.ErrorInternalServerError
				}
			}

			totalSum, err := s.findingPositionInCheck(checkingCheck.Data.JSON.Items, goodRespCheck.ID, check.UserID)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with findingPositionInCheck in HandWriteCheck"))
				return nil, infrastruct.ErrorInternalServerError
			}

			if totalSum != 0 {
				check.Check = types.Accepted
				check.Winning = totalSum
				check.PhotoID = "Ручной ввод"
			} else {
				return nil, infrastruct.ErrorCheckNotValid
			}

		} else {
			return nil, infrastruct.ErrorDuplicate
		}

	default:
		return nil, infrastruct.ErrorCheckNotValid
	}

	if err = s.repo.CreateCheque(check); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateCheque in HandWriteCheck"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return check, nil
}

func (s *ProfileService) checkingCheckApi(check *types.Cheque) (*types.Check, int, error) {

	payload := fmt.Sprintf("fn=%s&fd=%s&fp=%s&n=1&s=%s&t=%s&qr=0&token=%s",
		check.FN,
		check.FD,
		check.FP,
		check.CheckAmount,
		check.Date,
		s.check.Token,
	)

	logg := &types.LoggerReqCheck{
		Logger: payload,
		UserID: check.UserID,
	}

	if err := s.repo.CreateLoggerReqCheck(logg); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateLoggerReqCheck in checkingCheckApi"))
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, s.check.Url, strings.NewReader(payload))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest in checkingCheckApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with client.Do in checkingCheckApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ReadAll in checkingCheckApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	if err = s.repo.CreateLoggerRespCheck(&types.LoggerRespCheck{
		Logger:           string(body),
		UserID:           check.UserID,
		LoggerReqCheckID: logg.ID,
	}); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateLoggerRespCheck in checkingCheckApi"))
	}

	if res.StatusCode != 200 {
		return nil, 0, infrastruct.ErrorServiceUnavailable
	}

	preCheck := new(types.PreloadAnswer)
	err = json.Unmarshal(body, preCheck)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode in checkingCheckApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	checkingCheck := new(types.Check)
	if preCheck.Code == 1 {
		err = json.Unmarshal(body, checkingCheck)
		if err != nil {
			logger.LogError(errors.New(fmt.Sprintf("err with Decode checkingCheck in checkingCheckApi, AnswerAPI: %s", string(body))))
			return nil, 0, infrastruct.ErrorInternalServerError
		}
	}

	return checkingCheck, preCheck.Code, nil
}

func (s *ProfileService) findingPositionInCheck(items []types.Items, goodRespCheck, userID string) (float64, error) {

	whiskey, err := s.repo.GetAllWhiskey()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllWhiskey in findingPositionInCheck"))
		return 0, infrastruct.ErrorInternalServerError
	}

	var totalSumWinning float64

	for i := range items {
		for j := range whiskey {
			if items[i].Name == whiskey[j].Name || strings.Contains(items[i].Name, whiskey[j].Name) {

				if err := s.repo.UpdateMarkerPositionItemsOnCheck(goodRespCheck, items[i].Name); err != nil {
					logger.LogError(errors.Wrap(err, "err with UpdateMarkerPositionItemsOnCheck in findingPositionInCheck"))
				}

				price := items[i].Price
				count := items[i].Quantity
				price = price * count
				price = price / 100

				user, err := s.repo.GetUserByID(userID)
				if err != nil {
					logger.LogError(errors.Wrap(err, "err with GetProfile in findingPositionInCheck"))
					return 0, infrastruct.ErrorInternalServerError
				}

				totalSumWinning = totalSumWinning + price
				price = price + user.Balance

				if err = s.repo.UpdateBalanceCheck(user.ID, price); err != nil {
					logger.LogError(errors.Wrap(err, "err with UpdateBalance in findingPositionInCheck"))
					return 0, infrastruct.ErrorInternalServerError
				}

				break
			}
		}
	}

	return totalSumWinning, nil
}

func (s *ProfileService) checkingCheckPhotoApi(file *multipart.FileHeader, userID string) (*types.Check, int, error) {

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	f, err := file.Open()
	defer f.Close()

	part, err := writer.CreateFormFile("qrfile", file.Filename)
	_, err = io.Copy(part, f)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateFormFile in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	_ = writer.WriteField("token", s.check.Token)
	err = writer.Close()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with WriteField in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, s.check.Url, payload)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Header.Set in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with ReadAll in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	if err = s.repo.CreateLoggerRespCheck(&types.LoggerRespCheck{
		Logger:           string(body),
		UserID:           userID,
		LoggerReqCheckID: "Photo",
	}); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateLoggerRespCheck in checkingCheckPhotoApi"))
	}

	if res.StatusCode != 200 {
		return nil, 0, infrastruct.ErrorServiceUnavailable
	}

	preCheck := new(types.PreloadAnswer)
	err = json.Unmarshal(body, preCheck)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Decode in checkingCheckPhotoApi"))
		return nil, 0, infrastruct.ErrorInternalServerError
	}

	checkingCheck := new(types.Check)
	if preCheck.Code == 1 {
		err = json.Unmarshal(body, checkingCheck)
		if err != nil {
			logger.LogError(errors.New(fmt.Sprintf("err with Decode checkingCheck in checkingCheckPhotoApi, AnswerAPI: %s", string(body))))
			return nil, 0, infrastruct.ErrorInternalServerError
		}
	}

	return checkingCheck, preCheck.Code, nil
}

func (s *ProfileService) checkShops(shop *types.JSON) (bool, error) {

	shops, err := s.repo.GetAllShops()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllShops in checkShops"))
		return false, infrastruct.ErrorInternalServerError
	}

	for i, _ := range shops {
		if shop.Userinn == shops[i].Inn || strings.Contains(shop.Userinn, shops[i].Inn) {
			return true, nil
		}
	}

	return false, err
}

func (s *ProfileService) PrizeLogic(prize *types.Prize, userID string) (types.Products, error) {
	tx := s.repo.Tx()
	defer tx.Rollback()

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in PrizeLogic"))
		return nil, infrastruct.ErrorInternalServerError
	}

	products := make(types.Products, 0)

	for _, product := range prize.Products {

		gift, err := s.repo.GetGift(product.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			logger.LogError(errors.Wrap(err, "err with GetGift in PrizeLogic"))
			return nil, infrastruct.ErrorInternalServerError
		}

		if gift == nil {
			return nil, infrastruct.ErrorWrongIDGift
		}

		if product.Count == 0 {
			return nil, infrastruct.ErrorGiftCountZero
		}

		user.Balance = user.Balance - (gift.Price * float64(product.Count))
		if user.Balance < 0 {
			return nil, infrastruct.ErrorBalanceEmpty
		}

		user.Action = user.Action - (gift.Sum * float64(product.Count))
		if user.Action < 0 {
			return nil, infrastruct.ErrorTotalLimit
		}

		if err = s.repo.UpdateBalancePrize(tx, user); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateBalancePrize in PrizeLogic"))
			return nil, infrastruct.ErrorInternalServerError
		}

		prod := &types.Product{
			GiftID: gift.ID,
			Count:  int(product.Count),
			UserID: user.ID,
		}

		if err = s.repo.CreateProduct(tx, prod); err != nil {
			logger.LogError(errors.Wrap(err, "err with CreateProduct in PrizeLogic"))
			return nil, infrastruct.ErrorInternalServerError
		}

		products = append(products, *prod)

		for i := 0; i < int(product.Count); i++ {
			if err = s.repo.CreateRequestGift(tx, &types.RequestGift{
				Phone:     user.Phone,
				GiftID:    gift.ID,
				PrizeName: gift.Name,
				ProductID: prod.ID,
				UserID:    user.ID,
			}); err != nil {
				logger.LogError(errors.Wrap(err, "err with CreateRequestGift in PrizeLogic"))
				return nil, infrastruct.ErrorInternalServerError
			}
		}
	}

	if err = tx.Commit().Error; err != nil {
		logger.LogError(errors.Wrap(err, "err with tx.Commit() in PrizeLogic"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i, product := range products {
		products[i].Gift, err = s.repo.GetGift(product.GiftID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetGift in GetProfile"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	return products, nil
}

func (s *ProfileService) GetAllNotVerifiPhoto(userID string) (types.Files, error) {

	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in GetAllNotVerifiPhoto"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user.Role != "ADMIN" {
		return nil, infrastruct.ErrorPermissionDenied
	}

	files, err := s.repo.GetAllNotVerifFiles()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllNotVerifFiles in GetAllNotVerifiPhoto"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return files, nil
}

func (s *ProfileService) VerifiPhoto(cheque *types.Cheque, photoID string, userID string) error {

	check, err := s.repo.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in VerifiPhoto"))
		return infrastruct.ErrorInternalServerError
	}

	if check.Role != "ADMIN" {
		return infrastruct.ErrorPermissionDenied
	}

	if cheque.FP != "" && cheque.FD != "" && cheque.FN != "" {
		countCheques, err := s.repo.CountCheques(cheque)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with CountCheques in VerifiPhoto"))
			return err
		}

		if countCheques >= 1 {
			return infrastruct.ErrorDuplicate
		}
	}

	countCheques, err := s.repo.CountAdminCheques(photoID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CountAdminCheques in VerifiPhoto"))
		return err
	}

	if countCheques >= 1 {
		return infrastruct.ErrorDuplicate
	}

	switch cheque.Check {
	case types.NotAccepted:

		if err = s.forChecking(cheque, photoID); err != nil {
			logger.LogError(errors.Wrap(err, "err with forChecking in VerifiPhoto"))
			return err
		}

	case types.Accepted:

		checkingCheck, code, err := s.checkingCheckApi(cheque)
		if err != nil {
			if err == infrastruct.ErrorServiceUnavailable {
				logger.LogError(errors.Wrap(err, "check service unavailable"))
				return err
			}

			logger.LogError(errors.Wrap(err, "err with checkingCheckApi in VerifiPhoto"))
			return err
		}

		switch code {
		case 1:
			count, err := s.repo.CountGoodRespCheck(&checkingCheck.Data.JSON)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with CountGoodRespCheck in VerifiPhoto"))
				return infrastruct.ErrorInternalServerError
			}

			if count <= 1 {

				goodRespCheck := &types.GoodRespCheck{
					Code:                    checkingCheck.Data.JSON.Code,
					User:                    checkingCheck.Data.JSON.User,
					FnsUrl:                  checkingCheck.Data.JSON.Fnsurl,
					KktRegID:                checkingCheck.Data.JSON.Kktregid,
					RetailPlace:             checkingCheck.Data.JSON.Retailplace,
					RetailPlaceAddress:      checkingCheck.Data.JSON.Retailplaceaddress,
					UserInn:                 checkingCheck.Data.JSON.Userinn,
					DateTime:                checkingCheck.Data.JSON.Datetime,
					RequestNumber:           checkingCheck.Data.JSON.Requestnumber,
					TotalSum:                checkingCheck.Data.JSON.Totalsum,
					ShiftNumber:             checkingCheck.Data.JSON.Shiftnumber,
					OperationType:           checkingCheck.Data.JSON.Operationtype,
					FiscalDriveNumber:       checkingCheck.Data.JSON.Fiscaldrivenumber,
					FiscalDocumentNumber:    checkingCheck.Data.JSON.Fiscaldocumentnumber,
					FiscalSign:              checkingCheck.Data.JSON.Fiscalsign,
					FiscalDocumentFormatVer: checkingCheck.Data.JSON.Fiscaldocumentformatver,
					Url:                     fmt.Sprintf("https://proverkacheka.com/check/%s-%d-%d", checkingCheck.Data.JSON.Fiscaldrivenumber, checkingCheck.Data.JSON.Fiscaldocumentnumber, checkingCheck.Data.JSON.Fiscalsign),
					UserID:                  cheque.UserID,
				}

				if err = s.repo.CreateGoodRespCheck(goodRespCheck); err != nil {
					logger.LogError(errors.Wrap(err, "err with CreateGoodRespCheck in VerifiPhoto"))
					return infrastruct.ErrorInternalServerError
				}

				for i, _ := range checkingCheck.Data.JSON.Items {
					position := &types.PositionInCheck{
						Name:            checkingCheck.Data.JSON.Items[i].Name,
						Price:           checkingCheck.Data.JSON.Items[i].Price,
						Count:           checkingCheck.Data.JSON.Items[i].Quantity,
						Sum:             checkingCheck.Data.JSON.Items[i].Sum,
						GoodRespCheckID: goodRespCheck.ID,
						UserID:          cheque.UserID,
					}

					if err := s.repo.CreatePositionItemsOnCheck(position); err != nil {
						logger.LogError(errors.Wrap(err, "err with CreatePositionItemsOnCheck in VerifiPhoto"))
						return infrastruct.ErrorInternalServerError
					}
				}

				totalSum, err := s.findingPositionInCheck(checkingCheck.Data.JSON.Items, goodRespCheck.ID, cheque.UserID)
				if err != nil {
					logger.LogError(errors.Wrap(err, "err with findingPositionInCheck in VerifiPhoto"))
					return infrastruct.ErrorInternalServerError
				}

				if totalSum != 0 {

					cheque.Check = types.Accepted
					cheque.Winning = totalSum

				} else {

					cheque.Check = types.NotAccepted

				}

			} else {

				return infrastruct.ErrorDuplicate

			}

		default:

			cheque.Check = types.NotAccepted

		}

		if err = s.forChecking(cheque, photoID); err != nil {
			logger.LogError(errors.Wrap(err, "err with forChecking in VerifiPhoto"))
			return err
		}
	}

	return nil
}

func (s *ProfileService) checkCities(shop *types.JSON) (bool, error) {

	cities, err := s.repo.GetAllCities()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetAllCities in checkCities"))
		return false, infrastruct.ErrorInternalServerError
	}

	for i, _ := range cities {
		if strings.Contains(shop.Retailplaceaddress, cities[i].City) {
			return true, nil
		}
	}

	return false, err
}

func (s *ProfileService) addForChecking(file *multipart.FileHeader, userID string) error {

	object, err := s.minio.Add(file, s.minio.BadBucket, userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with addForChecking in badBucket addForChecking"))
		return infrastruct.ErrorInternalServerError
	}

	object.Url = s.getURLFile(s.path, object.Bucket, userID, object.Name)

	if err = s.repo.UploadFile(object); err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFile in addForChecking"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *ProfileService) forChecking(cheque *types.Cheque, photoID string) error {

	file, err := s.repo.GetFileByID(photoID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetFileByID in forChecking"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.minio.Delete(file, s.minio.BadBucket); err != nil {
		logger.LogError(errors.Wrap(err, "err with Delete in forChecking"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.repo.UpdateCheckFileTrue(file.ID); err != nil {
		logger.LogError(errors.Wrap(err, "err with UpdateCheckFileTrue in forChecking"))
		return infrastruct.ErrorInternalServerError
	}

	cheque.PhotoID = photoID
	if err = s.repo.CreateCheque(cheque); err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateCheque in forChecking"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
