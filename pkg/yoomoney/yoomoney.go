package yoomoney

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/logger"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Yoomoney struct {
	db       *gorm.DB
	yoomoney *config.Yoomoney
}

func NewYoomoney(db *gorm.DB, yoomoney *config.Yoomoney) (*Yoomoney, error) {
	return &Yoomoney{
		db:       db,
		yoomoney: yoomoney,
	}, nil
}

type PayRes struct {
	Request_id      string  `json:"request_id"`
	Contract_amount float64 `json:"contract_amount"`
	Title           string  `json:"title"`
	Error           string  `json:"error"`
	Status          string  `json:"status"`
}

type PayedRes struct {
	Status     string `json:"status"`
	Payment_id string `json:"payment_id"`
	Invoice_id string `json:"invoice_id"`
	Error      string `json:"error"`
}

func (y *Yoomoney) SendYoomoney(user *types.User, amount string) error {

	number := strings.TrimPrefix(user.Phone, "+")

	yoomoneyLogs, err := y.urlRequestStart(number, amount, user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with urlRequestStart in SendYoomoney"))
		return err
	}

	err = y.urlRequestFinish(yoomoneyLogs, number, amount, user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with urlRequestFinish in SendYoomoney"))
		return err
	}

	return nil
}

func (y *Yoomoney) urlRequestStart(number string, amount, userID string) (*types.YoomoneyLogs, error) {

	data := url.Values{}
	data.Set("amount", amount)
	data.Set("pattern_id", "phone-topup")
	data.Set("phone-number", number)

	r, err := http.NewRequest(http.MethodPost, y.yoomoney.ReqStart, strings.NewReader(data.Encode()))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.NewRequest in urlRequestStart"))
		return nil, infrastruct.ErrorInternalServerError
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", y.yoomoney.Auth)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with DefaultClient.Do in urlRequestStart"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if res.StatusCode == http.StatusUnauthorized {
		logger.LogError(errors.Wrap(err, "err with StatusUnauthorized in urlRequestStart"))
		return nil, infrastruct.ErrorInternalServerError
	}

	logStart := PayRes{}
	if err = json.NewDecoder(res.Body).Decode(&logStart); err != nil {
		logger.LogError(errors.Wrap(err, "err with NewDecoder in urlRequestStart"))
		return nil, infrastruct.ErrorInternalServerError
	}

	yooMoneyLog := &types.YoomoneyLogs{
		UserID:    userID,
		Type:      "Start",
		Error:     logStart.Error,
		RequestID: logStart.Request_id,
		Status:    logStart.Status,
		Amount:    amount,
		Phone:     number,
	}

	if err = y.db.Debug().Table("yoomoney_log").Create(yooMoneyLog).Error; err != nil {
		logger.LogError(errors.Wrap(err, "err with yoomoney_log in urlRequestStart"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return yooMoneyLog, nil
}

func (y *Yoomoney) urlRequestFinish(yoomoneyLogs *types.YoomoneyLogs, number string, amount, userID string) error {

	data := url.Values{}
	data.Set("request_id", yoomoneyLogs.RequestID)

	r, err := http.NewRequest(http.MethodPost, y.yoomoney.ReqFinish, strings.NewReader(data.Encode()))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err NewRequest in urlRequestFinish"))
		return infrastruct.ErrorInternalServerError
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", y.yoomoney.Auth)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err DefaultClient.Do in urlRequestFinish"))
		return infrastruct.ErrorInternalServerError
	}

	if res.StatusCode == http.StatusUnauthorized {
		logger.LogError(errors.Wrap(err, "err StatusUnauthorized in urlRequestFinish"))
		return infrastruct.ErrorInternalServerError
	}

	logFinal := PayedRes{}
	if err = json.NewDecoder(res.Body).Decode(&logFinal); err != nil {
		logger.LogError(errors.Wrap(err, "err NewDecoder in urlRequestFinish"))
		return infrastruct.ErrorInternalServerError
	}

	if err = y.db.Debug().Table("yoomoney_log").Create(&types.YoomoneyLogs{
		UserID:    userID,
		Type:      "Finish",
		Error:     logFinal.Error,
		RequestID: yoomoneyLogs.RequestID,
		Status:    logFinal.Status,
		Amount:    amount,
		PaymentID: logFinal.Payment_id,
		InvoiceID: logFinal.Invoice_id,
		Phone:     number,
	}).Error; err != nil {
		logger.LogError(errors.Wrap(err, "err yoomoney_log in urlRequestFinish"))
		return infrastruct.ErrorInternalServerError
	}

	if logFinal.Error != "" {
		if logFinal.Error == "not_enough_funds" {
			logger.LogInfo("err with payment(YOOMONEY): ДЕНЬГИ ЗАКОНЧИЛИСЬ - РАСЧЕХЛЯЙ КРЕДИТКУ ЖЕНЕ4КА")
			logger.LogInfo(fmt.Sprintf("Деньги из за этого не зачислились: %s", userID))
			return infrastruct.ErrorInternalServerError
		}

		logger.LogError(errors.Wrap(errors.Errorf(logFinal.Error), "err with send money in payment(YOOMONEY)"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
