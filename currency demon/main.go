package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"example.com/m/v2/internal/repositories"
	usecases "example.com/m/v2/internal/use_cases"
	"github.com/go-co-op/gocron/v2"
	"gopkg.in/yaml.v3"
)

type conf struct {
	Currencies []string `yaml:currencies`
}

func main() {
	var conf conf
	file, err := os.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		log.Fatal(err)
	}
	repository, err := repositories.NewRepo()
	if err != nil {
		log.Fatal(err)
	}
	parser := usecases.NewInteractor(repository)

	parser.Call(conf.Currencies)
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	j, err := s.NewJob(
		gocron.DurationJob(
			12*time.Hour,
		),
		gocron.NewTask(parser.Call, conf.Currencies),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(j.ID())

	s.Start()

	select {}
}
