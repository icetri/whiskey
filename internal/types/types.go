package types

import "time"

type Status string

var Accepted Status = "Принят"
var NotAccepted Status = "Не принят"
var OnCheck Status = "На проверке"

type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	Code      string    `json:"code" gorm:"column:code"`
	Balance   float64   `json:"balance" gorm:"column:balance;default:0"`
	Action    float64   `json:"action" gorm:"column:action;default:4000"`
	Role      string    `json:"role" gorm:"column:role;default:USER"`
}

type CheckUser struct {
	ID    string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Phone string `json:"phone" gorm:"column:phone"`
	Code  string `json:"code" gorm:"column:code"`
}

type Token struct {
	Token string `json:"token"`
}

type Profile struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	Role      string    `json:"role" gorm:"column:role;default:USER"`
	Balance   float64   `json:"balance" gorm:"column:balance"`

	Cheques  Cheques  `json:"cheques" gorm:"foreignKey:UserID"`
	Products Products `json:"products" gorm:"foreignKey:UserID"`
}

type Cheques []Cheque
type Cheque struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	Date        string    `json:"date" gorm:"column:date"`
	CheckAmount string    `json:"check_amount" gorm:"column:check_amount"`
	FN          string    `json:"fn" gorm:"column:fn"`
	FD          string    `json:"fd" gorm:"column:fd"`
	FP          string    `json:"fp" gorm:"column:fp"`
	Check       Status    `json:"check" gorm:"column:check"`
	Winning     float64   `json:"winning" gorm:"column:winning"`
	UserID      string    `json:"user_id" gorm:"column:user_id"`

	PhotoID string `json:"-" gorm:"column:photo_id"`
}

type Products []Product
type Product struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	GiftID    int       `json:"gift_id" gorm:"column:gift_id"`
	Count     int       `json:"count" gorm:"column:count"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`

	Gift *Gift `json:"gift" gorm:"-"`
}

type Gifts []Gift
type Gift struct {
	ID    int     `json:"id" gorm:"column:id"`
	Name  string  `json:"name" gorm:"column:name"`
	Price float64 `json:"price" gorm:"column:price"`
	Sum   float64 `json:"-" gorm:"column:sum"`
}

type Sms struct {
	Messages []SMSMessages `json:"messages"`
	Login    string        `json:"login"`
	Password string        `json:"password"`
}

type SMSMessages struct {
	Phone    string `json:"phone"`
	ClientID string `json:"clientId"`
	Text     string `json:"text"`
	Sender   string `json:"sender"`
}

type LogSMS struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Resp      string    `json:"resp" gorm:"column:resp"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
}

type LogEmail struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Data      string    `json:"data" gorm:"column:data"`
	Email     string    `json:"email" gorm:"column:email"`
}

type PreloadAnswer struct {
	Code int `json:"code"`
}

type Check struct {
	Code int  `json:"code"`
	Data Data `json:"data"`
}

type Items struct {
	Sum      float64 `json:"sum"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity float64 `json:"quantity"`
}

type JSON struct {
	Code                    int     `json:"code" gorm:"column:email"`
	User                    string  `json:"user" gorm:"column:email"`
	Items                   []Items `json:"items" gorm:"column:email"`
	Fnsurl                  string  `json:"fnsUrl" gorm:"column:email"`
	Userinn                 string  `json:"userInn" gorm:"column:email"`
	Datetime                string  `json:"dateTime" gorm:"column:email"`
	Kktregid                string  `json:"kktRegId" gorm:"column:email"`
	Totalsum                int     `json:"totalSum" gorm:"column:email"`
	Fiscalsign              int     `json:"fiscalSign" gorm:"column:email"`
	Retailplace             string  `json:"retailPlace" gorm:"column:email"`
	Shiftnumber             int     `json:"shiftNumber" gorm:"column:email"`
	Operationtype           int     `json:"operationType" gorm:"column:email"`
	Requestnumber           int     `json:"requestNumber" gorm:"column:email"`
	Fiscaldrivenumber       string  `json:"fiscalDriveNumber" gorm:"column:email"`
	Retailplaceaddress      string  `json:"retailPlaceAddress" gorm:"column:email"`
	Fiscaldocumentnumber    int     `json:"fiscalDocumentNumber" gorm:"column:email"`
	Fiscaldocumentformatver int     `json:"fiscalDocumentFormatVer" gorm:"column:email"`
}

type Data struct {
	JSON JSON `json:"json"`
}

