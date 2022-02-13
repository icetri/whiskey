package mail

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/whiskey-back/internal/config"
	"github.com/whiskey-back/internal/types"
	"github.com/whiskey-back/pkg/logger"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
	"time"
)

type Mail struct {
	email *config.ForSendEmail
	db    *gorm.DB
}

func NewMail(cfg *config.ForSendEmail, db *gorm.DB) (*Mail, error) {
	return &Mail{
		email: cfg,
		db:    db,
	}, nil
}

func (r *Mail) SendMail(subject, text, to string) error {

	m := gomail.NewMessage()

	m.SetAddressHeader("From", r.email.EmailSender, r.email.NameSender)
	m.SetAddressHeader("To", to, to)

	m.SetHeader("From", fmt.Sprintf("%s <%s>", r.email.NameSender, r.email.EmailSender))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetHeader("MIME-Version:", "1.0")
	m.SetHeader("Reply-To", r.email.EmailUnsubscribe)

	m.SetBody("text/plain", text)

	d := gomail.NewDialer(r.email.EmailHost, r.email.EmailPort, r.email.EmailLogin, r.email.EmailPass)

	stopMail := 0
	for stopMail < 1 {
		stopMail++

		if err := d.DialAndSend(m); err != nil {

			time.Sleep(time.Second * 30)
			if stopMail < 1 {
				continue
			}

			if errInsert := r.db.Debug().Table("email_problems").Create(&types.LogEmail{
				Email: to,
				Data:  err.Error(),
			}).Error; errInsert != nil {
				logger.LogError(errors.Wrap(errInsert, "err with errInsert"))
			}

			return errors.Wrap(err, "err with DialAndSend mail")
		}

		break
	}

	return nil
}
