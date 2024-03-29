package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"crud/config"
	"crud/models"

	"github.com/stretchr/testify/assert"
)

type BookResponse struct {
	Message string         `json:"message"`
	Data    []models.Books `json:"data"`
}

type DeleteBookResponse struct {
	Message string `json:"message"`
}

func InsertDataBookForGetBooks() error {
	book := models.Books{
		Judul:    "informatika",
		Penulis:  "ipsum dolor",
		Penerbit: "lorem",
	}

	var err error
	if err = config.DB.Save(&book).Error; err != nil {
		return err
	}
	return nil
}

func TestGetBooksController(t *testing.T) {
	var testCases = []struct {
		name       string
		path       string
		expectCode int
		sizeData   int
	}{
		{
			name:       "get book normal",
			path:       "/books",
			expectCode: http.StatusOK,
			sizeData:   1,
		},
	}

	e := InitEchoTestAPI()
	InsertDataBookForGetBooks()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	for _, testCase := range testCases {
		c.SetPath(testCase.path)

		if assert.NoError(t, GetBooksController(c)) {
			assert.Equal(t, testCase.expectCode, rec.Code)
			body := rec.Body.String()

			var book BookResponse
			err := json.Unmarshal([]byte(body), &book)

			if err != nil {
				assert.Error(t, err, "error")
			}
			assert.Equal(t, testCase.sizeData, len(book.Data))
		}

	}

}

func TestGetBookController(t *testing.T) {

	book := models.Books{
		Judul:    "Alta",
		Penulis:  "ipsum dolor",
		Penerbit: "lorem",
	}
	if err := config.DB.Save(&book).Error; err != nil {
		t.Errorf("Error setting up test book: %v", err)
	}

	e := InitEchoTestAPI()
	InsertDataBookForGetBooks()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/books/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(1))

	if assert.NoError(t, GetBookController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var respBody struct {
			Message string       `json:"message"`
			Data    models.Books `json:"data"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &respBody)
		if err != nil {
			t.Errorf("Error unmarshalling response body: %v", err)
		}
		assert.Equal(t, "get book successfully", respBody.Message)
		assert.Equal(t, book.Judul, respBody.Data.Judul)
		assert.Equal(t, book.Penulis, respBody.Data.Penulis)
		assert.Equal(t, book.Penerbit, respBody.Data.Penerbit)
	}

	if err := config.DB.Delete(&book).Error; err != nil {
		t.Errorf("Error cleaning up test book: %v", err)
	}
}
func TestCreateBookController(t *testing.T) {
	var testCases = []struct {
		name       string
		request    string
		expectCode int
	}{
		{
			name:       "create book success",
			request:    `{"title":"Test","author":"Rifal","publisher":"Luthfi"}`,
			expectCode: http.StatusOK,
		},
	}

	e := InitEchoTestAPI()
	for _, testCase := range testCases {
		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(testCase.request))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if assert.NoError(t, CreateBookController(c)) {
			assert.Equal(t, testCase.expectCode, rec.Code)

			if testCase.expectCode == http.StatusCreated {
				body := rec.Body.String()
				var book models.Books
				err := json.Unmarshal([]byte(body), &book)

				if err != nil {
					assert.Error(t, err, "error")
				}
				assert.Equal(t, testCase.request[0], book.Judul)
				assert.Equal(t, testCase.request[2], book.Penerbit)
			}
		}
		// if assert.Error(t, CreateBookController(c)) {
		// 	t.Fatal("error")
		// }
	}
}
func TestDeleteBookController(t *testing.T) {
	e := InitEchoTestAPI()
	InsertDataBookForGetBooks()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, DeleteBookController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()

		var result DeleteBookResponse
		err := json.Unmarshal([]byte(body), &result)

		if err != nil {
			assert.Error(t, err, "error")
		}

		assert.Equal(t, "book deleted successfully", result.Message)
	}
}
func TestUpdateBookController(t *testing.T) {
	e := InitEchoTestAPI()

	book := models.Books{
		Judul:    "Alta",
		Penulis:  "ipsum dolor",
		Penerbit: "lorem",
	}
	if err := config.DB.Save(&book).Error; err != nil {
		t.Errorf("failed to create test book: %v", err)
	}

	requestBody := map[string]interface{}{
		"title": "Jane Doe",
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Errorf("failed to marshal request body: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(int(book.ID)), bytes.NewBuffer(requestBodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(int(book.ID)))

	if assert.NoError(t, UpdateBookController(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var bookResponse models.Books
		err = json.Unmarshal(requestBodyBytes, &bookResponse)
		if err != nil {
			t.Errorf("failed to unmarshal response body: %v", err)
		}
		assert.Equal(t, requestBody["title"], bookResponse.Judul)
	}
}
