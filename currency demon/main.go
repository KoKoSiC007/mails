package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"example.com/m/v2/internal/repositories"
	usecases "example.com/m/v2/internal/use_cases"
	"github.com/go-co-op/gocron/v2"
)

type Config struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Enable   bool   `json:"enable"`
	Schedule string `json:"schedule"`
}

func main() {
	conf := GetConf()

	repository, err := repositories.NewRepo()
	if err != nil {
		log.Fatal(err)
	}
	parser := usecases.NewInteractor(repository)

	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range *conf {
		j, err := s.NewJob(
			gocron.CronJob(val.Schedule, true),
			gocron.NewTask(parser.Call, []string{val.Name}),
		)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(j.ID())
	}

	s.Start()

	select {}
}

func GetConf() *[]Config {
	urlP, err := url.Parse("http://172.16.0.2:8080")
	if err != nil {
		log.Fatal(err)
	}
	urlP.Path += "/api/v1/currencies"

	resp, err := http.Get(urlP.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var data []Config
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Fatal(err)
	}

	return &data
}
