package models



type InputPost struct {
    Title     string    `json:"title"`
    Content   string    `json:"content"`
}

type OutputPost struct {
    ID        int    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt string`json:"created_at"`
}
