package errors

import "github.com/gin-gonic/gin"

type CustomError interface {
	Error() string
	Code() int
}

var (
	BindingJSONErr         = "BINDING_JSON_ERR"
	CreatingSegmentErr     = "CREATING_SEGMENTS_ERR"
	DeletingSegmentErr     = "DELETING_SEGMENTS_ERR"
	EmptyUserIDErr         = "EMPTY_USER_ID"
	EditingUserErr         = "EDITING_USER_ERR"
	InvalidPercentErr      = "INVALID_PERCENT_ERR"
	ConvertingUserIdErr    = "CONVERTING_USER_ID_ERR"
	GettingUserSegmentsErr = "GETTING_USER_SEGMENTS_ERR"
	EmptyDateErr           = "EMPTY_DATE"

	SegmentNotFoundErr400        = "SEGMENT_NOT_FOUND"
	EmptySegmentNameErr400       = "EMPTY_SEGMENT_NAME"
	SegmentAlreadyExist400       = "SEGMENT_ALREADY_EXIST"
	MissingNamesErr400           = "MISSING_NAMES: "
	DateParsingErr400            = "TIME_PARSING_ERR"
	UserDoesNotHaveSegmentErr400 = "USER_DOES_NOT_HAVE_SEGMENT"
	UserAlreadyHasSegmentErr400  = "USER_ALREADY_HAS_SEGMENT"

	DeleteSegmentsErr500   = "DELETE_SEGMENTS_ERR"
	UpdatingUserErr500     = "CREATING_OR_UPDATING_USER_ERR"
	AddingLogsErr500       = "ADDING_LOGS_ERR"
	AddingPercentErr500    = "ADDING_PERCENT_SEGMENTS_ERR"
	CountUsersNumberErr500 = "COUNT_USERS_NUMBER_ERR"
	GetSegmentByNameErr500 = "GET_SEGMENT_BY_NAME_ERR"
	CreatingFileErr500     = "CREATING_FILE_ERR"
	WritingFileErr500      = "WRITING_FILE_ERR"
	GettingLogsErr500      = "GETTING_LOGS_ERR"
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
