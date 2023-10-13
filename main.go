package main

import (
	"fmt"
	"log"
	"simple_bot/responses"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"

	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	responsFunc := map[string]func(tele.Context) error{
		"/hello": responses.Hello,
	}

	pref := tele.Settings{
		Token:  "6216320540:AAHv5cUkrYWrLU4h5LJ7N-kUKx7LPLAuC54",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Use(middleware.Logger())
	b.Use(middleware.AutoRespond())

	b.Handle(tele.OnText, func(c tele.Context) error {
		userId := string(rune(c.Message().Chat.ID))

		state, err := rdb.Get(userId).Result()
		if err == redis.Nil {
			log.Println("State was no define")
		} else if err != nil {
			panic(err)
		}
		fmt.Println(state)
		if state != "" {
			rdb.Set(userId, false, 0).Err()
			if err != nil {
				panic(err)
			}
			return responses.WeatherMessage(c.Message().Text, c)
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
