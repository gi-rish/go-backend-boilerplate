package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NewStreetTechnologies/go-backend-boilerplate/controllers"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/models"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/routes"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/test/mocks"

	localizationPb "github.com/NewStreetTechnologies/go-grpc-proto/localization-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTranslationController_CreateTranslation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success", func(t *testing.T) {
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("CreateTranslations", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		m, b := map[string]interface{}{
			"key":          "agri.email.success",
			"code":         "eu-UK",
			"translations": "Hello World Template",
		}, new(bytes.Buffer)
		json.NewEncoder(b).Encode(m)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/translations/create", b)
		request.Header.Set("appname", "mvpmp")
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("Fail Case : calling without appname header", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:      r,
			Logger: logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/translations/create", nil)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"please provide appname header value\",\"error_code\":\"PLS0002\",\"result\":\"fail\"}")
	})

	t.Run("Fail Case : Empty body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:      r,
			Logger: logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/translations/create", nil)
		request.Header.Set("appname", "mvpmp")
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"wrong request format\",\"error_code\":\"PLS0004\",\"result\":\"fail\"}")
	})

	t.Run("Fail Case : calling with invalid appname header value(Not Alphabets)", func(t *testing.T) {
		appname := "mvpmp."
		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:      r,
			Logger: logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodPost, "/api/v1/translations/create", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"appname header value is not valid\",\"error_code\":\"PLS0003\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : unable to create translations", func(t *testing.T) {
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("CreateTranslations", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("unable to create translations"))
		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		m, b := map[string]interface{}{
			"key":          "agri.email.success",
			"code":         "eu-UK",
			"translations": "Hello World Template",
		}, new(bytes.Buffer)
		json.NewEncoder(b).Encode(m)
		request, err := http.NewRequest(http.MethodPost, "/api/v1/translations/create", b)
		request.Header.Set("appname", "agri")
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"unable to create translations\",\"error_code\":\"\",\"result\":\"fail\"}")
	})
}

func TestTranslationController_GetAllTranslations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success : calling with proper parameters", func(t *testing.T) {
		appname := "mvpmp"
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetAllTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/all", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, fmt.Sprintf("%+v", w.Body), appname)
	})
	t.Run("Fail Case : calling without appname header", func(t *testing.T) {
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetAllTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/all", nil)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"please provide appname header value\",\"error_code\":\"PLS0002\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : calling with invalid appname header value(Not Alphabets)", func(t *testing.T) {
		appname := "mvpmp."
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetAllTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/all", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"appname header value is not valid\",\"error_code\":\"PLS0003\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : No data found when calling with wrong appname", func(t *testing.T) {
		appname := "appname"
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := make([]*models.Translation, 0)

		mockTranslationRepo.On("GetAllTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/all", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"no translations found\",\"error_code\":\"PLS0001\",\"result\":\"fail\"}")
	})
}
func TestTranslationController_GetDefaultLanguageTranslations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success : calling with proper parameters", func(t *testing.T) {
		appname := "mvpmp"
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetDefaultLanguageTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, fmt.Sprintf("%+v", w.Body), appname)
	})
	t.Run("Fail Case : calling without appname header", func(t *testing.T) {
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetDefaultLanguageTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/", nil)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"please provide appname header value\",\"error_code\":\"PLS0002\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : calling with invalid appname header value(Not Alphabets)", func(t *testing.T) {
		appname := "mvpmp."
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := GetTranslations()

		mockTranslationRepo.On("GetDefaultLanguageTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"appname header value is not valid\",\"error_code\":\"PLS0003\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : no data found when calling with wrong appname", func(t *testing.T) {
		appname := "appname"
		mockTranslationRepo := new(mocks.TranslationRepository)
		translations := make([]*models.Translation, 0)

		mockTranslationRepo.On("GetDefaultLanguageTranslations", mock.Anything).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}

		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, "/api/v1/translations/", nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"no translations found\",\"error_code\":\"PLS0001\",\"result\":\"fail\"}")
	})
}

