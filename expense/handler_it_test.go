package expense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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

func TestPostExpense(t *testing.T) {
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
