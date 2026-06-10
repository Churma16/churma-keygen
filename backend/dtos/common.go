package dtos

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type MetaResponse struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type APIResponse struct {
	Meta MetaResponse `json:"meta"`
	Data interface{}  `json:"data"`
}

func NewSuccessResponse(code int, message string, data interface{}) APIResponse {
	return APIResponse{
		Meta: MetaResponse{
			Code:    code,
			Status:  "success",
			Message: message,
		},
		Data: data,
	}
}

func NewErrorResponse(code int, message string) APIResponse {
	return APIResponse{
		Meta: MetaResponse{
			Code:    code,
			Status:  "error",
			Message: message,
		},
		Data: nil,
	}
}
