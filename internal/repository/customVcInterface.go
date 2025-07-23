package repository

import "github.com/Sush1sui/sushi-vc-bot-go/internal/models"

type CustomVcInterface interface {
	GetAllVcs() ([]*models.CustomVcModel, error)
	CreateVc(ownerId string, channelId string) (*models.CustomVcModel, error)
	GetByOwnerOrChannelId(ownerId string, channelId string) (*models.CustomVcModel, error)
	DeleteByOwnerOrChannelId(ownerId string, channelId string) (int, error)
	ChangeOwnerByChannelId(channelId string, newOwnerId string) (int, error)
}

var CustomVcService CustomVcInterface