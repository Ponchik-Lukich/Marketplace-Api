package errors

import "github.com/gin-gonic/gin"

type CustomError interface {
	Error() string
	Message() string
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

	ctx.JSON(status, response)
}

func HandleCustomError(ctx *gin.Context, status int, err CustomError) {
	response := gin.H{
		"error": err.Error(),
	}

	if err.Message() != "" {
		response["message"] = err.Message()
	}

	ctx.JSON(status, response)
}

type SegmentNotFound struct{}

func (e SegmentNotFound) Error() string {
	return SegmentNotFoundErr400
}

func (e SegmentNotFound) Code() int {
	return 400
}

func (e SegmentNotFound) Message() string {
	return ""
}

type EmptySegmentName struct{}

func (e EmptySegmentName) Error() string {
	return EmptySegmentNameErr400
}

func (e EmptySegmentName) Code() int {
	return 400
}

func (e EmptySegmentName) Message() string {
	return ""
}

type SegmentAlreadyExist struct{}

func (e SegmentAlreadyExist) Error() string {
	return SegmentAlreadyExist400
}

func (e SegmentAlreadyExist) Code() int {
	return 400
}

func (e SegmentAlreadyExist) Message() string {
	return ""
}

type MissingNames struct {
	err string
}

func (e MissingNames) Error() string {
	return MissingNamesErr400
}

func (e MissingNames) Code() int {
	return 400
}

func (e MissingNames) Message() string {
	return e.err
}

type DateParsing struct {
	Err string
}

func (e DateParsing) Error() string {
	return DateParsingErr400
}

func (e DateParsing) Code() int {
	return 400
}

func (e DateParsing) Message() string {
	return e.Err
}

type UserDoesNotHaveSegment struct {
	Err string
}

func (e UserDoesNotHaveSegment) Error() string {
	return UserDoesNotHaveSegmentErr400
}

func (e UserDoesNotHaveSegment) Code() int {
	return 400
}

func (e UserDoesNotHaveSegment) Message() string {
	return e.Err
}

type UserAlreadyHasSegment struct {
	Err string
}

func (e UserAlreadyHasSegment) Error() string {
	return UserAlreadyHasSegmentErr400
}

func (e UserAlreadyHasSegment) Code() int {
	return 400
}

func (e UserAlreadyHasSegment) Message() string {
	return e.Err
}

type DeleteSegments struct {
	Err string
}

func (e DeleteSegments) Error() string {
	return DeleteSegmentsErr500
}

func (e DeleteSegments) Code() int {
	return 500
}

func (e DeleteSegments) Message() string {
	return e.Err
}

type UpdatingUser struct {
	Err string
}

func (e UpdatingUser) Error() string {
	return UpdatingUserErr500
}

func (e UpdatingUser) Code() int {
	return 500
}

func (e UpdatingUser) Message() string {
	return e.Err
}

type AddingLogs struct {
	Err string
}

func (e AddingLogs) Error() string {
	return AddingLogsErr500
}

func (e AddingLogs) Code() int {
	return 500
}

func (e AddingLogs) Message() string {
	return e.Err
}

type AddingPercent struct {
	Err string
}

func (e AddingPercent) Error() string {
	return AddingPercentErr500
}

func (e AddingPercent) Code() int {
	return 500
}

func (e AddingPercent) Message() string {
	return e.Err
}

type CountUsersNumber struct {
	Err string
}

func (e CountUsersNumber) Error() string {
	return CountUsersNumberErr500
}

func (e CountUsersNumber) Code() int {
	return 500
}

func (e CountUsersNumber) Message() string {
	return e.Err
}

type GetSegmentByName struct {
	Err string
}

func (e GetSegmentByName) Error() string {
	return GetSegmentByNameErr500
}

func (e GetSegmentByName) Code() int {
	return 500
}

func (e GetSegmentByName) Message() string {
	return e.Err
}

type CreatingFile struct {
	Err string
}

func (e CreatingFile) Error() string {
	return CreatingFileErr500
}

func (e CreatingFile) Code() int {
	return 500
}

func (e CreatingFile) Message() string {
	return e.Err
}

type WritingFile struct {
	Err string
}

func (e WritingFile) Error() string {
	return WritingFileErr500
}

func (e WritingFile) Code() int {
	return 500
}

func (e WritingFile) Message() string {
	return e.Err
}

type GettingLogs struct {
	Err string
}

func (e GettingLogs) Error() string {
	return GettingLogsErr500
}

func (e GettingLogs) Code() int {
	return 500
}

func (e GettingLogs) Message() string {
	return e.Err
}
