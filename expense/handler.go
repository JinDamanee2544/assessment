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

func CreateExpense(c echo.Context) error {
	if c.Request().ContentLength == 0 {
		return c.JSON(http.StatusBadRequest, Error{Message: "No body found"})
	}

	var err error
	e := Expense{}
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

// For testing purpose
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

func GetExpenseByID(c echo.Context) error {
	id := c.Param("id")

	row := db.QueryRow("SELECT * FROM expenses WHERE id = $1", id)

	e := Expense{}
	err := row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}

func UpdateExpenseByID(c echo.Context) error {
	id := c.Param("id")

	e := Expense{}
	err := c.Bind(&e)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Error{Message: err.Error()})
	}

	row := db.QueryRow("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id", e.Title, e.Amount, e.Note, pq.Array(&e.Tags), id)

	err = row.Scan(&e.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Error{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
