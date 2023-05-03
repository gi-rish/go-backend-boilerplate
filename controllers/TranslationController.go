package controllers

import (
	"context"
	"errors"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/models"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/repository"
	"net/http"
	"strings"

	localizationPB "github.com/NewStreetTechnologies/go-grpc-proto/localization-service"

	_json "github.com/Jeffail/gabs/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TranslationController struct {
	TranslationRepository repository.TranslationRepository
	Logger                *logrus.Logger
	localizationPB.UnimplementedLocalizationServiceServer
}

// InitTranslationController Constructor to initialize controller for grpc call
func InitTranslationController(TranslationRepository repository.TranslationRepository, Logger *logrus.Logger) *TranslationController {
	return &TranslationController{
		TranslationRepository: TranslationRepository,
		Logger:                Logger,
	}
}

// TranslationReq struct for create translation request
type TranslationReq struct {
	Key          string `json:"key"`
	Code         string `json:"code"`
	Translations string `json:"translations"`
}

// CreateTranslation Store translations into Database
func (tc *TranslationController) CreateTranslation(c *gin.Context) {
	errCode, errHeader := ValidateHeader(*c.Request)
	if errHeader != nil {
		tc.Logger.Errorf("Error: %+v", errHeader)
		ErrorResult(c, errHeader, "", errCode, http.StatusOK)
		return
	}

	var translationReq TranslationReq
	if err := c.ShouldBindJSON(&translationReq); err != nil {
		tc.Logger.Errorf("Error parsing request %+v", err)
		ErrorResult(c, errors.New(ErrorCodes["PLS0004"]), ErrorCodes["PLS0004"], "PLS0004", http.StatusInternalServerError)
		return
	}

	err := tc.TranslationRepository.CreateTranslations(*c.Request, translationReq.Code, translationReq.Key, translationReq.Translations)
	if err != nil {
		tc.Logger.Errorf("Error creating translations %+v", err)
		ErrorResult(c, err, "", "", http.StatusInternalServerError)
		return
	}

	SuccessResult(c, gin.H{
		"message": "success",
		"code":    "",
		"data":    gin.H{},
	})
}

// GetTranslationsByKey Get Translation by Key for requested language code
func (tc *TranslationController) GetTranslationsByKey(ctx context.Context, req *localizationPB.GetTranslationsRequest) (*localizationPB.GetTranslationsResponse, error) {
	//check if request parameters are empty
	if len(req.Key) <= 0 || len(req.Code) <= 0 {
		tc.Logger.Errorf("Error: %+v", ErrorCodes["PLS0005"])
		return &localizationPB.GetTranslationsResponse{
			Code:    "PLS0005",
			Message: ErrorCodes["PLS0005"],
			Result:  nil,
		}, status.Errorf(codes.InvalidArgument, ErrorCodes["PLS0005"])
	}
	//check if any of the key value is empty
	for _, k := range req.Key {
		if len(strings.TrimSpace(k)) <= 0 {
			tc.Logger.Errorf("Error: %+v", ErrorCodes["PLS0005"]+" : key should not be empty")
			return &localizationPB.GetTranslationsResponse{
				Code:    "PLS0005",
				Message: ErrorCodes["PLS0005"],
				Result:  nil,
			}, status.Errorf(codes.InvalidArgument, ErrorCodes["PLS0005"]+" : key should not be empty")
		}
	}
	var translations []*models.Translation
	var err error
	//get translations
	if translations, err = tc.TranslationRepository.GetTranslationsByKey(ctx, req.Key, req.Code); err != nil {
		tc.Logger.Errorf("Error: %+v", err)
		return &localizationPB.GetTranslationsResponse{
			Code:    "PLS0006",
			Message: err.Error(),
			Result:  nil,
		}, status.Errorf(codes.Internal, err.Error())
	}
	//check if no translations exist for requested keys or length of translations is not same as supplied no of keys
	if len(translations) == 0 || len(translations) != len(req.Key) {
		tc.Logger.Errorf("Error: %+v", ErrorCodes["PLS0001"]+" for the key/s : "+strings.Join(req.Key, ", "))
		return &localizationPB.GetTranslationsResponse{
			Code:    "PLS0001",
			Message: ErrorCodes["PLS0001"],
			Result:  nil,
		}, status.Errorf(codes.NotFound, ErrorCodes["PLS0001"]+" for the key/s : "+strings.Join(req.Key, ", "))
	}

	res := make(map[string]string)
	for _, t := range translations {
		translation := t.Translations
		//check if header and footer is available in translation. if yes then need to replace header and footer with translation value
		if strings.Contains(translation, "$header$") || strings.Contains(translation, "$footer$") {
			//get app-name to fetch header and footer details
			appName := strings.Split(t.Key, ".")[0]
			var keys []string
			keys = append(keys, appName+".email.header")
			keys = append(keys, appName+".email.footer")
			var trans []*models.Translation
			//fetch header and footer details
			if trans, err = tc.TranslationRepository.GetTranslationsByKey(ctx, keys, "en-US"); err != nil {
				tc.Logger.Errorf("Error while fetching header and footer: %+v", err)
				return &localizationPB.GetTranslationsResponse{
					Code:    "PLS0006",
					Message: err.Error(),
					Result:  nil,
				}, status.Errorf(codes.Internal, err.Error())
			}
			//check if any of the translations for header or footer is not available
			if len(trans) != 2 {
				tc.Logger.Errorf("Translation/s for header or footer not available: %+v", strings.Join(keys, ", "))
				return &localizationPB.GetTranslationsResponse{
					Code:    "PLS0001",
					Message: ErrorCodes["PLS0001"],
					Result:  nil,
				}, status.Errorf(codes.NotFound, "Translation/s for header or footer not available : "+strings.Join(keys, ", "))
			}
			for _, tran := range trans {
				//check for header details
				if tran.Key == keys[0] {
					translation = strings.Replace(translation, "$header$", tran.Translations, 1)
				}
				//check for footer details
				if tran.Key == keys[1] {
					translation = strings.Replace(translation, "$footer$", tran.Translations, 1)
				}
			}
		}
		res[t.Key] = translation
	}
	tc.Logger.Debugf("Response: %+v", res)
	return &localizationPB.GetTranslationsResponse{
		Result: res,
	}, nil
}

