package handler

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
