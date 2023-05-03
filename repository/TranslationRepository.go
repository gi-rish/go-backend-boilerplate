package repository

import (
	"context"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/models"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TranslationRepository interface {
	CreateTranslations(req http.Request, code, key, translations string) error
	GetTranslationsByKey(ctx context.Context, key []string, code string) ([]*models.Translation, error)
	GetAllTranslations(req http.Request) ([]*models.Translation, error)
	GetTranslationsByLanguage(req http.Request, code string) ([]*models.Translation, error)
	GetDefaultLanguageTranslations(req http.Request) ([]*models.Translation, error)
	UpdateTranslations(req http.Request, code, key, translations string) error
}

type TranslationRepo struct {
	DB *gorm.DB
}

func (r *TranslationRepo) CreateTranslations(req http.Request, code, key, translations string) error {
	now := time.Now().UTC()
	translation := &models.Translation{
		Id:           uuid.New(),
		Code:         code,
		Key:          key,
		Translations: translations,
		DateCreated:  &now,
		DateUpdated:  &now,
	}
	if result := r.DB.Create(&translation); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TranslationRepo) GetTranslationsByKey(ctx context.Context, key []string, code string) ([]*models.Translation, error) {
	translations := make([]*models.Translation, 0)
	err := r.DB.Find(&translations, "key IN (?) and code = ?", key, code).Error
	return translations, err
}

func (r *TranslationRepo) GetAllTranslations(req http.Request) ([]*models.Translation, error) {
	appname := strings.ToLower(req.Header.Get("appname"))
	translations := make([]*models.Translation, 0)
	err := r.DB.Find(&translations, "key like ? ", appname+"%").Error
	return translations, err
}

func (r *TranslationRepo) GetTranslationsByLanguage(req http.Request, code string) ([]*models.Translation, error) {
	appname := strings.ToLower(req.Header.Get("appname"))
	translations := make([]*models.Translation, 0)
	err := r.DB.Find(&translations, "code = ? AND key like ? ", code, appname+"%").Error
	return translations, err
}

func (r *TranslationRepo) GetDefaultLanguageTranslations(req http.Request) ([]*models.Translation, error) {
	appname := strings.ToLower(req.Header.Get("appname"))
	translations := make([]*models.Translation, 0)
	err := r.DB.Model(&models.Translation{}).Select("translations.code, translations.key, translations.translations").Joins("JOIN languages as l ON l.code = translations.code").Where("l.default = ? AND translations.key like ?", true, appname+"%").Find(&translations).Error
	return translations, err
}

func (r *TranslationRepo) UpdateTranslations(req http.Request, code, key, translations string) error {
	var translation *models.Translation
	err := r.DB.Model(&translation).Update("translations", translation).Where(models.Translation{Code: code, Key: key}).Error
	return err
}
