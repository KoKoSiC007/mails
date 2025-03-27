package parser

import (
	"fmt"
	"io"
	"net/http"
	"odyssey/m/v2/repositories"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	Url      string
	rateRepo *repositories.RateRepo
}

func NewParser(repo *repositories.RateRepo) *Parser {
	return &Parser{
		Url:      "https://www.cnb.cz/en/financial_markets/foreign_exchange_market/exchange_rate_fixing/year.txt?",
		rateRepo: repo,
	}
}

func (p *Parser) Parse(start, end time.Time, currencies *[]string) error {
	startYear := start.Year()
	endYear := end.Year()
	count := 0
	diffYear := endYear - startYear
	if diffYear == 0 {
		p.sync(strconv.Itoa(startYear), currencies)
	} else {
		for count <= diffYear {
			fmtYear := strconv.Itoa(start.AddDate(count, 0, 0).Year())
			p.sync(fmtYear, currencies)
			count += 1
		}
	}
	return nil
}

func (p *Parser) sync(year string, currencies *[]string) error {
	data, err := p.request(year)
	if err != nil {
		return err
	}

	records, err := p.collectData(&data, currencies)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = p.save(records)
	if err != nil {
		return err
	}

	return nil
}

func (p *Parser) request(year string) ([]string, error) {
	url := fmt.Sprintf("%vyear=%v", p.Url, year)
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

	return strings.Split(string(body), "\n"), nil
}

func (p *Parser) collectData(data *[]string, currencies *[]string) (*[]repositories.DbRate, error) {
	header := strings.Split((*data)[0], "|")
	colSeq := make(map[int]string)
	for i, value := range header[1:] {
		hCurName := strings.Split(value, " ")[1]
		if slices.Contains(*currencies, hCurName) {
			colSeq[i+1] = hCurName
		}
	}
	rawData := (*data)[1:]

	var result []repositories.DbRate
	for _, value := range rawData {
		row := strings.Split(value, "|")
		if len(row) == 1 {
			continue
		}
		date, err := time.Parse("02.01.2006", row[0])
		if err != nil {
			return nil, err
		}
		for i, curName := range colSeq {
			rate, err := strconv.ParseFloat(row[i], 32)
			if err != nil {
				return nil, err
			}

			dbRate := repositories.DbRate{
				Name:      curName,
				Rate:      float32(rate),
				CreatedAt: date,
			}

			result = append(result, dbRate)
		}
	}
	return &result, nil
}

func (p *Parser) save(records *[]repositories.DbRate) error {
	for _, record := range *records {
		err := p.rateRepo.Create(record)
		if err != nil {
			return err
		}
	}

	return nil
}
