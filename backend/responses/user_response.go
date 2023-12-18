package responses

import (
	"gofr.dev/pkg/gofr"
)

type UserResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}
