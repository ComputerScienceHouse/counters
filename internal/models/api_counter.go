package models

type GetCountersFilter struct {
	Timestamp      *int64  `json:"timestamp"`
	SortDescending bool    `json:"sort_descending"`
	Name           *string `json:"name"`
}

type GetCountersCounterPart struct {
	CounterID        string `json:"counter_id"`
	CounterTitle     string `json:"counter_title"`
	ResetDescription string `json:"reset_description"`
	ResetUsername    string `json:"reset_username"`
	ResetTimestamp   int64  `json:"reset_timestamp"`
}

type GetCountersInput struct {
	CounterID string            `json:"counter_id"`
	Filter    GetCountersFilter `json:"filter"`
}

type GetCountersOutput struct {
	Counters  []*GetCountersCounterPart `json:"counters"`
	HasMore   bool                      `json:"has_more"`
	Timestamp *int64                    `json:"timestamp"`
}

type CreateCounterWithResetInput struct {
	Counter CreateCounterInput `json:"counter"`
	Reset   CreateResetInput   `json:"reset"`
}

type CreateCounterInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeleteCounterInput struct {
	RowID string `json:"counter_id"`
}
