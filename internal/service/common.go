package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/repository"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/infrastruct"
	"github.com/whiskey-back/pkg/logger"
	"gorm.io/gorm"
	"net/http"
	url2 "net/url"
	"strings"
	"time"
)

type CommonService struct {
	repo repository.Common
}

func NewCommonService(repo repository.Common) *CommonService {
	return &CommonService{
		repo: repo,
	}
}

func (s *CommonService) CheckUsers(userID string) error {

	user, err := s.repo.GetUserByID(userID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if user == nil {
		return infrastruct.ErrorPermissionDenied
	}

	return nil
}

func (s *CommonService) GetGifts() (types.Gifts, error) {

	gifts, err := s.repo.GetGifts()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetGifts in GetGifts"))
		return nil, infrastruct.ErrorBadRequest
	}

	return gifts, nil
}

func (s *CommonService) SupportSend(sup *types.Support) error {

	preText := fmt.Sprintf("От: %s\n Телефон: %s\n Email: %s\n Тема обращения: %s\n Сообщение: %s",
		sup.FIO,
		sup.Phone,
		sup.Email,
		sup.Theme,
		sup.Text,
	)

	preText = url2.QueryEscape(preText)

	text := fmt.Sprintf("%s [%s]: %s", "INFO", time.Now().Format("2006-01-02T15:04:05"), preText)

	url := fmt.Sprintf("https://api.telegram.org/botSECRET/sendMessage?chat_id=SECRET&text=%s", text)
	method := "GET"

	newUrl := strings.ReplaceAll(url, " ", "+")

	client := &http.Client{}
	req, err := http.NewRequest(method, newUrl, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