func TestTranslationController_GetTranslationsByLanguage(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success : calling with proper parameters", func(t *testing.T) {
		translations := GetTranslations()
		langCode := "en-US"
		appname := "mvpmp"
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByLanguage", mock.Anything, langCode).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/translations/%s", langCode), nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, fmt.Sprintf("%+v", w.Body), fmt.Sprintf("{\"result\":{\"%s\":", langCode))
	})
	t.Run("Fail Case : calling with wrong language code", func(t *testing.T) {
		translations := make([]*models.Translation, 0)
		langCode := "en-USP"
		appname := "mvpmp"
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByLanguage", mock.Anything, langCode).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/translations/%s", langCode), nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"no translations found\",\"error_code\":\"PLS0001\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : no data found when calling with wrong appname", func(t *testing.T) {
		translations := make([]*models.Translation, 0)
		langCode := "en-US"
		appname := "appname"
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByLanguage", mock.Anything, langCode).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/translations/%s", langCode), nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"no translations found\",\"error_code\":\"PLS0001\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : calling without appname header", func(t *testing.T) {
		translations := GetTranslations()
		langCode := "en-US"
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByLanguage", mock.Anything, langCode).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/translations/%s", langCode), nil)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"please provide appname header value\",\"error_code\":\"PLS0002\",\"result\":\"fail\"}")
	})
	t.Run("Fail Case : calling with invalid appname header value(Not Alphabets)", func(t *testing.T) {
		translations := GetTranslations()
		langCode := "en-US"
		appname := "mvpmp."
		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByLanguage", mock.Anything, langCode).Return(translations, nil)

		w := httptest.NewRecorder()
		r := gin.Default()
		routerConfig := &routes.RouterConfig{
			R:                     r,
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		}
		routerConfig.InitRouter()
		request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/translations/%s", langCode), nil)
		request.Header.Set("appname", appname)
		assert.NoError(t, err)
		r.ServeHTTP(w, request)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, fmt.Sprintf("%+v", w.Body), "{\"error\":\"appname header value is not valid\",\"error_code\":\"PLS0003\",\"result\":\"fail\"}")
	})

}

