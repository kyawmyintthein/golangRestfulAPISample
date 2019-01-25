package controller

import (
	"github.com/go-chi/chi"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/service"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"net/http"
)

type UserController struct{
	BaseController
	UserService service.UserServiceInterface
}

// Register godoc
// @Summary Create new user
// @Description create new user
// @Tags  user
// @Accept  json
// @Produce  json
// @Param payload body model.CreateUserRequest true "Request Payload"
// @Success 201 {object} model.User
// @Failure 500 {object} model.ErrorResponse
// @Router /users [post]
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req model.CreateUserRequest
	err := c.decodeAndValidate(r, &req)
	if err != nil{
		c.WriteError(r, w, err)
		return
	}

	user, err := c.UserService.Create(r.Context(), &req)
	if err != nil{
	}

	c.WriteJSON(r, w, http.StatusCreated, user)
}


// Update godoc
// @Summary Update user profile
// @Description update user's profile with ID
// @Tags  user
// @Accept  json
// @Produce  json
// @Param payload body model.UpdateUserRequest true "Request Payload"
// @Success 200 {object} model.User
// @Failure 500 {object} model.ErrorResponse
// @Router /users [put]
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var req model.UpdateUserRequest
	err := c.decodeAndValidate(r, &req)
	if err != nil{
		c.WriteError(r, w, err)
		return
	}

	user, err := c.UserService.Update(r.Context(), &req)
	if err != nil{
		c.WriteError(r, w, err)
		return
	}

	c.WriteJSON(r, w, http.StatusOK, user)
}


// GetProfile godoc
// @Summary Get user profile
// @Description retrieve user's profile with ID
// @Tags  user
// @Accept  json
// @Produce  json
// @Success 200 {object} model.User
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{user_id} [get]
func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == ""{
		err := errors.New(ecodes.InvalidRequestParameters, constant.InvalidRequestParameter)
		c.WriteError(r, w, err)
		return
	}

	user, err := c.UserService.FindByID(r.Context(), userID)
	if err != nil{
		c.WriteError(r, w, err)
		return
	}

	c.WriteJSON(r, w, http.StatusOK, user)
}


// Remove godoc
// @Summary Remove user profile
// @Description remove user's profile with ID
// @Tags  user
// @Accept  json
// @Produce  json
// @Success 204 {object} model.User
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{user_id} [delete]
func (c *UserController) Remove(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if userID == ""{
		err := errors.New(ecodes.InvalidRequestParameters, constant.InvalidRequestParameter)
		c.WriteError(r, w, err)
		return
	}

	err := c.UserService.DeleteByID(r.Context(), userID)
	if err != nil{
		c.WriteError(r, w, err)
		return
	}

	c.WriteWithStatus(w, http.StatusNoContent)
}


// GetAllUsers godoc
// @Summary get all user profile
// @Description retrieve all user's profile
// @Tags  user
// @Accept  json
// @Produce  json
// @Success 200 {array} model.User
// @Failure 500 {object} model.ErrorResponse
// @Router /users [get]
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.UserService.FindAll(r.Context())
	if err != nil {
		c.WriteError(r, w, err)
		return
	}
	c.WriteJSON(r, w, http.StatusOK, users)
}