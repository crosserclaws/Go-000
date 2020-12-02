package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/crosserclaws/Go-000/Week02/biz"
	customError "github.com/crosserclaws/Go-000/Week02/custom/error"
	"github.com/gin-gonic/gin"
)

type userID struct {
	ID int `uri:"id" binding:"required"`
}

// GetUserInfoByID is an endpoint return an user's information by a given user ID. Return 404 if an invalid user ID is given.
func GetUserInfoByID(c *gin.Context) {
	var uid userID
	if err := c.ShouldBindUri(&uid); err != nil {
		log.Printf("Failed to get user id=%d: %v", uid.ID, err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	user, err := biz.FetchUserInfo(uid.ID)
	if err != nil {
		log.Printf("Failed to get user id=%d: %+v", uid.ID, err)
		if errors.Is(err, customError.EmptyResultError) {
			c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"ID": user.ID, "name": user.Name})
}

// GetUsers is an endpoint return an array of users. Return an empty array when no users.
func GetUsers(c *gin.Context) {
	users, err := biz.ListUsers()
	if err != nil {
		log.Printf("Failed to list users: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
