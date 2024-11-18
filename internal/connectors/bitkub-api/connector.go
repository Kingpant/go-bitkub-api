package bitkubapi

import (
	"bitkub-port-avg/internal/types"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type IBitkubApiClient interface {
	RequestOrderHistories(tokenSymbol string, startTimestamp *uint64) ([]types.OrderHistory, error)
	RequestDepositHistories(tokenSymbol string) ([]types.DepositHistory, error)
	RequestWithdrawHistories(tokenSymbol string) ([]types.WithdrawHistory, error)
}

type bitkubApiClient struct {
	baseUrl   string
	apiKey    string
	apiSecret string
}

func NewBitkubApiClient(baseUrl, apiKey, apiSecret string) IBitkubApiClient {
	return &bitkubApiClient{
		baseUrl:   baseUrl,
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
}

func (b *bitkubApiClient) RequestOrderHistories(tokenSymbol string, startTimestamp *uint64) ([]types.OrderHistory, error) {
	tokenSymbol = strings.ToLower(tokenSymbol)
	page := uint64(1)
	limit := "100"
	orderHistories := []types.OrderHistory{}

	for {
		orderHistoryPath := "/api/v3/market/my-order-history"
		queryParma := fmt.Sprintf("?sym=%s_thb&p=%d&lmt=%s", tokenSymbol, page, limit)
		if startTimestamp != nil {
			queryParma += fmt.Sprintf("&start=%d", *startTimestamp)
		}

		body, err := b.httpRequest(orderHistoryPath, queryParma)
		if err != nil {
			return nil, err
		}

		orderHistoryResponse := types.OrderHistoryResponse{}
		err = json.Unmarshal(body, &orderHistoryResponse)
		if err != nil {
			return nil, err
		}

		orderHistories = append(orderHistories, orderHistoryResponse.Result...)

		if orderHistoryResponse.Pagination.Next == 0 {
			break
		}
		page = orderHistoryResponse.Pagination.Next
	}

	return orderHistories, nil
}

func (b *bitkubApiClient) RequestDepositHistories(tokenSymbol string) ([]types.DepositHistory, error) {
	tokenSymbol = strings.ToLower(tokenSymbol)
	page := uint64(1)
	limit := "100"
	depositHistories := []types.DepositHistory{}

	for {
		depositHistoryPath := "/api/v3/crypto/deposit-history"
		queryParma := fmt.Sprintf("?p=%d&lmt=%s", page, limit)

		body, err := b.httpRequest(depositHistoryPath, queryParma)
		if err != nil {
			return nil, err
		}

		depositHistoryResponse := types.DepositHistoryResponse{}
		err = json.Unmarshal(body, &depositHistoryResponse)
		if err != nil {
			return nil, err
		}

		for _, depositHistory := range depositHistoryResponse.Result {
			if depositHistory.Currency == tokenSymbol {
				depositHistories = append(depositHistories, depositHistory)
			}
		}

		if depositHistoryResponse.Pagination.Last == page {
			break
		}
		page++
	}

	return depositHistories, nil
}

func (b *bitkubApiClient) RequestWithdrawHistories(tokenSymbol string) ([]types.WithdrawHistory, error) {
	tokenSymbol = strings.ToLower(tokenSymbol)
	page := uint64(1)
	limit := "100"
	withdrawHistories := []types.WithdrawHistory{}

	for {
		withdrawHistoryPath := "/api/v3/crypto/withdraw-history"
		queryParma := fmt.Sprintf("?p=%d&lmt=%s", page, limit)

		body, err := b.httpRequest(withdrawHistoryPath, queryParma)
		if err != nil {
			return nil, err
		}

		withdrawHistoryResponse := types.WithdrawHistoryResponse{}
		err = json.Unmarshal(body, &withdrawHistoryResponse)
		if err != nil {
			return nil, err
		}

		for _, depositHistory := range withdrawHistoryResponse.Result {
			if depositHistory.Currency == tokenSymbol {
				withdrawHistories = append(withdrawHistories, depositHistory)
			}
		}

		if withdrawHistoryResponse.Pagination.Last == page {
			break
		}
		page++
	}

	return withdrawHistories, nil
}

func (b *bitkubApiClient) httpRequest(path, queryParam string) ([]byte, error) {
	nowMilliSec := time.Now().UnixMilli()
	nowMilliSecStr := strconv.FormatInt(nowMilliSec, 10)
	method := "GET"
	payload := nowMilliSecStr + method + path + queryParam

	signature := b.genSignature(payload)

	client := &http.Client{}
	req, err := http.NewRequest("GET", b.baseUrl+path+queryParam, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-BTK-TIMESTAMP", nowMilliSecStr)
	req.Header.Add("X-BTK-APIKEY", b.apiKey)
	req.Header.Add("X-BTK-SIGN", signature)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (b *bitkubApiClient) genSignature(payload string) string {
	secretBytes := []byte(b.apiSecret)
	payloadBytes := []byte(payload)

	hmac := hmac.New(sha256.New, secretBytes)
	hmac.Write(payloadBytes)
	signature := hmac.Sum(nil)

	return hex.EncodeToString(signature)
}
