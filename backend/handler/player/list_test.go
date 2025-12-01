package player

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/importandrew/MaCRa/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testData := map[string]struct {
		queryParams  map[string]string
		mockSetup    func(*MockPlayerService)
		expectedCode int
		expectedBody func(*testing.T, *httptest.ResponseRecorder)
	}{
		"ok with name and limit": {
			queryParams: map[string]string{
				"name":  "test",
				"limit": "5",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().ListByNameLimit(
					mock.Anything,
					"test",
					5,
				).Return([]models.Player{
					{ID: 1, Name: "Test User 1"},
					{ID: 2, Name: "Test User 2"},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var players []models.Player
				err := json.Unmarshal(w.Body.Bytes(), &players)
				assert.NoError(t, err)
				assert.Len(t, players, 2)
				assert.Equal(t, 1, players[0].ID)
				assert.Equal(t, "Test User 1", players[0].Name)
				assert.Equal(t, 2, players[1].ID)
				assert.Equal(t, "Test User 2", players[1].Name)
			},
		},
		"ok with name only": {
			queryParams: map[string]string{
				"name": "john",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().ListByNameLimit(
					mock.Anything,
					"john",
					0,
				).Return([]models.Player{
					{ID: 3, Name: "John Doe"},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var players []models.Player
				err := json.Unmarshal(w.Body.Bytes(), &players)
				assert.NoError(t, err)
				assert.Len(t, players, 1)
				assert.Equal(t, 3, players[0].ID)
				assert.Equal(t, "John Doe", players[0].Name)
			},
		},
		"ok with limit only": {
			queryParams: map[string]string{
				"limit": "10",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().ListByNameLimit(
					mock.Anything,
					"",
					10,
				).Return([]models.Player{}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var players []models.Player
				err := json.Unmarshal(w.Body.Bytes(), &players)
				assert.NoError(t, err)
				assert.Empty(t, players)
			},
		},
		"ok empty query": {
			queryParams: map[string]string{},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().ListByNameLimit(
					mock.Anything,
					"",
					0,
				).Return([]models.Player{}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var players []models.Player
				err := json.Unmarshal(w.Body.Bytes(), &players)
				assert.NoError(t, err)
				assert.Empty(t, players)
			},
		},
		"invalid limit": {
			queryParams: map[string]string{
				"name":  "test",
				"limit": "not-a-number",
			},
			mockSetup:    func(m *MockPlayerService) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		"service error": {
			queryParams: map[string]string{
				"name":  "test",
				"limit": "5",
			},
			mockSetup: func(m *MockPlayerService) {
				m.EXPECT().ListByNameLimit(
					mock.Anything,
					"test",
					5,
				).Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "database error", response["error"])
			},
		},
	}

	for name, testCase := range testData {
		t.Run(name, func(t *testing.T) {
			mockPlayerService := NewMockPlayerService(t)
			testCase.mockSetup(mockPlayerService)

			handler := NewHandler(mockPlayerService, "", nil)

			req := httptest.NewRequest(http.MethodGet, "/api/players/list", nil)
			q := req.URL.Query()
			for key, value := range testCase.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			handler.List(ctx)

			assert.Equal(t, testCase.expectedCode, w.Code)
			if testCase.expectedBody != nil {
				testCase.expectedBody(t, w)
			}
		})
	}
}
