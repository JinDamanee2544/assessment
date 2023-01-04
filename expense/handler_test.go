package expense

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func InitMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err.Error())
	}
	setMockDB(db)

	closeDB := func() {
		db.Close()
	}

	return mock, closeDB
}

// with http/test
func TestCreateExpense(t *testing.T) {
	mock, closeDB := InitMockDB(t)
	defer closeDB()

	var seed = Expense{
		Title:  "strawberry smoothie C++",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	sqlCreate := "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id"
	sqlResult := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta(sqlCreate)).
		WithArgs(seed.Title, seed.Amount, seed.Note, pq.Array(seed.Tags)).
		WillReturnRows((sqlResult))

	// ----------------------------

	body := bytes.NewBufferString(
		`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/expense", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TOKEN"))

	e := echo.New()
	c := e.NewContext(req, rec)
	err := CreateExpense(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	ex := Expense{}
	err = json.NewDecoder(rec.Body).Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, rec.Code)
	assert.NotEqual(t, seed.ID, ex.ID)
	assert.Equal(t, seed.Title, ex.Title)
	assert.Equal(t, seed.Amount, ex.Amount)
	assert.Equal(t, seed.Note, ex.Note)
	assert.Equal(t, seed.Tags, ex.Tags)
}
