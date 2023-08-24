package errors

import "github.com/gin-gonic/gin"

var (
	BindingJSONErr = "BINDING_JSON_ERR"

	CreatingSegmentErr  = "CREATING_SEGMENTS_ERR"
	DeletingSegmentErr  = "DELETING_SEGMENTS_ERR"
	SegmentNotFoundErr  = "SEGMENT_NOT_FOUND"
	SegmentsNotFoundErr = "SEGMENTS_NOT_FOUND"
	EmptySegmentNameErr = "EMPTY_SEGMENT_NAME"
	SegmentAlreadyExist = "SEGMENT_ALREADY_EXIST"
)

func HandleError(ctx *gin.Context, status int, errMsg string, err error) {
	response := gin.H{
		"error": errMsg,
	}

	if err != nil {
		response["message"] = err.Error()
	}

	ctx.JSON(status, response)
}
