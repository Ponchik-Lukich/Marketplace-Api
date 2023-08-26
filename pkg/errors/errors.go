package errors

import "github.com/gin-gonic/gin"

type CustomError interface {
	Error() string
	Message() string
	Code() int
}

var (
	BindingJSONErr      = "BINDING_JSON_ERR"
	EmptyUserIDErr      = "EMPTY_USER_ID"
	InvalidPercentErr   = "INVALID_PERCENT_ERR"
	ConvertingUserIdErr = "CONVERTING_USER_ID_ERR"
	EmptyDateErr        = "EMPTY_DATE"

	SegmentNotFoundErr400        = "SEGMENT_NOT_FOUND"
	EmptySegmentNameErr400       = "EMPTY_SEGMENT_NAME"
	SegmentAlreadyExist400       = "SEGMENT_ALREADY_EXIST"
	MissingNamesErr400           = "SEGMENTS_DOES_NOT_EXISTS: "
	DateParsingErr400            = "TIME_PARSING_ERR"
	UserDoesNotHaveSegmentErr400 = "USER_DOES_NOT_HAVE_SEGMENT"
	UserAlreadyHasSegmentErr400  = "USER_ALREADY_HAS_SEGMENT"

	GetSegmentsByUserIdsErr500 = "GET_SEGMENTS_BY_USER_IDS_ERR"
	DeleteSegmentsErr500       = "DELETE_SEGMENTS_ERR"
	UpdatingUserErr500         = "CREATING_OR_UPDATING_USER_ERR"
	AddingLogsErr500           = "ADDING_LOGS_ERR"
	AddingPercentErr500        = "ADDING_PERCENT_SEGMENTS_ERR"
	CountUsersNumberErr500     = "COUNT_USERS_NUMBER_ERR"
	GetSegmentByNameErr500     = "GET_SEGMENT_BY_NAME_ERR"
	CreatingFileErr500         = "CREATING_FILE_ERR"
	WritingFileErr500          = "WRITING_FILE_ERR"
	GettingLogsErr500          = "GETTING_LOGS_ERR"
	CreatingSegmentErr500      = "CREATING_SEGMENT_ERR"
	GetMissingNamesErr500      = "GET_MISSING_NAMES_ERR"
	AddingSegmentsErr500       = "ADDING_SEGMENTS_ERR"
	CommitErr500               = "COMMIT_ERR"
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
	Err string
}

func (e MissingNames) Error() string {
	return MissingNamesErr400
}

func (e MissingNames) Code() int {
	return 400
}

func (e MissingNames) Message() string {
	return e.Err
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

type UpdateUser struct {
	Err string
}

func (e UpdateUser) Error() string {
	return UpdatingUserErr500
}

func (e UpdateUser) Code() int {
	return 500
}

func (e UpdateUser) Message() string {
	return e.Err
}

type AddLogs struct {
	Err string
}

func (e AddLogs) Error() string {
	return AddingLogsErr500
}

func (e AddLogs) Code() int {
	return 500
}

func (e AddLogs) Message() string {
	return e.Err
}

type AddPercent struct {
	Err string
}

func (e AddPercent) Error() string {
	return AddingPercentErr500
}

func (e AddPercent) Code() int {
	return 500
}

func (e AddPercent) Message() string {
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

type CreateFile struct {
	Err string
}

func (e CreateFile) Error() string {
	return CreatingFileErr500
}

func (e CreateFile) Code() int {
	return 500
}

func (e CreateFile) Message() string {
	return e.Err
}

type WriteFile struct {
	Err string
}

func (e WriteFile) Error() string {
	return WritingFileErr500
}

func (e WriteFile) Code() int {
	return 500
}

func (e WriteFile) Message() string {
	return e.Err
}

type GetLogs struct {
	Err string
}

func (e GetLogs) Error() string {
	return GettingLogsErr500
}

func (e GetLogs) Code() int {
	return 500
}

func (e GetLogs) Message() string {
	return e.Err
}

type CreateSegment struct {
	Err string
}

func (e CreateSegment) Error() string {
	return CreatingSegmentErr500
}

func (e CreateSegment) Code() int {
	return 500
}

func (e CreateSegment) Message() string {
	return e.Err
}

type GetSegmentsByUserId struct {
	Err string
}

func (e GetSegmentsByUserId) Error() string {
	return GetSegmentsByUserIdsErr500
}

func (e GetSegmentsByUserId) Code() int {
	return 500
}

func (e GetSegmentsByUserId) Message() string {
	return e.Err
}

type GetMissingNames struct {
	Err string
}

func (e GetMissingNames) Error() string {
	return GetMissingNamesErr500
}

func (e GetMissingNames) Code() int {
	return 500
}

func (e GetMissingNames) Message() string {
	return e.Err
}

type Transaction struct {
	Err string
}

func (e Transaction) Error() string {
	return CommitErr500
}

func (e Transaction) Code() int {
	return 500
}

func (e Transaction) Message() string {
	return e.Err
}

type AddSegments struct {
	Err string
}

func (e AddSegments) Error() string {
	return AddingSegmentsErr500
}

func (e AddSegments) Code() int {
	return 500
}

func (e AddSegments) Message() string {
	return e.Err
}
