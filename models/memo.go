package models

type Memo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type CreateMemoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdateMemoRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
} 