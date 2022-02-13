package sms

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/logger"
	"io/ioutil"
	"net/http"
)

type ApiSMS struct {
	Sender string
	sms    *config.Sms
}

func NewSMS(sms *config.Sms) (*ApiSMS, error) {
	return &ApiSMS{
		sms:    sms,
		Sender: sms.SmsSender,
	}, nil
}

func (s *ApiSMS) SendSMS(mes []types.SMSMessages) ([]byte, error) {
	smsBody := types.Sms{Login: s.sms.SmsLog,
		Password: s.sms.SmsPass,
		Messages: mes,
	}

	jsonStr, err := json.Marshal(smsBody)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Marshal in SendSMS"))
		return nil, err
	}

	req, err := http.NewRequest("POST", s.sms.SMSURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.LogError(errors.Wrap(err, "Err with http.NewRequest in SendSMS"))
		return nil, infrastruct.ErrorInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "Err with http.Client in SendSMS"))
		return nil, infrastruct.ErrorInternalServerError
	}

	respLogerBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in ChoiceCertificate for logger"))
	}
	defer resp.Body.Close()

	return respLogerBytes, nil
}
