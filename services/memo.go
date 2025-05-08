package services

import (
	"github.com/kenta-kenta/gaishi-serverside/models"
	"github.com/kenta-kenta/gaishi-serverside/repositories"
)

func CreateMemo(req models.CreateMemoRequest) (*models.Memo, error) {
	return repositories.CreateMemo(req)
}

func GetMemoByID(id int) (*models.Memo, error) {
	return repositories.GetMemoByID(id)
}

func UpdateMemo(id int, req models.UpdateMemoRequest) (*models.Memo, error) {
	return repositories.UpdateMemo(id, req)
}

func DeleteMemo(id int) error {
	return repositories.DeleteMemo(id)
}