func TestTranslationController_GetTranslationsByKey(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		// Creating Server
		l, err := net.Listen("tcp", "localhost:8000")
		assert.NoError(t, err)
		s := grpc.NewServer()

		// Creating Client
		clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(l.Addr().String(), clientOptions...)
		assert.NoError(t, err)

		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByKey", mock.Anything, mock.Anything, mock.Anything).Return(GetMockTranslations(), nil)
		localizationPb.RegisterLocalizationServiceServer(s, &controllers.TranslationController{
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		})
		go func() {
			if err := s.Serve(l); err != nil {
				panic(fmt.Sprintf("Failed to start gRPC server: %+v !!!!!!", err))
			}
		}()
		defer s.Stop()
		client := localizationPb.NewLocalizationServiceClient(conn)
		response, err := client.GetTranslationsByKey(ctx, &localizationPb.GetTranslationsRequest{
			Key: []string{
				"agri.login.success.email",
				"agri.email.header",
				"agri.email.footer",
			},
			Code: "en-US",
		})
		assert.NoError(t, err)
		assert.Equal(t, response.GetResult(), GetMockTranslationResp())
	})

	t.Run("NoKey", func(t *testing.T) {
		ctx := context.Background()
		// Creating Server
		l, err := net.Listen("tcp", "localhost:8000")
		assert.NoError(t, err)
		s := grpc.NewServer()

		// Creating Client
		clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(l.Addr().String(), clientOptions...)
		assert.NoError(t, err)

		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByKey", mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.InvalidArgument, "Key & Code are required"))
		localizationPb.RegisterLocalizationServiceServer(s, &controllers.TranslationController{
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		})
		go func() {
			if err := s.Serve(l); err != nil {
				panic(fmt.Sprintf("Failed to start gRPC server: %+v !!!!!!", err))
			}
		}()
		defer s.Stop()
		client := localizationPb.NewLocalizationServiceClient(conn)
		response, err := client.GetTranslationsByKey(ctx, &localizationPb.GetTranslationsRequest{
			Key: []string{
				"agri.login.success.email",
				"agri.email.header",
				"agri.email.footer",
			},
			Code: "en-US",
		})
		assert.Error(t, err)
		assert.Equal(t, response.GetResult(), map[string]string(nil))
		if er, ok := status.FromError(err); ok {
			assert.Equal(t, er.Code(), codes.Internal)
			assert.Equal(t, er.Message(), "rpc error: code = InvalidArgument desc = Key & Code are required")
		}
	})

	t.Run("NoCode", func(t *testing.T) {
		ctx := context.Background()
		// Creating Server
		l, err := net.Listen("tcp", "localhost:8000")
		assert.NoError(t, err)
		s := grpc.NewServer()

		// Creating Client
		clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(l.Addr().String(), clientOptions...)
		assert.NoError(t, err)

		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByKey", mock.Anything, mock.Anything, mock.Anything).Return(nil, status.Errorf(codes.InvalidArgument, "Key & Code are required"))
		localizationPb.RegisterLocalizationServiceServer(s, &controllers.TranslationController{
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		})
		go func() {
			if err := s.Serve(l); err != nil {
				panic(fmt.Sprintf("Failed to start gRPC server: %+v !!!!!!", err))
			}
		}()
		defer s.Stop()
		client := localizationPb.NewLocalizationServiceClient(conn)
		response, err := client.GetTranslationsByKey(ctx, &localizationPb.GetTranslationsRequest{
			Key: []string{
				"agri.login.success.email",
				"agri.email.header",
				"agri.email.footer",
			},
			Code: "en-US",
		})
		assert.Error(t, err)
		assert.Equal(t, response.GetResult(), map[string]string(nil))
		if er, ok := status.FromError(err); ok {
			assert.Equal(t, er.Code(), codes.Internal)
			assert.Equal(t, er.Message(), "rpc error: code = InvalidArgument desc = Key & Code are required")
		}
	})

	t.Run("No Translation/s found", func(t *testing.T) {
		ctx := context.Background()
		// Creating Server
		l, err := net.Listen("tcp", "localhost:8000")
		assert.NoError(t, err)
		s := grpc.NewServer()

		// Creating Client
		clientOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		conn, err := grpc.Dial(l.Addr().String(), clientOptions...)
		assert.NoError(t, err)

		mockTranslationRepo := new(mocks.TranslationRepository)
		mockTranslationRepo.On("GetTranslationsByKey", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
		localizationPb.RegisterLocalizationServiceServer(s, &controllers.TranslationController{
			TranslationRepository: mockTranslationRepo,
			Logger:                logrus.New(),
		})
		go func() {
			if err := s.Serve(l); err != nil {
				panic(fmt.Sprintf("Failed to start gRPC server: %+v !!!!!!", err))
			}
		}()
		defer s.Stop()
		client := localizationPb.NewLocalizationServiceClient(conn)
		response, err := client.GetTranslationsByKey(ctx, &localizationPb.GetTranslationsRequest{
			Key: []string{
				"agri.login.success.email",
				"agri.email.header",
				"agri.email.footer",
			},
			Code: "en-US",
		})
		assert.Error(t, err)
		assert.Equal(t, response.GetResult(), map[string]string(nil))
		if er, ok := status.FromError(err); ok {
			assert.Equal(t, er.Code(), codes.NotFound)
			assert.Equal(t, er.Message(), "no translations found for the key/s : agri.login.success.email, agri.email.header, agri.email.footer")
		}
	})
}

func GetTranslations() []*models.Translation {
	arr := make([]*models.Translation, 0)
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "en-US",
		Key:          "mvpmp.language.chinese",
		Translations: "Chinese",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "en-US",
		Key:          "mvpmp.language.english",
		Translations: "English",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "zh-CHS",
		Key:          "mvpmp.language.chinese",
		Translations: "中文",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "zh-CHS",
		Key:          "mvpmp.language.chinese",
		Translations: "英文",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "en-US",
		Key:          "agri.login.success.email",
		Translations: "英文",
	})
	return arr
}

func GetMockTranslations() []*models.Translation {
	arr := make([]*models.Translation, 0)
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Code:         "en-US",
		Key:          "agri.login.success.email",
		Translations: "agri.login.success.email",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Translations: "agri.email.header",
		Key:          "agri.email.header",
		Code:         "en-US",
	})
	arr = append(arr, &models.Translation{
		Id:           GetRandomUUID(),
		Translations: "agri.email.footer",
		Key:          "agri.email.footer",
		Code:         "en-US",
	})
	return arr
}

func GetMockTranslationResp() map[string]string {
	return map[string]string{
		"agri.email.footer":        "agri.email.footer",
		"agri.email.header":        "agri.email.header",
		"agri.login.success.email": "agri.login.success.email",
	}
}

func GetRandomUUID() uuid.UUID {
	uniqueId, _ := uuid.NewRandom()
	return uniqueId
}
