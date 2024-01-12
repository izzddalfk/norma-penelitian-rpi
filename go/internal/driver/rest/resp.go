package rest

import "net/http"

type Response struct {
	ServiceID string      `json:"svc_id"`
	Code      int         `json:"code"`
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}

func NewSuccessResponse(data interface{}, svcID string) Response {
	return Response{
		ServiceID: svcID,
		Code:      http.StatusOK,
		Status:    "OK",
		Data:      data,
	}
}

func NewInternalServerErrorResponse(errorMessage interface{}) Response {
	return Response{
		Code:   http.StatusInternalServerError,
		Status: "ERR_INTERNAL_SERVER",
		Errors: errorMessage,
	}
}

func NewBadRequestErrorResponse(errorMessage interface{}) Response {
	return Response{
		Code:   http.StatusBadRequest,
		Status: "ERR_BAD_REQUEST",
		Errors: errorMessage,
	}
}
