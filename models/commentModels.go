package models

import "time"

type AuthContext struct {
	UserID int64
	Role   string
}

type CreateCommentRequest struct {
	Text          string `json:"text"`
	PrevCommentID *int64 `json:"prev_comment_id"`
}

type UpdateCommentRequest struct {
	Text string `json:"text"`
}

type Comment struct {
	ID            int64     `json:"id"`
	ReviewID      int64     `json:"review_id"`
	UserID        int64     `json:"user_id"`
	Text          string    `json:"text"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	PrevCommentID int64     `json:"prev_comment_id"`
	EditedFlag    bool      `json:"edited_flag"`
	EditedAt      time.Time `json:"edited_at"`
}
