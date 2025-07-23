package repository

import "github.com/Sush1sui/sushi-vc-bot-go/internal/models"

type CategoryJTCInterface interface {
	GetAllJTCs() ([]*models.CategoryJTCModel, error)
	CreateCategoryJTC(interfaceId, interfaceMessageId, jtcChannelId, categoryId string) (*models.CategoryJTCModel, error)
	DeleteAll() (int, error)
}

var CategoryJTCService CategoryJTCInterface