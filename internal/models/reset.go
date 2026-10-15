package models

type Reset struct {
	ResetId     string `json:"reset_id"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
}

type ResetListPart struct {
	ResetID          string `json:"reset_id"`
	ResetDescription string `json:"reset_description"`
	ResetOwner       string `json:"reset_owner"`
	ResetUsername    string `json:"reset_username"`
	ResetOccurredAt  int64  `json:"reset_occurred_at"`
}
