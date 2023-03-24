package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/CyrilSbrodov/passManager.git/server/cmd/config"
	"github.com/CyrilSbrodov/passManager.git/server/cmd/loggers"
	"github.com/CyrilSbrodov/passManager.git/server/internal/mocks"
	"github.com/CyrilSbrodov/passManager.git/server/internal/models"
	"github.com/CyrilSbrodov/passManager.git/server/internal/storage/repositories"
	"github.com/CyrilSbrodov/passManager.git/server/pkg/client/postgres"
)

var (
	CFG config.Config
)

func TestMain(m *testing.M) {
	cfg := config.ConfigInit()
	CFG = *cfg
	os.Exit(m.Run())
}

func newRepo() *repositories.Store {
	logger := loggers.NewLogger()
	client, _ := postgres.NewClient(context.Background(), 5, &CFG, logger)
	repo, _ := repositories.NewStore(client, &CFG, logger)
	return repo
}

func TestHandler_Registration(t *testing.T) {
	tests := []struct {
		name         string
		body         models.User
		answerID     interface{}
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.User{
				Login:    "test",
				Password: "123456",
			},
			answerID:     "1",
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 409",
			body: models.User{
				Login:    "test",
				Password: "123456",
			},
			answerID:     "",
			answerError:  errors.New("err"),
			expectedCode: http.StatusConflict,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(bodyJSON))
			s.EXPECT().Register(gomock.Any()).Return(tt.answerID, tt.answerError)
			h.Registration().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name         string
		body         models.User
		answerID     interface{}
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.User{
				Login:    "test",
				Password: "123456",
			},
			answerID:     "1",
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 401",
			body: models.User{
				Login:    "test",
				Password: "123456",
			},
			answerID:     "",
			answerError:  errors.New("err"),
			expectedCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(bodyJSON))
			s.EXPECT().Login(gomock.Any()).Return(tt.answerID, tt.answerError)
			h.Login().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_CollectCards(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoCard
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoCard{
				Name:   []byte("test"),
				Number: []byte("123456"),
				CVC:    []byte("123"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoCard{
				Name:   []byte("test"),
				Number: []byte("1234256"),
				CVC:    []byte("123"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/cards", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().CollectCard(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.CollectCards().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetCards(t *testing.T) {
	tests := []struct {
		name         string
		answerCode   int
		answerError  error
		answerData   []models.CryptoCard
		expectedCode int
	}{
		{
			name:        "Test ok",
			answerCode:  200,
			answerError: nil,
			answerData: []models.CryptoCard{
				{
					UID:    1,
					Number: nil,
					Name:   nil,
					CVC:    nil,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test 500",
			answerCode:   500,
			answerError:  errors.New("err"),
			answerData:   nil,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/data/cards", nil)
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().GetCards(gomock.Any()).Return(tt.answerCode, tt.answerData, tt.answerError)
			h.GetCards().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_UpdateCards(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoCard
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoCard{
				Name:   []byte("test"),
				Number: []byte("123456"),
				CVC:    []byte("123"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoCard{
				Name:   []byte("test"),
				Number: []byte("1234256"),
				CVC:    []byte("123"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/update/cards", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().UpdateCard(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.UpdateCards().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_DeleteCards(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoCard
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoCard{
				UID:    1,
				Name:   []byte("test"),
				Number: []byte("123456"),
				CVC:    []byte("123"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoCard{
				Name:   []byte("test"),
				Number: []byte("1234256"),
				CVC:    []byte("123"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/delete/cards", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().DeleteCard(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.DeleteCards().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_CollectPassword(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoPassword
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoPassword{
				Login: []byte("test"),
				Pass:  []byte("123"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoPassword{
				Login: []byte("test"),
				Pass:  []byte("1234256"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/password", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().CollectPassword(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.CollectPassword().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetPasswords(t *testing.T) {
	tests := []struct {
		name         string
		answerCode   int
		answerError  error
		answerData   []models.CryptoPassword
		expectedCode int
	}{
		{
			name:        "Test ok",
			answerCode:  200,
			answerError: nil,
			answerData: []models.CryptoPassword{
				{
					UID:   1,
					Login: nil,
					Pass:  nil,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test 500",
			answerCode:   500,
			answerError:  errors.New("err"),
			answerData:   nil,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/data/passwords", nil)
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().GetPassword(gomock.Any()).Return(tt.answerCode, tt.answerData, tt.answerError)
			h.GetPasswords().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_UpdatePassword(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoPassword
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoPassword{
				Login: []byte("test"),
				Pass:  []byte("123456"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoPassword{
				Login: []byte("test"),
				Pass:  []byte("1234256"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/update/password", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().UpdatePassword(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.UpdatePassword().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_DeletePassword(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoPassword
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoPassword{
				UID:   1,
				Login: []byte("test"),
				Pass:  []byte("123456"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoPassword{
				Login: []byte("test"),
				Pass:  []byte("1234256"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/delete/password", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().DeletePassword(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.DeletePassword().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_CollectText(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoTextData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoTextData{
				Text: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoTextData{
				Text: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/text", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().CollectText(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.CollectText().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetText(t *testing.T) {
	tests := []struct {
		name         string
		answerCode   int
		answerError  error
		answerData   []models.CryptoTextData
		expectedCode int
	}{
		{
			name:        "Test ok",
			answerCode:  200,
			answerError: nil,
			answerData: []models.CryptoTextData{
				{
					UID:  1,
					Text: nil,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test 500",
			answerCode:   500,
			answerError:  errors.New("err"),
			answerData:   nil,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/data/text", nil)
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().GetText(gomock.Any()).Return(tt.answerCode, tt.answerData, tt.answerError)
			h.GetText().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_UpdateText(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoTextData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoTextData{
				Text: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoTextData{
				Text: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/update/text", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().UpdateText(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.UpdateText().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_DeleteText(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoTextData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoTextData{
				UID:  1,
				Text: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoTextData{
				Text: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/delete/text", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().DeleteText(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.DeleteText().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_CollectBinary(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoBinaryData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoBinaryData{
				Data: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoBinaryData{
				Data: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/binary", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().CollectBinary(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.CollectBinary().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_GetBinary(t *testing.T) {
	tests := []struct {
		name         string
		answerCode   int
		answerError  error
		answerData   []models.CryptoBinaryData
		expectedCode int
	}{
		{
			name:        "Test ok",
			answerCode:  200,
			answerError: nil,
			answerData: []models.CryptoBinaryData{
				{
					UID:  1,
					Data: nil,
				},
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "Test 500",
			answerCode:   500,
			answerError:  errors.New("err"),
			answerData:   nil,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/data/binary", nil)
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().GetBinary(gomock.Any()).Return(tt.answerCode, tt.answerData, tt.answerError)
			h.GetBinary().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_UpdateBinary(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoBinaryData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoBinaryData{
				Data: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoBinaryData{
				Data: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/update/binary", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().UpdateBinary(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.UpdateBinary().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}

func TestHandler_DeleteBinary(t *testing.T) {
	tests := []struct {
		name         string
		body         models.CryptoBinaryData
		answerCode   int
		answerError  error
		expectedCode int
	}{
		{
			name: "Test ok",
			body: models.CryptoBinaryData{
				UID:  1,
				Data: []byte("test"),
			},
			answerCode:   200,
			answerError:  nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "Test 500",
			body: models.CryptoBinaryData{
				Data: []byte("test"),
			},
			answerCode:   500,
			answerError:  errors.New("err"),
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			s := mocks.NewMockStorage(ctrl)
			logger := loggers.NewLogger()
			h := &Handler{
				Storage: s,
				logger:  *logger,
			}

			bodyJSON, err := json.Marshal(tt.body)
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/data/delete/binary", bytes.NewBuffer(bodyJSON))
			req = req.WithContext(context.WithValue(context.Background(), "user_id", "1"))
			s.EXPECT().DeleteBinary(gomock.Any(), gomock.Any()).Return(tt.answerCode, tt.answerError)
			h.DeleteBinary().ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
