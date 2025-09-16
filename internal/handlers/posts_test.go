package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RomanKovalev007/mai_news/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllPostsHandler(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockPoster)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			mockSetup: func(mp *MockPoster) {
				mp.On("GetAllPosts").Return([]models.OutputPost{
					{ID: 1, Title: "Test Post 1", Content: "Content 1", CreatedAt: ""},
					{ID: 2, Title: "Test Post 2", Content: "Content 2", CreatedAt: ""},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"title":"Test Post 1","content":"Content 1","created_at":""},{"id":2,"title":"Test Post 2","content":"Content 2","created_at":""}]`+"\n",
		},
		{
			name: "not found",
			mockSetup: func(mp *MockPoster) {
				mp.On("GetAllPosts").Return([]models.OutputPost{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Post not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPoster := NewMockPoster(t)
			tt.mockSetup(mockPoster)

			handler := GetAllPostsHandler(mockPoster, slog.Default())
			req := httptest.NewRequest("GET", "/posts", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockPoster.AssertExpectations(t)
		})
	}
}

func TestGetPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockSetup      func(*MockPoster)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			postID: "1",
			mockSetup: func(mp *MockPoster) {
				mp.On("GetPost", 1).Return(models.OutputPost{
					ID: 1, Title: "Test Post", Content: "Test Content",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Test Post","content":"Test Content","created_at":""}` + "\n",
		},
		{
			name:   "invalid id",
			postID: "invalid",
			mockSetup: func(mp *MockPoster) {
				// No mock setup needed for invalid ID
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid post ID\n",
		},
		{
			name:   "not found",
			postID: "999",
			mockSetup: func(mp *MockPoster) {
				mp.On("GetPost", 999).Return(models.OutputPost{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Post not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPoster := NewMockPoster(t)
			tt.mockSetup(mockPoster)

			handler := GetPostHandler(mockPoster, slog.Default())
			req := httptest.NewRequest("GET", "/posts/"+tt.postID, nil)
			req.SetPathValue("id", tt.postID)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockPoster.AssertExpectations(t)
		})
	}
}

func TestCreatePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockPoster)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			requestBody: models.InputPost{
				Title:   "New Post",
				Content: "New Content",
			},
			mockSetup: func(mp *MockPoster) {
				mp.On("SavePost", mock.AnythingOfType("models.InputPost")).Return(models.OutputPost{
					ID: 1, Title: "New Post", Content: "New Content",
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":1,"title":"New Post","content":"New Content","created_at":""}` + "\n",
		},
		{
			name:           "invalid json",
			requestBody:    `invalid json`,
			mockSetup:      func(mp *MockPoster) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request payload\n",
		},
		{
			name: "save error",
			requestBody: models.InputPost{
				Title:   "Error Post",
				Content: "Error Content",
			},
			mockSetup: func(mp *MockPoster) {
				mp.On("SavePost", mock.AnythingOfType("models.InputPost")).Return(models.OutputPost{}, errors.New("save error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "failed to save post\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPoster := NewMockPoster(t)
			tt.mockSetup(mockPoster)

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			handler := CreatePostHandler(mockPoster, slog.Default())
			req := httptest.NewRequest("POST", "/posts", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockPoster.AssertExpectations(t)
		})
	}
}

func TestPatchPostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		requestBody    interface{}
		mockSetup      func(*MockPoster)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			postID: "1",
			requestBody: models.InputPost{
				Title:   "Updated Post",
				Content: "Updated Content",
			},
			mockSetup: func(mp *MockPoster) {
				mp.On("PatchPost", 1, mock.AnythingOfType("models.InputPost")).Return(models.OutputPost{
					ID: 1, Title: "Updated Post", Content: "Updated Content",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"title":"Updated Post","content":"Updated Content","created_at":""}` + "\n",
		},
		{
			name:           "invalid id",
			postID:         "invalid",
			requestBody:    models.InputPost{Title: "Test"},
			mockSetup:      func(mp *MockPoster) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid post ID\n",
		},
		{
			name:           "invalid json",
			postID:         "1",
			requestBody:    `invalid json`,
			mockSetup:      func(mp *MockPoster) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request payload\n",
		},
		{
			name:   "not found",
			postID: "999",
			requestBody: models.InputPost{
				Title:   "Non-existent",
				Content: "Content",
			},
			mockSetup: func(mp *MockPoster) {
				mp.On("PatchPost", 999, mock.AnythingOfType("models.InputPost")).Return(models.OutputPost{}, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Post not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPoster := NewMockPoster(t)
			tt.mockSetup(mockPoster)

			var bodyBytes []byte
			switch v := tt.requestBody.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			handler := PatchPostHandler(mockPoster, slog.Default())
			req := httptest.NewRequest("PATCH", "/posts/"+tt.postID, bytes.NewReader(bodyBytes))
			req.SetPathValue("id", tt.postID)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockPoster.AssertExpectations(t)
		})
	}
}

func TestDeletePostHandler(t *testing.T) {
	tests := []struct {
		name           string
		postID         string
		mockSetup      func(*MockPoster)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "success",
			postID: "1",
			mockSetup: func(mp *MockPoster) {
				mp.On("DeletePost", 1).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "invalid id",
			postID:         "invalid",
			mockSetup:      func(mp *MockPoster) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid post ID\n",
		},
		{
			name:   "not found",
			postID: "999",
			mockSetup: func(mp *MockPoster) {
				mp.On("DeletePost", 999).Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Post not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPoster := NewMockPoster(t)
			tt.mockSetup(mockPoster)

			handler := DeletePostHandler(mockPoster, slog.Default())
			req := httptest.NewRequest("DELETE", "/posts/"+tt.postID, nil)
			req.SetPathValue("id", tt.postID)
			w := httptest.NewRecorder()

			handler(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
			mockPoster.AssertExpectations(t)
		})
	}
}