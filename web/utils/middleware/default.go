package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"ihome/web/logic/user"
	"net/http"
)

func LoginCheck(c *gin.Context) {
	s := sessions.Default(c)
	userName := s.Get(user.SESSION_USER_NAME)
	if userName == nil {
		c.Abort()
		c.Redirect(http.StatusUnauthorized, "/home/login.html")
	}
}
