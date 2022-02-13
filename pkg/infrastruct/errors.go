package infrastruct

import "net/http"

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

var (
	ErrorInternalServerError = NewError("internal server error", http.StatusInternalServerError)
	ErrorBadRequest          = NewError("bad query input", http.StatusBadRequest)
	ErrorInvalidLink         = NewError("invalid link", http.StatusBadRequest)
	ErrorCheckNotValid       = NewError("check is not valid", http.StatusBadRequest)
	ErrorUserNotExist        = NewError("phone does not exist", http.StatusBadRequest)
	ErrorIncorrectCode       = NewError("incorrect code", http.StatusBadRequest)
	ErrorServiceUnavailable  = NewError("check service unavailable", http.StatusBadRequest)
	ErrorDuplicate           = NewError("duplicate", http.StatusBadRequest)
	ErrorDate                = NewError("outside the date of the promotion", http.StatusBadRequest)
	ErrorEmailIsExist        = NewError("email already registered", http.StatusConflict)
	ErrorPhoneIsExist        = NewError("phone already registered", http.StatusConflict)
	ErrorShop                = NewError("the store does not participate in the promotion", http.StatusBadRequest)
	ErrorCities              = NewError("the city does not participate in the promotion", http.StatusBadRequest)
	ErrorGiftCountZero       = NewError("the number of gifts cannot be zero", http.StatusTeapot)
	ErrorJWTIsBroken         = NewError("jwt spoiled", http.StatusForbidden)
	ErrorWrongIDGift         = NewError("wrong id of the gift", http.StatusBadRequest)
	ErrorInvalidHeader       = NewError("invalid auth header", http.StatusForbidden)
	ErrorTotalLimit          = NewError("the total amount has exceeded the limit", http.StatusBadRequest)
	ErrorInvalidRole         = NewError("role does not match", http.StatusForbidden)
	ErrorPermissionDenied    = NewError("you don't have enough rights", http.StatusForbidden)
	ErrorBalanceEmpty        = NewError("insufficient funds on the balance sheet", http.StatusBadRequest)
)
