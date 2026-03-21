package repository

import (
	"LifeNavigator/internal/models"
	"errors"
	"math/rand/v2"

	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

func (r *baseRepository) create(model models.Model) error {
	var ID uint64
	for { //分配随机 ID
		ID = rand.Uint64()
		err := r.db.Model(model).Where(ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return err
		}
	}
	model.SetID(ID)
	return r.db.Create(model).Error
}

func (r *baseRepository) createWithTX(tx *gorm.DB, model models.Model) error {
	var ID uint64
	for { //分配随机 ID
		ID = rand.Uint64()
		err := tx.Model(model).Where(ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return err
		}
	}
	model.SetID(ID)
	return tx.Create(model).Error
}
