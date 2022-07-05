package controllers

import (
	"balance/pkg/db"
	"balance/pkg/models"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func SendMoney(c echo.Context) error {
	data := &models.Balance{}
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(fmt.Sprintf("Bind failed: %v", err)))
	}

	balance, err := db.SendMoney(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(fmt.Sprintf("Couldn't execute function: %v", err)))
	}

	u := &models.Balance{
		UserId:  data.UserId,
		Balance: balance,
	}
	return c.JSON(http.StatusOK, u)
}

func PayBetweenUsers(c echo.Context) error {
	data := &models.Balance{}
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(fmt.Sprintf("Bind failed: %v", err)))
	}

	list, err := db.UpdateTables(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse("failed: %v", err))
	}

	return c.JSON(http.StatusOK, list)
}

func SelectBalance(c echo.Context) error {
	data := &models.Balance{}
	err := c.Bind(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.NewErrorResponse(fmt.Sprintf("Bind failed: %v", err)))
	}

	balance, err := db.SelectBalance(data.UserId)
	if err != nil {
		return err
	}

	user := &models.Balance{
		UserId:  data.UserId,
		Balance: balance,
	}
	return c.JSON(http.StatusOK, user)

}
