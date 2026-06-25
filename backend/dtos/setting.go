package dtos

type UpdateSettingRequest struct {
	Value string `json:"value" binding:"required"`
}
