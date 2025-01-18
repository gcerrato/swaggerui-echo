package main

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

type ResponseType string

const (
	Error ResponseType = "error message"
	Info  ResponseType = "info message"
)

type apiResponse struct {
	Code     int          `json:"code"`
	RespType ResponseType `json:"type"`
	Message  string       `json:"message"`
}

type Pet struct {
	Name    string `json:"name"`
	PetType string `json:"pet_type"`
}

var petDB = []Pet{
	{
		Name:    "Bobby",
		PetType: "cat",
	},
	{
		Name:    "Ralph",
		PetType: "dog",
	},
	{
		Name:    "Shirley",
		PetType: "Armadillo",
	},
}
var dbLock sync.RWMutex

func petHandler(c echo.Context) error {
	type petResponse struct {
		ID int `json:"id"`
		Pet
	}

	idStr := strings.TrimPrefix(c.Request().URL.Path, "/pet/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, apiResponse{
			Code:     http.StatusBadRequest,
			RespType: Error,
			Message:  "invalid id",
		})
	}

	switch c.Request().Method {
	case "PUT":
		if c.Request().Header.Get("X-API-KEY") == "" {
			return c.JSON(http.StatusUnauthorized, apiResponse{
				Code:     http.StatusUnauthorized,
				RespType: Error,
				Message:  "unauthorized",
			})
		}
		dbLock.Lock()
		defer dbLock.Unlock()
		if id < 0 || id >= len(petDB) {
			return c.JSON(http.StatusNotFound, apiResponse{
				Code:     http.StatusNotFound,
				RespType: Error,
				Message:  "pet not found",
			})
		}
		newPet := &Pet{}
		if err := c.Bind(newPet); err != nil {
			return c.JSON(http.StatusBadRequest, apiResponse{
				Code:     http.StatusBadRequest,
				RespType: Error,
				Message:  "invalid pet object",
			})
		}
		if newPet.Name != "" {
			petDB[id].Name = newPet.Name
		}
		if newPet.PetType != "" {
			petDB[id].PetType = newPet.PetType
		}
		return c.JSON(http.StatusOK, apiResponse{
			Code:     http.StatusOK,
			RespType: Info,
			Message:  "pet updated",
		})

	case "GET":
		dbLock.RLock()
		defer dbLock.RUnlock()
		if id < 0 || id >= len(petDB) {
			return c.JSON(http.StatusNotFound, apiResponse{
				Code:     http.StatusNotFound,
				RespType: Error,
				Message:  "pet not found",
			})
		}
		return c.JSON(http.StatusOK, petResponse{ID: id, Pet: petDB[id]})

	default:
		return c.JSON(http.StatusMethodNotAllowed, apiResponse{
			Code:     http.StatusMethodNotAllowed,
			RespType: Error,
			Message:  "method not allowed",
		})
	}
}
