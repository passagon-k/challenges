package handler

import (
	"clallenges/challenges/model"
	"fmt"
	"strings"
)

type SummaryResult struct {
	SuccessPay int64
	FailPay    int64
	Currency   string
	TotalPay   int64
	TopDonate  []*model.Payment
}

func (donation *SummaryResult) SortTopPayments(payment *model.Payment) {
	len := len(donation.TopDonate)
	for i, p := range donation.TopDonate {
		if payment.Amount > p.Amount {
			donation.TopDonate = append(donation.TopDonate[:i], append([]*model.Payment{payment}, donation.TopDonate[i:len-1]...)...)
			break
		}
	}
}

func (donation SummaryResult) ShowSummary() {

	//donation.Currency = strings.ToUpper(donation.Currency)

	fmt.Printf("\n%30s %s %17.2f\n", "total received:", strings.ToUpper(donation.Currency), float64((donation.SuccessPay+donation.FailPay)/100))
	fmt.Printf("%30s %s %17.2f\n", "successfully donated:", strings.ToUpper(donation.Currency), float64(donation.SuccessPay/100))
	fmt.Printf("%30s %s %17.2f\n\n", "faulty donation:", strings.ToUpper(donation.Currency), float64(donation.FailPay/100))
	fmt.Printf("%30s %s %17.2f\n", "average per person:", strings.ToUpper(donation.Currency), (float64(donation.SuccessPay+donation.FailPay)/float64(donation.TotalPay))/100)
	fmt.Printf("%30s", "top donors:")
	for i, payment := range donation.TopDonate {
		if payment == nil || payment.Card == nil {
			// reserving line for future top donors, clean-up linces will be mess
			fmt.Println()
			continue
		}
		if i == 0 {
			fmt.Printf(" %s\n", payment.Card.Name)
		} else {
			fmt.Printf("%30s %s\n", "", payment.Card.Name)
		}
	}
}
