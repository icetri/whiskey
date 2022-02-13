package config

type Config struct {
	HTTP        *HTTP         `yaml:"http"`
	PostgresDsn string        `yaml:"postgres_dsn"`
	JWTKey      string        `yaml:"jwtkey"`
	PathUrl     string        `yaml:"path_url"`
	DateStart   string        `yaml:"date_start"`
	Telegram    *Telegram     `yaml:"telegram"`
	SMS         *Sms          `yaml:"sms"`
	Email       *ForSendEmail `yaml:"server_email"`
	FileStorage *Storage      `yaml:"file_storage"`
	Check       *Check        `yaml:"check"`
	Yoomoney    *Yoomoney     `yaml:"yoomoney"`
}

type HTTP struct {
	Port string `yaml:"server_port"`
}

type Telegram struct {
	TelegramToken string `yaml:"telegram_token"`
	ChatID        string `yaml:"chat_id"`
}

type Sms struct {
	SMSURL    string `yaml:"sms_url"`
	SmsPass   string `yaml:"sms_pass"`
	SmsLog    string `yaml:"sms_login"`
	SmsSender string `yaml:"sms_sender"`
}

type ForSendEmail struct {
	EmailHost        string `yaml:"host"`
	EmailPort        int    `yaml:"port"`
	EmailLogin       string `yaml:"login"`
	EmailFrom        string `yaml:"from"`
	EmailPass        string `yaml:"pass"`
	EmailUnsubscribe string `yaml:"email_unsubscribe"`
	NameSender       string `yaml:"email_name_sender"`
	EmailSender      string `yaml:"email_sender"`
}

type Storage struct {
	Host      string `yaml:"host"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	SSL       bool   `yaml:"ssl"`
	Bucket    string `yaml:"bucket"`
	BucketBad string `yaml:"bucket_bad"`
}

type Check struct {
	Url   string `yaml:"url"`
	Token string `yaml:"token"`
}

type Yoomoney struct {
	Auth      string `yaml:"auth"`
	ReqStart  string `yaml:"req_start"`
	ReqFinish string `yaml:"req_finish"`
}
