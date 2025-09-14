package functions

import (
	"net/http"

	"github.com/U-T-kuroitigo/Saikyo_UI/configuration"
	"github.com/U-T-kuroitigo/Saikyo_UI/models"
	"github.com/U-T-kuroitigo/Saikyo_UI/response"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		r := response.Model{
			Code:    "400",
			Message: "Incorrect structure",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if err := models.ValidateUser(user); err != nil {
		return c.JSON(http.StatusBadRequest, response.Model{
			Code:    "400",
			Message: "Failed validation",
			Data:    err.Error(),
		})
	}

	db := configuration.GetDB()
	if err := db.Create(&user).Error; err != nil {
		r := response.Model{
			Code:    "500",
			Message: "Error creating user",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	r := response.Model{
		Code:    "201",
		Message: "Created Successfully",
		Data:    user,
	}
	return c.JSON(http.StatusCreated, r)
}

func GetAllUsers(c echo.Context) error {
	users := []models.User{}
	db := configuration.GetDB()
	if err := db.Find(&users).Error; err != nil {
		r := response.Model{
			Code:    "500",
			Message: "Query error",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	r := response.Model{
		Code:    "200",
		Message: "Correctly consulted",
		Data:    users,
	}
	return c.JSON(http.StatusOK, r)
}

func DeleteUser(c echo.Context) error {
	var user models.User
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.Model{
				Code:    "404",
				Message: "User not found",
				Data:    err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Model{
				Code:    "500",
				Message: "Query error",
				Data:    err.Error(),
			})
		}
	}

	if err := db.Delete(&user).Error; err != nil {
		r := response.Model{
			Code:    "500",
			Message: "Delete error",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	r := response.Model{
		Code:    "202",
		Message: "Correctly Deleted",
		Data:    user,
	}
	return c.JSON(http.StatusAccepted, r)
}

func UpdateUser(c echo.Context) error {
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	user := models.User{}
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.Model{
				Code:    "404",
				Message: "User not found",
				Data:    err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Model{
				Code:    "500",
				Message: "Query error",
				Data:    err.Error(),
			})
		}
	}

	var requestBody map[string]interface{}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, response.Model{
			Code:    "400",
			Message: "Invalid request body",
			Data:    err.Error(),
		})
	}

	updates := make(map[string]interface{})
	for key, value := range requestBody {
		updates[key] = value
	}

	if err := db.Model(&models.User{}).Where("user_id = ?", ui).Updates(updates).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.Model{
			Code:    "500",
			Message: "Error updating",
			Data:    err.Error(),
		})
	}

	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, response.Model{
			Code:    "500",
			Message: "Query error",
			Data:    err.Error(),
		})
	}

	r := response.Model{
		Code:    "202",
		Message: "Updated successfully",
		Data:    user,
	}
	return c.JSON(http.StatusAccepted, r)
}

func GetUser(c echo.Context) error {
	ui := c.QueryParam("user_id")
	db := configuration.GetDB()

	var user models.User
	if err := db.Where("user_id = ?", ui).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, response.Model{
				Code:    "404",
				Message: "User not found",
				Data:    err.Error(),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, response.Model{
				Code:    "500",
				Message: "Query error",
				Data:    err.Error(),
			})
		}
	}

	r := response.Model{
		Code:    "200",
		Message: "Correctly consulted",
		Data:    user,
	}
	return c.JSON(http.StatusOK, r)
}
