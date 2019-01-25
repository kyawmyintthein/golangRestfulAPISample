package controller

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"net/http"
)

type HttpErrorController struct{
	BaseController
}

func (c *HttpErrorController) ResourceNotFound(w http.ResponseWriter, r *http.Request) {
	err := errors.New(ecodes.NotFound, constant.ResourceNotFound)
	c.WriteError(r, w, err)
	return
}