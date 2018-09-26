package repository

import (
	"1024casts/backend/model"
	"fmt"
	"1024casts/backend/pkg/constvar"
)

type UserRepo struct {
	db *model.Database
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: model.DB,
	}
}

func (repo *UserRepo) GetUserById(id int) (*model.UserModel, error) {
	user := model.UserModel{}
	result := repo.db.Self.Where("id = ?", id).First(&user)

	return &user, result.Error
}

func (repo *UserRepo) GetUserList(username string, offset, limit int) ([]*model.UserModel, uint64, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*model.UserModel, 0)
	var count uint64

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := repo.db.Self.Model(&model.UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := repo.db.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

func (repo *UserRepo) FindByEmail(email string) (*model.UserModel, error) {
	users := model.UserModel{}
	return &users, nil
}

func (repo *UserRepo) FindByChangePasswordHash(hash string) (*model.UserModel, error) {
	users := model.UserModel{}
	return &users, nil
}

func (repo *UserRepo) FindAll() ([]*model.UserModel, error) {
	users := make([]*model.UserModel, 0)

	//repo.db.Self.Model()

	return users, nil
}

func (repo *UserRepo) Update(user *model.UserModel) error {
	return nil
}

//func (repo *UserRepo) Store(user *model.UserModel) (entity.ID, error) {
//	users := model.UserModel{}
//	return &users, nil
//}

func (repo *UserRepo) FindByValidationHash(hash string) (*model.UserModel, error) {
	users := model.UserModel{}
	return &users, nil
}
