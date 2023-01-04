package expense

import (
	"bytes"
	"encoding/json"
	"io"
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

type sqlCommandMock struct {
	sqlCommand string
	sqlResult  *sqlmock.Rows
}

func setUpContext(body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TOKEN"))

	e := echo.New()
	c := e.NewContext(req, rec)

	return c, rec
}

func seedExpense(t *testing.T, mock sqlmock.Sqlmock) *Expense {
	// Similiar to Test CreateExpense
	var seed = Expense{
		Title:  "strawberry smoothie C++",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	sqlCreate := sqlCommandMock{
		sqlCommand: "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id",
		sqlResult:  sqlmock.NewRows([]string{"id"}).AddRow(1),
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlCreate.sqlCommand)).
		WithArgs(seed.Title, seed.Amount, seed.Note, pq.Array(seed.Tags)).
		WillReturnRows(sqlCreate.sqlResult)

	body := bytes.NewBufferString(
		`{
			"title": "strawberry smoothie C++",
			"amount": 79,
			"note": "night market promotion discount 10 bath",
			"tags": ["food", "beverage"]
		}`)

	c, rec := setUpContext(body)
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

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	return &ex
}

func TestCreateExpense(t *testing.T) {
	mock, closeDB := InitMockDB(t)
	defer closeDB()

	var seed = Expense{
		Title:  "strawberry smoothie C++",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	sqlCreate := sqlCommandMock{
		sqlCommand: "INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id",
		sqlResult:  sqlmock.NewRows([]string{"id"}).AddRow(1),
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlCreate.sqlCommand)).
		WithArgs(seed.Title, seed.Amount, seed.Note, pq.Array(seed.Tags)).
		WillReturnRows(sqlCreate.sqlResult)

	// ----------------------------

	body := bytes.NewBufferString(
		`{
		"title": "strawberry smoothie C++",
		"amount": 79,
		"note": "night market promotion discount 10 bath",
		"tags": ["food", "beverage"]
	}`)

	c, rec := setUpContext(body)
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

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetExpenseByID(t *testing.T) {
	mock, closeDB := InitMockDB(t)
	defer closeDB()

	seed := seedExpense(t, mock)

	var want = Expense{
		ID:     "1",
		Title:  "strawberry smoothie C++",
		Amount: 79,
		Note:   "night market promotion discount 10 bath",
		Tags:   []string{"food", "beverage"},
	}

	sqlGet := sqlCommandMock{
		sqlCommand: "SELECT * FROM expenses WHERE id = $1",
		sqlResult: sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
			AddRow(want.ID, want.Title, want.Amount, want.Note, pq.Array(want.Tags)),
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlGet.sqlCommand)).
		WithArgs(seed.ID).
		WillReturnRows(sqlGet.sqlResult)

	// ----------------------------

	c, rec := setUpContext(nil)
	c.SetParamNames("id")
	c.SetParamValues(seed.ID)
	err := GetExpenseByID(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	ex := Expense{}
	err = json.NewDecoder(rec.Body).Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rec.Code)
	assert.Equal(t, seed.ID, ex.ID)
	assert.Equal(t, seed.Title, ex.Title)
	assert.Equal(t, seed.Amount, ex.Amount)
	assert.Equal(t, seed.Note, ex.Note)
	assert.Equal(t, seed.Tags, ex.Tags)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateExpenseByID(t *testing.T) {
	mock, closeDB := InitMockDB(t)
	defer closeDB()

	seed := seedExpense(t, mock)

	var changeTo = Expense{
		Title:  "apple smoothie",
		Amount: 89,
		Note:   "no discount",
		Tags:   []string{"beverage"},
	}

	sqlUpdate := sqlCommandMock{
		sqlCommand: "UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id",
		sqlResult:  sqlmock.NewRows([]string{"id"}).AddRow(1),
	}

	mock.ExpectQuery(regexp.QuoteMeta(sqlUpdate.sqlCommand)).
		WithArgs(changeTo.Title, changeTo.Amount, changeTo.Note, pq.Array(changeTo.Tags), seed.ID).
		WillReturnRows(sqlUpdate.sqlResult)

	// ----------------------------

	body := bytes.NewBufferString(
		`{
			"title": "apple smoothie",
			"amount": 89,
			"note": "no discount",
			"tags": ["beverage"]
		}`)

	c, rec := setUpContext(body)
	c.SetParamNames("id")
	c.SetParamValues(seed.ID)
	err := UpdateExpenseByID(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	ex := Expense{}
	err = json.NewDecoder(rec.Body).Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, rec.Code)
	assert.Equal(t, seed.ID, ex.ID)
	assert.Equal(t, changeTo.Title, ex.Title)
	assert.Equal(t, changeTo.Amount, ex.Amount)
	assert.Equal(t, changeTo.Note, ex.Note)
	assert.Equal(t, changeTo.Tags, ex.Tags)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
