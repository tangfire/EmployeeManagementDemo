package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func GetCurrentUserID(c *gin.Context) (uint, error) {
	rawUserID, ok := c.Get("userID")
	if !ok {
		return 0, errors.New("userID not found")
	}
	userID, ok := rawUserID.(uint)
	if !ok {
		return 0, errors.New("invalid userID type")
	}
	return userID, nil
}

func GetCurrentUserRole(c *gin.Context) (string, error) {
	rawRole, ok := c.Get("userRole")
	if !ok {
		return "", errors.New("userRole not found")
	}
	role, ok := rawRole.(string)
	if !ok {
		return "", errors.New("invalid userRole type")
	}
	return role, nil
}
