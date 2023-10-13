package responses

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	reddb "simple_bot/redis"
	"time"

	"github.com/go-redis/redis"
	tele "gopkg.in/telebot.v3"
)

type Response struct {
	Data struct {
		Temp      float64 `json:"temp_c"`
		Wind      float64 `json:"wind_kph"`
		Wind_dir  string  `json:"wind_dir"`
		Humidity  int     `json:"humidity"`
		Cloud     int     `json:"cloud"`
		Feelslike float64 `json:"feelslike_c"`
	} `json:"current"`
}

func WeatherMain(bot tele.Context) error {
	return bot.Send("Пришли название города в котором ты хочешь узнать погоду")
}

func WeatherMessage(city string, bot tele.Context) error {
	wind_dir := map[string]string{
		"N":   "северное",
		"NNE": "севере-северо-восточное",
		"NE":  "северо-восточное",
		"ENE": "восточно-северо-восточое",
		"E":   "восточное",
		"ESE": "восточно-юго-восточное",
		"SE":  "юго-восточное",
		"SSE": "юго-юго-восточное",
		"S":   "южное",
		"SSW": "юго-юго-западное",
		"SW":  "юго-западное",
		"WSW": "западно-юго-западное",
		"W":   "западное",
		"WNW": "западно-северо-западное",
		"NW":  "северо-западное",
		"NNW": "северо-северо-западный"}

	rdb := reddb.GetRedis()

	val, err := rdb.Get(city).Result()
	if err == redis.Nil {
		fmt.Printf("%s does not exist", city)
	} else if err != nil {
		panic(err)
	} else if val != "0" {
		return bot.Send(val)
	}

	weather_data := Weather(city)
	message_text := "В городе %s сейчас %s, \n температура составляет %.1f градусов, \n Ощущается как: %.1f. \n Скорость ветра составляет %.1f км\\ч \n Направление %s \n Влажность составляет %d процентов \n"
	response := fmt.Sprintf(
		message_text,
		city,
		WeatherCondition(weather_data.Data.Cloud),
		weather_data.Data.Temp,
		weather_data.Data.Feelslike,
		weather_data.Data.Wind,
		wind_dir[weather_data.Data.Wind_dir],
		weather_data.Data.Humidity)
	rdb.Set(city, response, 10*time.Minute).Err()
	return bot.Send(response)
}

func Weather(city string) Response {
	api_key := "02d57a4ae70f484facf124135231210"
	link := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", api_key, city)
	request, err := http.Get(link)
	if err != nil {
		log.Fatal()
	}
	defer request.Body.Close()

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Fatal()
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatal()
	}

	return result
}

func WeatherCondition(cloud int) string {
	if cloud < 25 {
		return "ясно"
	} else if cloud < 50 {
		return "переменная облачность"
	} else if cloud < 75 {
		return "облачно с прояснениями"
	} else {
		return "пасмурно"
	}

}
