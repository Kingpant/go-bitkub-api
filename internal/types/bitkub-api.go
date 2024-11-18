package types

type DepositWithdrawPaginationDetail struct {
	Page uint64 `json:"page"`
	Last uint64 `json:"last"`
}

type WithdrawHistory struct {
	TxnId    string  `json:"txn_id"`
	ExtRef   string  `json:"ext_ref"`
	Hash     string  `json:"hash"`
	Currency string  `json:"currency"`
	Amount   string  `json:"amount"`
	Fee      float64 `json:"fee"`
	Address  string  `json:"address"`
	Memo     string  `json:"memo"`
	Status   string  `json:"status"`
	Note     string  `json:"note"`
	Time     uint64  `json:"time"`
}

type WithdrawHistoryResponse struct {
	Error      uint64                          `json:"error"`
	Result     []WithdrawHistory               `json:"result"`
	Pagination DepositWithdrawPaginationDetail `json:"pagination"`
}

type DepositHistory struct {
	TxnId    string  `json:"txn_id"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"`
	Time     uint64  `json:"time"`
}

type DepositHistoryResponse struct {
	Error      uint64                          `json:"error"`
	Result     []DepositHistory                `json:"result"`
	Pagination DepositWithdrawPaginationDetail `json:"pagination"`
}

type OrderHistory struct {
	TxnId           string `json:"txn_id"`
	OrderId         string `json:"order_id"`
	Hash            string `json:"hash"`
	ParentOrderId   string `json:"parent_order_id"`
	ParentOrderHash string `json:"parent_order_hash"`
	SuperOrderId    string `json:"super_order_id"`
	SuperOrderHash  string `json:"super_order_hash"`
	ClientId        string `json:"client_id"`
	TakenByMe       bool   `json:"taken_by_me"`
	IsMaker         bool   `json:"is_maker"`
	Side            string `json:"side"`
	Type            string `json:"type"`
	Rate            string `json:"rate"`
	Fee             string `json:"fee"`
	Credit          string `json:"credit"`
	Amount          string `json:"amount"`
	Ts              uint64 `json:"ts"`
}

type OrderPaginationDetail struct {
	Page uint64 `json:"page"`
	Last uint64 `json:"last"`
	Next uint64 `json:"next"`
	Prev uint64 `json:"prev"`
}

type OrderHistoryResponse struct {
	Error      uint64                `json:"error"`
	Result     []OrderHistory        `json:"result"`
	Pagination OrderPaginationDetail `json:"pagination"`
}
