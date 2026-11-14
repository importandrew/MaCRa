package player

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/importandrew/MaCRa/config"
	"github.com/importandrew/MaCRa/models"
)

func TestHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		url          string
		mockSetup    func(*MockPlayerService)
		expectedCode int
		want         []models.Player
	}{
		"the query reaches the search": {
			url: "/players?q=Iv",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Search(mock.Anything, "Iv").Return([]models.Player{{ID: 42, Name: "Ivan"}}, nil)
			},
			expectedCode: http.StatusOK,
			want:         []models.Player{{ID: 42, Name: "Ivan"}},
		},
		"no query": {
			url: "/players",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Search(mock.Anything, "").Return([]models.Player{}, nil)
			},
			expectedCode: http.StatusOK,
			want:         []models.Player{},
		},
		"service error": {
			url: "/players?q=Ivan",
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().Search(mock.Anything, "Ivan").Return(nil, errors.New("service error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockPlayer := NewMockPlayerService(t)
			testCase.mockSetup(mockPlayer)
			handler := NewHandler(mockPlayer, config.Config{})

			router := gin.New()
			router.GET("/players", handler.List)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, testCase.url, nil))

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.want != nil {
				var players []models.Player
				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &players))
				assert.Equal(t, testCase.want, players)
			}
		})
	}
}
