package domain

import (
	"github.com/google/uuid"
)

type ClientRepository interface {
	FindAll() ([]Client, error)
	FindByID(id uuid.UUID) (*Client, error)
	Create(client *Client) error
	Update(client *Client) error
	Delete(id uuid.UUID) error
	Count() (int64, error)
}
