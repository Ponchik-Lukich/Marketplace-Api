package dtos

type EditUserDtoRequest struct {
	ToCreate []string `json:"to_create"`
	ToDelete []string `json:"to_delete"`
	UserID   uint64   `json:"user_id"`
}

type GetLogsDtoRequest struct {
	UserID uint64 `json:"user_id"`
	Date   string `json:"date"`
}
