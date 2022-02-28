package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go-thai-dialect/models"

	"github.com/labstack/echo"
)

func getProvince(c echo.Context) error {

	province_list := models.GetAllProvince()

	return c.JSON(http.StatusOK, province_list)
}

func getDistrict(c echo.Context) error {

	province_id := c.QueryParam("province_id")

	district_list := models.GetDistrictByProvinceID(province_id)

	return c.JSON(http.StatusOK, district_list)
}

func postVolunteer(c echo.Context) error {

	var volunteer models.Volunteer

	if err := c.Bind(&volunteer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status":  false,
			"message": "Invalid request",
			"error":   err.Error(),
		})
	}

	b, err := json.Marshal(volunteer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))

	// fmt.Println(volunteer)

	if models.CheckExistVolunteer(volunteer) {
		create_status, volunteer := models.CreateVolunteer(volunteer)

		b, err := json.Marshal(volunteer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))

		if create_status {
			return c.JSON(http.StatusOK, volunteer)
		} else {
			return c.JSON(http.StatusOK, echo.Map{
				"status":  false,
				"message": "Create volunteer fail",
			})
		}
	}

	return c.JSON(http.StatusBadRequest, echo.Map{
		"status":  false,
		"message": "Something went wrong",
	})
}

func VolunteerDock(_echo *echo.Group) {

	_echo.GET("/volunteer/province", getProvince)
	_echo.POST("/volunteer", postVolunteer)
	_echo.GET("/volunteer/district/", getDistrict)

}
