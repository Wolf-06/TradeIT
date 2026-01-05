package engine

import "fmt"

type ErrorCodes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorCodes) Error() string {
	return fmt.Sprintf("code: %s\nmessage: %s", e.Code, e.Message)
}

type EngineErrorConfig struct {
	OrderInvalidSide           *ErrorCodes
	OrderInvalidType           *ErrorCodes
	OrderNotInMatcher          *ErrorCodes
	OrderPartialCancelled      *ErrorCodes
	OrderUnsupportedType       *ErrorCodes
	OrderNotFoundOrProcessed   *ErrorCodes
	OrderInvalidQuantity       *ErrorCodes
	OrderPartialDecreaseDenied *ErrorCodes
	OrderAlreadyProcessed      *ErrorCodes
	OrderInvalidPrice          *ErrorCodes
}
