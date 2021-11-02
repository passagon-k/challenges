package main

import (
	"clallenges/challenges/constant"
	"clallenges/challenges/handler"
	"clallenges/challenges/model"
	"fmt"
	"github.com/omise/omise-go"
	"io"
	"omise-go-challenge/cipher"
	"os"
	"sync"
	"time"
)

var (
	options *model.Options
)

func main() {
	fmt.Printf("performing donations...\n")
	jobs := make(chan *model.Payment, constant.NumberOfMaxJobs)
	results := make(chan *omise.Charge, constant.NumberOfMaxJobs)
	client, err := omise.NewClient(constant.OmisePublicKey, constant.OmiseSecretKey)
	start := time.Now()
	fileOpen, err := os.Open(constant.FilePath)
	if err != nil {
		panic(err)
	}
	defer fileOpen.Close()

	var reader io.Reader
	reader = fileOpen

	decoder, err := cipher.NewRot128Reader(reader)
	if err != nil {
		panic(err)
	}
	handler.ReadFileCsv(decoder, jobs)

	result := &handler.SummaryResult{TopDonate: []*model.Payment{{}, {}, {}}, Currency: constant.DefaultCurrency}
	wg := new(sync.WaitGroup)
	wg.Add(constant.NumberOfWorker)
	checkWait := 0
	for w := 1; w <= constant.NumberOfWorker; w++ {
		if checkWait == constant.NumberOfCheckWait {
			time.Sleep(time.Second * 2)
		}
		go handler.ProcessToCharge(client, jobs, results, result, wg)
		checkWait++
	}
	wg.Wait()
	fmt.Printf("done.\n")
	result.ShowSummary()
	elapsed := time.Since(start)
	fmt.Printf("Excution Time %s", elapsed)

}