type LoggerReqCheck struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Logger    string    `json:"logger" gorm:"column:logger"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
}

type LoggerRespCheck struct {
	ID               string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt        time.Time `json:"created_at" gorm:"column:created_at"`
	Logger           string    `json:"logger" gorm:"column:logger"`
	UserID           string    `json:"user_id" gorm:"column:user_id"`
	LoggerReqCheckID string    `json:"logger_req_check_id" gorm:"column:logger_req_check_id"`
}

type GoodRespCheck struct {
	ID                      string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt               time.Time `json:"created_at" gorm:"column:created_at"`
	Code                    int       `json:"code" gorm:"column:code"`
	User                    string    `json:"user" gorm:"column:user"`
	FnsUrl                  string    `json:"fns_url" gorm:"column:fns_url"`
	UserInn                 string    `json:"user_inn" gorm:"column:user_inn"`
	DateTime                string    `json:"data_time" gorm:"column:data_time"`
	KktRegID                string    `json:"kkt_req_id" gorm:"column:kkt_req_id"`
	TotalSum                int       `json:"total_sum" gorm:"column:total_sum"`
	FiscalSign              int       `json:"fiscal_sign" gorm:"column:fiscal_sign"`
	RetailPlace             string    `json:"retail_place" gorm:"column:retail_place"`
	ShiftNumber             int       `json:"shift_number" gorm:"column:shift_number"`
	OperationType           int       `json:"operation_type" gorm:"column:operation_type"`
	RequestNumber           int       `json:"request_number" gorm:"column:request_number"`
	FiscalDriveNumber       string    `json:"fiscal_driver_number" gorm:"column:fiscal_driver_number"`
	RetailPlaceAddress      string    `json:"retail_place_address" gorm:"column:retail_place_address"`
	FiscalDocumentNumber    int       `json:"fiscal_document_number" gorm:"column:fiscal_document_number"`
	FiscalDocumentFormatVer int       `json:"fiscal_document_format_ver" gorm:"column:fiscal_document_format_ver"`
	Url                     string    `json:"url" gorm:"column:url"`
	UserID                  string    `json:"user_id" gorm:"column:user_id"`
}
type PositionInCheck struct {
	ID              string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	Name            string    `json:"name" gorm:"column:name"`
	Price           float64   `json:"price" gorm:"column:price"`
	Count           float64   `json:"count" gorm:"column:count"`
	Sum             float64   `json:"sum" gorm:"column:sum"`
	Marker          string    `json:"marker" gorm:"column:marker"`
	GoodRespCheckID string    `json:"good_resp_check_id" gorm:"column:good_resp_check_id"`
	UserID          string    `json:"user_id" gorm:"column:user_id"`
}

type Whiskey struct {
	ID   int    `json:"id" gorm:"column:id"`
	Name string `json:"product_name" gorm:"column:product_name"`
}

type Shop struct {
	ID   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Inn  string `json:"inn" gorm:"column:inn"`
}

type Prize struct {
	Products []Present `json:"products"`
}

type Present struct {
	ID    int  `json:"id"`
	Count uint `json:"count"`
}
type RequestGifts []RequestGift
type RequestGift struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	Phone       string    `json:"phone" gorm:"column:phone"`
	GiftID      int       `json:"gift_id" gorm:"column:gift_id"`
	PrizeName   string    `json:"prize_name" gorm:"column:prize_name"`
	Certificate string    `json:"certificate" gorm:"column:certificate"`
	ProductID   string    `json:"product_id" gorm:"column:product_id"`
	UserID      string    `json:"user_id" gorm:"column:user_id"`
	Send        bool      `json:"sent" gorm:"column:sent;default:false"`
}

type Support struct {
	FIO   string `json:"fio"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Theme string `json:"theme"`
	Text  string `json:"text"`
}

type ForCSVTable struct {
	ID          string `csv:"id"`
	CreatedAt   string `csv:"created_at"`
	Phone       string `csv:"phone"`
	GiftID      string `csv:"gift_id"`
	PrizeName   string `csv:"prize_name"`
	Certificate string `csv:"certificate"`
	ProductID   string `csv:"product_id"`
	UserID      string `csv:"user_id"`
}

type YoomoneyLogs struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UserID    string    `json:"user_id" gorm:"column:user_id"`
	Type      string    `json:"type" gorm:"column:type"`
	Error     string    `json:"error" gorm:"column:error"`
	RequestID string    `json:"request_id" gorm:"column:request_id"`
	Status    string    `json:"status" gorm:"column:status"`
	Amount    string    `json:"amount" gorm:"column:amount"`
	PaymentID string    `json:"payment_id" gorm:"column:payment_id"`
	InvoiceID string    `json:"invoice_id" gorm:"column:invoice_id"`
	Phone     string    `json:"phone" gorm:"column:phone"`
}

type Parameters struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type PagList struct {
	Items interface{} `json:"items"`
	Count int64       `json:"count"`
}

type Files []File
type File struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	Name      string    `json:"name" gorm:"column:name"`
	Url       string    `json:"url" gorm:"column:url"`
	Length    int64     `json:"length" gorm:"column:length"`
	MimeType  string    `json:"mime" gorm:"column:mime"`
	Bucket    string    `json:"bucket" gorm:"column:bucket"`
	Object    string    `json:"object" gorm:"column:object"`
	UserId    string    `json:"user_id" gorm:"column:user_id"`
	Review    bool      `json:"-" gorm:"column:review;default:false"`
}

type Statistics struct {
	TotalChecks    int64 `json:"total_checks" gorm:"column:total_checks"`
	ChecksAccepted int64 `json:"checks_accepted" gorm:"column:checks_accepted"`
	ChecksRejected int64 `json:"checks_rejected" gorm:"column:checks_rejected"`
	ChecksOnCheck  int64 `json:"checks_on_check" gorm:"column:checks_on_check"`
	CheckOrders    int64 `json:"check_orders" gorm:"column:check_orders"`
}

type Cities struct {
	ID   int    `json:"id" gorm:"column:id"`
	City string `json:"city" gorm:"column:city"`
}

type UsersClients []UserClient
type UserClient struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	Email     string    `json:"email" gorm:"column:email"`
	Phone     string    `json:"phone" gorm:"column:phone"`
	FirstName string    `json:"first_name" gorm:"column:first_name"`
	Code      string    `json:"code" gorm:"column:code"`
	Balance   float64   `json:"balance" gorm:"column:balance;default:0"`
	Action    float64   `json:"action" gorm:"column:action;default:4000"`
	Role      string    `json:"role" gorm:"column:role;default:USER"`
	Client    string    `json:"client" gorm:"column:client"`
}

type CountUsersGift struct {
	CountUsersSite int64 `json:"count_users_site"`
	CountUsersBot  int64 `json:"count_users_bot"`
	CountGift      int64 `json:"count_gift"`
}
