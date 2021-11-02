package handler

import (
	"clallenges/challenges/constant"
	"clallenges/challenges/model"
	"encoding/csv"
	"flag"
	"io"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var options *model.Options

func ReadFileCsv(decoder io.Reader, jobs chan<- *model.Payment) {
	commaFlag := flag.String("comma", ",", "delimiter for CSV files")
	options = &model.Options{
		Comma:   []rune(*commaFlag)[0],
		Headers: []string{"Name", "AmountSubunits", "CCNumber", "CVV", "ExpMonth", "ExpYear"},
	}
	options.Extension = filepath.Ext(constant.FilePath)
	csvReader := csv.NewReader(decoder)
	csvReader.FieldsPerRecord = len(options.Headers)
	csvReader.Comma = options.Comma
	csvReader.Read()
	csvLines, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		for _, line := range csvLines {
			amount, _ := strconv.ParseInt(line[1], 10, 64)
			month, err := strconv.ParseInt(line[4], 10, 0)
			if err != nil {
				panic(err)
			}
			year, err := strconv.ParseInt(line[5], 10, 0)
			if err != nil {
				panic(err)
			}
			payment := model.Payment{
				Amount: amount,
				Card: &model.Card{
					ExpirationMonth: time.Month(month),
					ExpirationYear:  int(year),
					Name:            line[0],
					Number:          line[2],
					SecurityCode:    line[3],
				},
			}
			jobs <- &payment

		}
		close(jobs)
	}()
	wg.Wait()
}
