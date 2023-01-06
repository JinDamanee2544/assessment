//go:build integration

package expense

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func validateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")

		correctToken := os.Getenv("TOKEN")
		if correctToken == "" {
			log.Fatal("TOKEN is not set")
		}

		if token != os.Getenv("TOKEN") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		return next(c)
	}
}

func seedExpenseIT(t *testing.T) (*Expense, func()) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 89,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
		}`)

	e := Expense{}
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)
	if err != nil {
		panic(err)
	}

	deleteSeed := func() {
		res := request(http.MethodDelete, uri("expenses", e.ID), nil)
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
	}
	return &e, deleteSeed
}

func TestITPostExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"id": "1",
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
		}`)

	e := Expense{}
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.NotEqualValues(t, 0, e.ID)
	assert.EqualValues(t, "strawberry smoothie", e.Title)
	assert.EqualValues(t, 79, e.Amount)
	assert.EqualValues(t, "night market promotion discount 10 bath", e.Note)
	assert.EqualValues(t, []string{"food", "beverage"}, e.Tags)

	// clean up
	res = request(http.MethodDelete, uri("expenses", e.ID), nil)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
}

func TestITPostExpenseNoBody(t *testing.T) {
	body := bytes.NewBufferString("")

	e := Expense{}
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, res.StatusCode)
}

func TestITGetExpenseByID(t *testing.T) {
	insertE, deleteSeed := seedExpenseIT(t)
	defer deleteSeed()

	res := request(http.MethodGet, uri("expenses", insertE.ID), nil)

	e := Expense{}
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, insertE.ID, e.ID)
	assert.EqualValues(t, insertE.Title, e.Title)
	assert.EqualValues(t, insertE.Amount, e.Amount)
	assert.EqualValues(t, insertE.Note, e.Note)
	assert.EqualValues(t, insertE.Tags, e.Tags)
}

func TestITUpdateExpenseByID(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)

	insertE, deleteSeed := seedExpenseIT(t)
	defer deleteSeed()

	res := request(http.MethodPut, uri("expenses", insertE.ID), body)

	editedEx := Expense{}
	err := res.Decode(&editedEx)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, insertE.ID, editedEx.ID)
	assert.EqualValues(t, "apple smoothie", editedEx.Title)
	assert.EqualValues(t, 89, editedEx.Amount)
	assert.EqualValues(t, "no discount", editedEx.Note)
	assert.EqualValues(t, []string{"beverage"}, editedEx.Tags)
}

func TestITGetAllExpenses(t *testing.T) {
	checkRes := request(http.MethodGet, uri("expenses"), nil)

	// Check former expenses length (offset)
	formerEx := []Expense{}
	err := checkRes.Decode(&formerEx)
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, checkRes.StatusCode)
	formerLen := len(formerEx)

	insertE, deleteSeed := seedExpenseIT(t)
	defer deleteSeed()

	insertE2, deleteSeed2 := seedExpenseIT(t)
	defer deleteSeed2()

	res := request(http.MethodGet, uri("expenses"), nil)

	e := []Expense{}
	err = res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.EqualValues(t, formerLen+2, len(e))
	assert.Contains(t, e, *insertE)
	assert.Contains(t, e, *insertE2)
}
