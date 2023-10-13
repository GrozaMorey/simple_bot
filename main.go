package main

import (
	"fmt"
	"log"
	reddb "simple_bot/redis"
	"simple_bot/responses"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/text/language"

	"github.com/bregydoc/gtranslate"
	tele "gopkg.in/telebot.v3"
)

func main() {
	rdb := reddb.GetRedis()
	responsFunc := map[string]func(tele.Context) error{
		"/hello": responses.Hello,
	}

	pref := tele.Settings{
		Token:  "token",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	//b.Use(middleware.Logger())
	//b.Use(middleware.AutoRespond())

	b.Handle(tele.OnText, func(c tele.Context) error {
		userId := string(rune(c.Message().Chat.ID))

		state, err := rdb.Get(userId).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(state, err)
		if state == "1" {
			rdb.Set(userId, false, 0).Err()
			if err != nil {
				panic(err)
			}

			city, err := gtranslate.Translate(c.Message().Text, language.Russian, language.English)
			if err != nil {
				fmt.Printf("Cant translate %s", c.Message().Text)
			}

			return responses.WeatherMessage(city, c)
		}

		if c.Message().Text == "/weather" {
			rdb.Set(userId, true, 0).Err()
			return responses.WeatherMain(c)
		}

		for n, f := range responsFunc {
			if c.Message().Text == n {
				return f(c)
			}
		}
		log.Println("unknown message text:", c.Message().Text)
		return c.Send("wtf")
	})

	b.Start()
}
