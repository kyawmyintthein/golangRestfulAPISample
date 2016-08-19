package controllers

import (
	"echo_rest_api/models"
	"github.com/labstack/echo"
	"net/http"
	"strconv"
)

//create user
func CreateUser(c echo.Context) error {
	user := new(models.User)
	var err error
	// marchal json to object
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	err = models.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusCreated, user)
}

//update user
func UpdateUser(c echo.Context) error {
	// Parse the content
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusForbidden, err)
	}

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"),10,64)
	m, err := models.GetUserById(id)
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}

	// update user data
	err = m.UpdateUser(user)
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}

	return c.JSON(http.StatusOK, m)
}

//delete user
func DeleteUser(c echo.Context) error {
	var err error

	// get the param id
	id, _ := strconv.ParseUint(c.Param("id"),10,64)
	m, err := models.GetUserById(id)
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}

	err = m.DeleteUser()
	return c.JSON(http.StatusNoContent, err)
}

// get one user
func ShowUser(c echo.Context) error {
	var (
		user models.User
		err error
	)
	id, _ := strconv.ParseUint(c.Param("id"),10,64)
	user, err = models.GetUserById(id)
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, user)
}

// get all users
func AllUsers(c echo.Context) error {
	var (
		users []models.User
		err error
	)
	users, err = models.GetUsers()
	if err != nil{
		return c.JSON(http.StatusForbidden, err)
	}
	return c.JSON(http.StatusOK, users)
}
