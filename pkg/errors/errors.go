package errors

import "github.com/gin-gonic/gin"

var (
	BindingJSONErr = "BINDING_JSON_ERR"

	CreatingSegmentErr        = "CREATING_SEGMENTS_ERR"
	DeletingSegmentErr        = "DELETING_SEGMENTS_ERR"
	SegmentNotFoundErr        = "SEGMENT_NOT_FOUND"
	EmptySegmentNameErr       = "EMPTY_SEGMENT_NAME"
	SegmentAlreadyExist       = "SEGMENT_ALREADY_EXIST"
	EmptyDateErr              = "EMPTY_DATE"
	DeleteSegmentsErr         = "DELETE_SEGMENTS_ERR"
	CreateSegmentsErr         = "CREATE_SEGMENTS_ERR"
	EmptyUserIDErr            = "EMPTY_USER_ID"
	EditingUserErr            = "EDITING_USER_ERR"
	UserDoesNotHaveSegmentErr = "USER_DOES_NOT_HAVE_SEGMENT"
	UserAlreadyHasSegmentErr  = "USER_ALREADY_HAS_SEGMENT"
	ConvertingUserIdErr       = "CONVERTING_USER_ID_ERR"
	GettingUserSegmentsErr    = "GETTING_USER_SEGMENTS_ERR"
	UpdatingUserErr           = "CREATING_OR_UPDATING_USER_ERR"
	AddingLogsErr             = "ADDING_LOGS_ERR"
	InvalidDateErr            = "INVALID_DATE_ERR"
	TimeParsingErr            = "TIME_PARSING_ERR"
	InvalidPercentErr         = "INVALID_PERCENT_ERR"
	AddingPercentErr          = "ADDING_PERCENT_SEGMENTS_ERR"
	CountUsersNumberErr       = "COUNT_USERS_NUMBER_ERR"
	GetSegmentByNameErr       = "GET_SEGMENT_BY_NAME_ERR"
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
