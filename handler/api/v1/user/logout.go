package user

import (
	"github.com/1024casts/1024casts/model"
	"github.com/1024casts/1024casts/pkg/app"

	"github.com/gin-gonic/gin"
)

// @Summary Login generates the authentication token
// @Produce  json
// @Param username body string true "Username"
// @Param password body string true "Password"
// @Success 200 {string} json "{"code":0,"message":"OK","data":{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MjgwMTY5MjIsImlkIjowLCJuYmYiOjE1MjgwMTY5MjIsInVzZXJuYW1lIjoiYWRtaW4ifQ.LjxrK9DuAwAzUD8-9v43NzWBN7HXsSLfebw92DKd1JQ"}}"
// @Router /logout [post]
func Logout(c *gin.Context) {

	app.Response(c, nil, model.Token{Token: ""})
}
