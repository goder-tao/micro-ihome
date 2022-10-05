package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"net/http"
)

func main() {
	r := gin.Default()
	// session存储在服务端
	store, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("ihome"))
	if err != nil {
		log.Fatalln("redis store, ", err)
	}
	r.Use(sessions.Sessions("ihome", store))

	r.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("ihome_cookie")
		// cookie提供临时（关闭后）和期限两种
		if err != nil {
			cookie = "cookie_value"
			c.SetCookie("ihome_cookie", cookie, 60, "", "", false, false)
		}
		fmt.Println("cookie value: ", cookie)
		c.JSON(http.StatusOK, nil)
	})
	r.GET("/session", func(c *gin.Context) {
		s := sessions.Default(c)
		// s.Options()
		v := s.Get("test")
		if v == nil {
			v = "session_value"
			s.Set("test", v)
			err := s.Save()
			if err != nil {
				log.Errorln("save session,", err)
				c.JSON(http.StatusInternalServerError, nil)
				return
			}
		}
		fmt.Println("session value:", v)
		c.JSON(http.StatusOK, nil)
	})

	if err := r.Run(":9090"); err != nil {
		log.Fatalln("route run", err)
	}
}
