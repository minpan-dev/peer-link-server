package repository

import (
	"context"
	"peer-link-server/internal/model"
	apperrors "peer-link-server/pkg/errors"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Create(ctx context.Context, room *model.Room) error
	FindByName(ctx context.Context, name string) (*model.Room, error)
	List(ctx context.Context, page, pageSize int) ([]model.Room, int64, error)
	Delete(ctx context.Context, name string) error
}

type roomRepository struct{ db *gorm.DB }

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *model.Room) error {
	if err := r.db.WithContext(ctx).Create(room).Error; err != nil {
		return apperrors.Wrap(409, 40900, "room already exists", err)
	}
	return nil
}

func (r *roomRepository) FindByName(ctx context.Context, name string) (*model.Room, error) {
	var room model.Room
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&room).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.Wrap(500, 50000, "db error", err)
	}
	return &room, nil
}

func (r *roomRepository) List(ctx context.Context, page, pageSize int) ([]model.Room, int64, error) {
	var rooms []model.Room
	var total int64
	offset := (page - 1) * pageSize
	r.db.WithContext(ctx).Model(&model.Room{}).Count(&total)
	err := r.db.WithContext(ctx).Offset(offset).Limit(pageSize).Find(&rooms).Error
	return rooms, total, err
}

func (r *roomRepository) Delete(ctx context.Context, name string) error {
	return r.db.WithContext(ctx).Where("name = ?", name).Delete(&model.Room{}).Error
}
