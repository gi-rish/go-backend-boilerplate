package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var ErrorCodes = map[string]string{
	"PLS0001": "no translations found",
	"PLS0002": "please provide appname header value",
	"PLS0003": "appname header value is not valid",
	"PLS0004": "wrong request format",
	"PLS0005": "invalid input",
	"PLS0006": "system error",
}

func ReadRequest(req *interface{}, c *gin.Context) *interface{} {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return req
}

// generate json result for error responses.
func ErrorResult(c *gin.Context, err error, msg, code string, errCode int) {
	c.JSON(errCode, gin.H{
		"result":     "fail",
		"error":      err.Error(),
		"error_code": code,
	})
}

// generate json result for error responses for V2.
func GenerateErrorResult(c *gin.Context, msg, code string, errCode int) {
	c.JSON(errCode, gin.H{
		"message": msg,
		"code":    code,
		"data": gin.H{
			"error": msg,
		},
	})
}

// generate json result for success responses.
func SuccessResult(c *gin.Context, res map[string]interface{}) {
	c.JSON(http.StatusOK, res)
}

// validate appname header from request
func ValidateHeader(req http.Request) (string, error) {
	appname := req.Header.Get("appname")
	if len(strings.TrimSpace(appname)) == 0 {
		err := errors.New(ErrorCodes["PLS0002"])
		return "PLS0002", err
	}
	isAlphabet := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	if !isAlphabet(appname) {
		err := errors.New(ErrorCodes["PLS0003"])
		return "PLS0003", err
	}
	return "", nil
}
