package handler

import (
	"clallenges/challenges/constant"
	"clallenges/challenges/model"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
	"log"
	"sync"
	"time"
)

func ProcessToCharge(client *omise.Client, jobs <-chan *model.Payment, results chan<- *omise.Charge, result *SummaryResult,wg *sync.WaitGroup) {
	for j := range jobs {
		result.TotalPay++
		if j.Card.ExpirationYear < time.Now().Year() || (j.Card.ExpirationYear == time.Now().Year() && j.Card.ExpirationMonth < time.Now().Month()) {
			result.FailPay += j.Amount
		} else {
			result.FailPay += j.Amount
			token, createToken := &omise.Token{}, &operations.CreateToken{
				Name:            j.Card.Name,
				Number:          j.Card.Number,
				ExpirationMonth: j.Card.ExpirationMonth,
				ExpirationYear:  j.Card.ExpirationYear,
				SecurityCode:    j.Card.SecurityCode,
			}
			if err := client.Do(token, createToken); err != nil {
				result.FailPay += j.Amount
				log.Printf("error from create Token %++v\n", err)

			} else {
				charge, createCharge := &omise.Charge{}, &operations.CreateCharge{
					Amount:   j.Amount,
					Currency: constant.DefaultCurrency,
					Card:     token.ID,
				}
				if err := client.Do(charge, createCharge); err != nil {
				//	log.Printf("err from create charge %++v\n", err)
				} else {
					result.SortTopPayments(j)
					result.SuccessPay += j.Amount
					results <- charge
				}
			}
		}
	}
	wg.Done()
}
