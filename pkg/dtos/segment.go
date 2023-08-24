package dtos

type SegmentDtoRequest struct {
	Name string `json:"slug"`
}

type SegmentDtoResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"slug"`
}
