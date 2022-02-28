package controllers

import (
	"go-thai-dialect/helper"
	"go-thai-dialect/models"
	"net/http"

	"github.com/labstack/echo"
)

func postUser(c echo.Context) error {
	user := models.User{}
	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	user.RegisterDate = helper.CurrentTime()

	user_id, err := models.InsertUser(user)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"status":  false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status":  true,
		"message": "Create user success",
		"user_id": user_id,
	})

}

func UserDock(_echo *echo.Group) {

	_echo.POST("/user", postUser)

}