// GetAllTranslations get all translations stored in DB
func (tc *TranslationController) GetAllTranslations(c *gin.Context) {
	errCode, errHeader := ValidateHeader(*c.Request)
	if errHeader != nil {
		tc.Logger.Errorf("Error: %+v", errHeader)
		ErrorResult(c, errHeader, "", errCode, http.StatusOK)
		return
	}

	translations, err := tc.TranslationRepository.GetAllTranslations(*c.Request)
	if err != nil {
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, "", "", http.StatusInternalServerError)
		return
	}
	if len(translations) == 0 {
		err := errors.New(ErrorCodes["PLS0001"])
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, ErrorCodes["PLS0001"], "PLS0001", http.StatusOK)
		return
	}
	translation := _json.New()
	for _, t := range translations {
		code := t.Code
		key := strings.Join([]string{code, t.Key}, ".")
		_, err := translation.SetP(t.Translations, key)
		if err != nil {
			return
		}
	}
	res := gin.H{
		"result": translation,
	}
	tc.Logger.Debugf("Response: %+v", res)
	SuccessResult(c, res)
}

// GetTranslationsByLanguage get translation scripts on the basis of requested language code
func (tc *TranslationController) GetTranslationsByLanguage(c *gin.Context) {
	errCode, errHeader := ValidateHeader(*c.Request)
	if errHeader != nil {
		tc.Logger.Errorf("Error: %+v", errHeader)
		ErrorResult(c, errHeader, "", errCode, http.StatusOK)
		return
	}
	translations, err := tc.TranslationRepository.GetTranslationsByLanguage(*c.Request, c.Param("id"))
	if err != nil {
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, "", "", http.StatusInternalServerError)
		return
	}
	if len(translations) == 0 {
		err := errors.New(ErrorCodes["PLS0001"])
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, ErrorCodes["PLS0001"], "PLS0001", http.StatusOK)
		return
	}
	translation := _json.New()
	for _, t := range translations {
		code := t.Code
		key := strings.Join([]string{code, t.Key}, ".")
		_, err := translation.SetP(t.Translations, key)
		if err != nil {
			return
		}
	}
	res := gin.H{
		"result": translation,
	}
	tc.Logger.Debugf("Response: %+v", res)
	SuccessResult(c, res)
}

// GetDefaultLanguageTranslations get default translation scripts (Default -- true [in languages table]])
func (tc *TranslationController) GetDefaultLanguageTranslations(c *gin.Context) {
	errCode, errHeader := ValidateHeader(*c.Request)
	if errHeader != nil {
		tc.Logger.Errorf("Error: %+v", errHeader)
		ErrorResult(c, errHeader, "", errCode, http.StatusOK)
		return
	}
	translations, err := tc.TranslationRepository.GetDefaultLanguageTranslations(*c.Request)
	if len(translations) == 0 {
		err := errors.New(ErrorCodes["PLS0001"])
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, ErrorCodes["PLS0001"], "PLS0001", http.StatusOK)
		return
	}
	if err != nil {
		tc.Logger.Errorf("Error: %+v", err)
		ErrorResult(c, err, "", "", http.StatusInternalServerError)
		return
	}
	translation := _json.New()
	for _, t := range translations {
		code := t.Code
		key := strings.Join([]string{code, t.Key}, ".")
		_, err := translation.SetP(t.Translations, key)
		if err != nil {
			return
		}
	}
	res := gin.H{
		"result": translation,
	}
	tc.Logger.Debugf("Response: %+v", res)
	SuccessResult(c, res)
}
