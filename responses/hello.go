package responses

import (
	"database/sql"
	"log"

	tele "gopkg.in/telebot.v3"
)

func Hello(bot tele.Context) error {
	db, err := sql.Open("postgres", "postgresql://postgres:123@localhost:5432/Go?sslmode=disable")
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INT PRIMARY KEY,
		username TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Exec(`INSERT INTO users VALUES ($1, $2) ON CONFLICT (id) do nothing`, bot.Message().Chat.ID, bot.Message().Chat.Username)
	if err != nil {
		log.Fatal()
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal()
	}
	if rowsAffected == 0 {
		return bot.Send("Hello again!")
	}
	return bot.Send("Hello!")
}
