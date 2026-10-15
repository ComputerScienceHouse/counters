package models

type GetResetsFilter struct {
	Timestamp      *int64 `json:"timestamp"`
	SortDescending bool   `json:"sort_descending"`
}

type GetResetsResetPart struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
	CanDelete   bool   `json:"can_delete"`
}

type GetResetsInput struct {
	CounterID string          `json:"counter_id"`
	Filter    GetResetsFilter `json:"filter"`
}

type GetResetsOutput struct {
	Resets    []*GetResetsResetPart `json:"resets"`
	HasMore   bool                  `json:"has_more"`
	Timestamp *int64                `json:"timestamp"`
}

type CreateResetInput struct {
	CounterID   string `json:"counter_id"`
	Description string `json:"description"`
	ResetTime   int64  `json:"reset_time"`
}

type DeleteResetInput struct {
	RowID string `json:"reset_id"`
}
