// Сервер GoNews.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"Agrigator/pkg/api"
	"Agrigator/pkg/rss"
	storage "Agrigator/pkg/storage/pstg"
)

type config struct {
	URLS   []string `json:"rss"`
	Period int      `json:"cicle"`
	Port   string   `json:"port"`
	DBadr  string   `json:"dbadr"`
}

func main() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	b, err := os.ReadFile("./cmd/config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.Println("Запуск службы...")
	db, err := storage.New(config.DBadr)
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	chPosts := make(chan []storage.Article)
	chErrs := make(chan error)
	for _, url := range config.URLS {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	go func() {
		for posts := range chPosts {
			db.SaveArticles(posts)
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()
	log.Println("Запуск сервера. Порт", config.Port)
	err = http.ListenAndServe(config.Port, api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// parseURL выполняет асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, db *storage.DB, posts chan<- []storage.Article, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
