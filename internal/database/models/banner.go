package models

import (
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"time"
)

type Banner struct {
	ID        int                `json:"banner_id" gorm:"primary_key;AUTO_INCREMENT"`
	Content   *datatypes.JSONMap `json:"content" gorm:"type:jsonb"`
	FeatureID *int               `json:"feature_id"`
	Tags      *pq.Int64Array     `json:"tag_ids" gorm:"type:int[]"`
	Active    bool               `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time          `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
}

func (banner *Banner) Create() (*map[string]interface{}, error) {

	err := GetDB().Create(&banner).Error

	if err != nil {
		return nil, err
	}
	resp := map[string]interface{}{"banner_id": banner.ID}

	return &resp, nil
}

func GetBanners(feature_id int, tag_id int, onlyActive bool, limit int, offset int) (*[]Banner, error) {
	banners := make([]Banner, 0)

	query := GetDB().Table("banners").Limit(limit).Offset(offset)
	if feature_id != -1 {
		query = query.Where("feature_id = ?", feature_id)
	}
	if tag_id != -1 {
		query = query.Where("? = ANY(tags)", tag_id)
	}
	if onlyActive {
		query = query.Where("active = true")
	}

	err := query.Find(&banners).Error
	if err != nil {
		return nil, err
	}

	return &banners, nil
}

func (banner *Banner) Save() error {

	err := GetDB().Save(&banner).Error
	if err != nil {
		return err
	}

	return nil
}

func (banner *Banner) Delete() error {

	err := GetDB().Delete(&banner).Error
	if err != nil {
		return err
	}

	return nil
}

func Exist(banner_id int) (bool, error) {

	var result struct {
		Found bool
	}

	row := db.Raw("SELECT EXISTS(SELECT 1 FROM banners WHERE id = ?) AS found", banner_id).Scan(&result)
	if row.Error != nil {
		return false, row.Error
	}

	return result.Found, nil
}
