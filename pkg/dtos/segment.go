package dtos

import "market/pkg/models"

type SegmentDtoRequest struct {
	Name    string `json:"slug"`
	Percent int    `json:"percent,omitempty"`
}

type SegmentDtoResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"slug"`
}

type CreateSegmentDto struct {
	Name       string `json:"slug"`
	DeleteTime string `json:"ttl,omitempty"`
}

func ToSegmentDto(segment *models.Segment) SegmentDtoResponse {
	return SegmentDtoResponse{
		ID:   segment.ID,
		Name: segment.Name,
	}
}
