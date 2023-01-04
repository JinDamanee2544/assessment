//go:build integration

package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func uri(paths ...string) string {
	host := "http://localhost:2565"

	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)

	req.Header.Set("Content-Type", "application/json")

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("TOKEN is not set")
	}
	req.Header.Set("Authorization", os.Getenv("TOKEN"))

	client := &http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func seedExpense() *Expense {
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
	return &e
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
	insertE := seedExpense()

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

	insertE := seedExpense()

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

func TestGetAllExpenses(t *testing.T) {
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
