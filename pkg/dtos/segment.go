package dtos

import "market/pkg/models"

type SegmentDtoRequest struct {
	Name string `json:"slug"`
}

type SegmentDtoResponse struct {
	ID   uint64 `json:"id"`
	Name string `json:"slug"`
}

func ToSegmentDto(segment *models.Segment) SegmentDtoResponse {
	return SegmentDtoResponse{
		ID:   segment.ID,
		Name: segment.Name,
	}
}
