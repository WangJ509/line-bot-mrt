package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/WangJ509/line-bot-mrt/mrt"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func main() {
	port := os.Getenv("PORT")

	bot, err := linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	if err != nil {
		log.Fatal("failed to new line bot", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(RequestLogger())

	router.GET("health-check", func(ctx *gin.Context) {
		ctx.String(200, "I am healthy!!!")
	})

	router.POST("/callback", getCallbackHandler(bot))
	router.Run(":" + port)
}

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		log.Println(string(body))
		log.Println(c.Request.Header)
		c.Next()
	}
}

func getCallbackHandler(bot *linebot.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		events, err := bot.ParseRequest(ctx.Request)

		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Println("Invalid signature")
				ctx.AbortWithStatus(400)
			} else {
				ctx.AbortWithStatus(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					replyString := handleTextMessage(message.Text)
					log.Println(replyString)

					replyMessage := linebot.NewTextMessage(replyString)
					if _, err = bot.ReplyMessage(event.ReplyToken, replyMessage).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

var (
	timeTableFormatMessage = "正確格式： 「時刻表 出發站 終點站 數量」, 範例： 「時刻表 景美 松山 5」"
)

func handleTextMessage(text string) string {
	mrtService := mrt.NewMRTService()

	inputs := strings.Split(text, " ")
	switch inputs[0] {
	case "時刻表":
		if len(inputs) < 4 {
			return timeTableFormatMessage
		}
		station := inputs[1]
		destination := inputs[2]
		number, err := strconv.Atoi(inputs[3])
		if err != nil {
			return timeTableFormatMessage
		}

		result, err := mrtService.GetUpcomingTimeTable(station, destination, number)
		if err != nil {
			return "Faces error: " + err.Error()
		}

		if len(result) < number {
			result = append(result, "末班車已駛離")
		}

		return strings.Join(result, "\n")
	}

	return ""
}
