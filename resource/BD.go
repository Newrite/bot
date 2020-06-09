package resource

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"sync"
)

var db *sql.DB
var once sync.Once

func SingleDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("sqlite3", "bot.db")
		if err != nil || db == nil {
			log.WithFields(log.Fields{
				"package":  "resource",
				"function": "sql.Open",
				"error":    err,
			}).Fatalln("Не удалось открыть ДБ.")
		}
	})
	return db
}

func AddQuoteDB(message string) {
	_, err := SingleDB().Exec(`
CREATE TABLE IF NOT EXISTS QUOTES
(
	ID INTEGER NOT NULL PRIMARY KEY autoincrement,
	CONTAIN TEXT NOT NULL
);`)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка создания таблицы в бд.")
	}
	message = strings.Replace(message, `'`, `''`, -1)
	message = strings.Replace(message, `&`, `&&`, -1)
	_, err = SingleDB().Exec(`
		INSERT INTO QUOTES
			(CONTAIN)
		VALUES
			($1)`, message)
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "Exec",
			"error":    err,
		}).Errorln("Ошибка добавления значения в бд.")
	}
}

func DBQuote() string {
	quotes := make([]string, 0)
	var lastID int
	rows, err := SingleDB().Query(`SELECT * FROM QUOTES`)
	if rows == nil || err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "Query",
			"error":    err,
		}).Errorln("Ошибка запроса с бд.")
		return err.Error()
	}
	for rows.Next() {
		var tmp string
		err = rows.Scan(&lastID, &tmp)
		if err != nil {
			log.WithFields(log.Fields{
				"package":  "main",
				"function": "Scan",
				"error":    err,
			}).Errorln("Ошибка скан запроса.")
			return err.Error()
		}
		quotes = append(quotes, tmp)
	}
	err = rows.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"package":  "main",
			"function": "Close",
			"error":    err,
		}).Errorln("Ошибка закрытия rows.")
		return err.Error()
	}
	return quotes[rand.Intn(lastID)]
}
