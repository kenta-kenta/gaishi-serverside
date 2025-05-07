package services

import (
	"github.com/kenta-kenta/gaishi-serverside/models"
	"github.com/kenta-kenta/gaishi-serverside/repositories"
)

func CreateMemo(repo *repositories.Repository, req models.CreateMemoRequest) (*models.Memo, error) {
	return repo.CreateMemo(req)
}

func GetMemoByID(repo *repositories.Repository, id int) (*models.Memo, error) {
	return repo.GetMemoByID(id)
}

func UpdateMemo(repo *repositories.Repository, id int, req models.UpdateMemoRequest) (*models.Memo, error) {
	return repo.UpdateMemo(id, req)
}

func DeleteMemo(repo *repositories.Repository, id int) error {
	return repo.DeleteMemo(id)
} 