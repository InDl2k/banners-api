package cache

import (
	"banners/internal/database/models"
)

type BannerCache interface {
	Set(key string, value *models.Banner)
	Get(key string) *models.Banner
}
