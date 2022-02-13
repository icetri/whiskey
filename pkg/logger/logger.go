package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/whiskey-back/internal/config"
	"net/http"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

var (
	chatID    = ""
	infoMsg   = "INFO"
	errMsg    = "ERROR"
	token     = ""
	logger    = zerolog.New(os.Stdout)
	loggerErr = zerolog.New(os.Stderr)
)

func NewLogger(t *config.Telegram) error {

	chatID = t.ChatID
	token = t.TelegramToken

	return nil
}

func LogError(err error) {
	loggerErr.Err(err).Time("time", time.Now()).Send()
	SendError(err)
}

func LogInfo(msg string) {
	logger.Info().Time("time", time.Now()).Msg(msg)
	SendMessage(msg)
}

func LogFatal(err error) {
	SendError(err)
	loggerErr.Fatal().Err(err).Time("time", time.Now()).Send()
}

func SendError(err error) {
	url := makeURLSendMessage(errMsg, url2.QueryEscape(err.Error()))
	if err := send(url); err != nil {
		loggerErr.Err(err).Send()
	}
}

func SendMessage(msg string) {
	url := makeURLSendMessage(infoMsg, url2.QueryEscape(msg))
	if err := send(url); err != nil {
		loggerErr.Err(err).Send()
	}
}

func makeURLSendMessage(typeMsg, text string) string {
	text = fmt.Sprintf("%s [%s]: %s", typeMsg, time.Now().Format("2006-01-02T15:04:05"), text)
	str := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		token, chatID, text)
	return strings.ReplaceAll(str, " ", "+")
}

func send(urlForSend string) error {
	req, err := http.NewRequest(http.MethodPost, urlForSend, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err = res.Body.Close(); err != nil {
			loggerErr.Err(err)
		}
	}()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d", res.StatusCode)
	}
	return nil
}
