package usecases

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/m/v2/internal/entities"
	"example.com/m/v2/internal/repositories"
)

type Interactor struct {
	repositoty repositories.Repo
}

func NewInteractor(repo *repositories.Repo) *Interactor {
	return &Interactor{repositoty: *repo}
}

func (i *Interactor) Call(currencies []string) error {
	data, err := i.request()
	if err != nil {
		return err
	}

	rates := i.unmarshal(data)
	// for _, rate := range rates {
	// 	fmt.Printf("%v - %f\n", rate.Name, rate.Rate)
	// }

	filteredRates := i.filter(rates, currencies)

	err = i.save(filteredRates)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interactor) request() ([]string, error) {
	date := time.Now().Format("02.01.2006")
	url := fmt.Sprintf("https://www.cnb.cz/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/daily.txt?date=%v", date)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 YaBrowser/25.2.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return strings.Split(string(body), "\n")[2:], nil
}

func (i *Interactor) unmarshal(marshaledRates []string) []entities.Rate {
	rates := make([]entities.Rate, len(marshaledRates)-1)
	for i, s := range marshaledRates {
		var rate entities.Rate
		parsedStr := strings.Split(s, "|")
		if len(parsedStr) == 1 {
			continue
		}
		rate.Name = parsedStr[3]
		parsedRate, err := strconv.ParseFloat(parsedStr[4], 32)
		if err != nil {
			fmt.Println("Can't parse rate %v, error %v", parsedStr[4], err)
			continue
		}
		rate.Rate = float32(parsedRate)
		rates[i] = rate
	}

	return rates
}

func (i *Interactor) filter(rates []entities.Rate, filter []string) []entities.Rate {
	var filteredRates = make(map[string]*entities.Rate, len(filter)-1)
	for _, currency := range filter {
		filteredRates[currency] = nil
	}

	for _, rate := range rates {
		_, present := filteredRates[rate.Name]
		if present {
			filteredRates[rate.Name] = &rate
		}
	}
	var result []entities.Rate
	for _, v := range filteredRates {
		if v != nil {
			result = append(result, *v)
		}
	}
	return result
}

func (i *Interactor) save(rates []entities.Rate) error {
	for _, rate := range rates {
		_, err := i.repositoty.Create(rate)
		if err != nil {
			fmt.Println("Can't save rate %v, error: %v", rate.Name, err)

			continue
		}
	}

	return nil
}
