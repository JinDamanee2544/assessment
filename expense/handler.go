package expense

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Expense struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Amount int      `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tag"`
}

func PostExpense(c echo.Context) error {
	e := Expense{}

	var err error
	err = c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	newE := Expense{}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&newE.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, newE)
}
