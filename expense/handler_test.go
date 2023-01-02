package expense

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// with http/test
func TestPostExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "http://localhost:2565/expense", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := PostExpense(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	ex := Expense{}
	err = json.NewDecoder(rec.Body).Decode(&ex)

	assert.Nil(t, err)
	assert.Equal(t, "1", ex.ID)
	assert.Equal(t, "strawberry smoothie", ex.Title)
	assert.Equal(t, 79, ex.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", ex.Note)
	assert.Equal(t, []string{"food", "beverage"}, ex.Tags)
}
