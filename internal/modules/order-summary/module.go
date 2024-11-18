package ordersummary

import (
	bitkubapi "bitkub-port-avg/internal/connectors/bitkub-api"
	"fmt"
	"os"
	"strconv"
)

type IOrderSummary interface {
	GetOrderSummary(tokenSymbol string, startTimestamp *uint64) (rateToFiatAmount map[float64]float64, rateToTokenAmount map[float64]float64, err error)
}

type orderSummary struct {
	bitkubApiClient bitkubapi.IBitkubApiClient
}

func NewOrderSummary(bitkubApiClient bitkubapi.IBitkubApiClient) IOrderSummary {
	return &orderSummary{
		bitkubApiClient: bitkubApiClient,
	}
}

func (o *orderSummary) GetOrderSummary(tokenSymbol string, startTimestamp *uint64) (map[float64]float64, map[float64]float64, error) {
	orderHistories, err := o.bitkubApiClient.RequestOrderHistories(tokenSymbol, startTimestamp)
	if err != nil {
		return nil, nil, err
	}

	rateToFiatAmount := map[float64]float64{}
	rateToTokenAmount := map[float64]float64{}
	tokenRemainingAmount := 0.0
	totalInvestmentFiat := 0.0
	for _, order := range orderHistories {
		amountFloat, _ := strconv.ParseFloat(order.Amount, 64)
		feeFloat, _ := strconv.ParseFloat(order.Fee, 64)
		creditFloat, _ := strconv.ParseFloat(order.Credit, 64)
		rateFloat, _ := strconv.ParseFloat(order.Rate, 64)
		if order.Side == "buy" {
			fiatAmount := amountFloat - feeFloat + creditFloat
			tokenAmount := fiatAmount / rateFloat

			tokenRemainingAmount += tokenAmount
			totalInvestmentFiat += amountFloat

			rateToFiatAmount[rateFloat] += fiatAmount
			rateToTokenAmount[rateFloat] += tokenAmount
		} else {
			fiatAmount := amountFloat*rateFloat - feeFloat + creditFloat

			totalInvestmentFiat -= fiatAmount
			tokenRemainingAmount -= amountFloat

			rateToFiatAmount[rateFloat] -= fiatAmount
			rateToTokenAmount[rateFloat] -= amountFloat
		}
	}

	fmt.Println("====================================================")
	fmt.Printf("Token %s remaining amount: %f\n", tokenSymbol, tokenRemainingAmount)
	fmt.Printf("Total investment in fiat: %f\n", totalInvestmentFiat)
	if tokenRemainingAmount > 0 {
		fmt.Println("Average price per token: ", totalInvestmentFiat/tokenRemainingAmount)
	} else {
		fmt.Println("Average price per token: 0")
	}

	if len(rateToFiatAmount) == 0 {
		if err := writeToFile(tokenSymbol, rateToFiatAmount, rateToTokenAmount); err != nil {
			return nil, nil, err
		}
	}

	return rateToFiatAmount, rateToTokenAmount, nil
}

func writeToFile(tokenSymbol string, rateToFiatAmount map[float64]float64, rateToTokenAmount map[float64]float64) error {
	keys := make([]float64, 0, len(rateToFiatAmount))
	for k := range rateToFiatAmount {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	fileContent := "rate,fiat_amount,token_amount\n"
	for _, key := range keys {
		fileContent += fmt.Sprintf("%f,%f,%f\n", key, rateToFiatAmount[key], rateToTokenAmount[key])
	}

	// Ensure the directory exists
	if err := os.MkdirAll("./reports", os.ModePerm); err != nil {
		return err
	}

	// write to file
	file, err := os.Create(fmt.Sprintf("./reports/%s_order_summary.csv", tokenSymbol))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fileContent)
	if err != nil {
		return err
	}

	return nil
}
