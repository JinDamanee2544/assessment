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
	Tags   []string `json:"tags"`
}

type Error struct {
	Message string `json:"message"`
}

func PostExpense(c echo.Context) error {
	e := Expense{}

	var err error
	err = c.Bind(&e)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	row := db.QueryRow("INSERT INTO expenses (title, amount, note, tags) VALUES ($1, $2, $3, $4) RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))
	err = row.Scan(&e.ID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}

func GetAllExpense(c echo.Context) error {
	rows, err := db.Query("SELECT * FROM expenses")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	expenses := []Expense{}
	for rows.Next() {
		e := Expense{}
		err = rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

		if err != nil {
			return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
		}

		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}
