package dtos

type ContactResponse struct {
	Phone       string `json:"phone"`
	WhatsAppURL string `json:"whatsapp_url"`
	Email       string `json:"email"`
}
