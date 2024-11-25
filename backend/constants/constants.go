package constants

const (
	Examitsu  = "examitsu"
	Success   = "success"
	JWTClaims = "jwt_claims"
	Error     = "error"
	Worker    = "worker"

	SystemUser  = "SYSTEM"
	ProctorUser = "PROCTOR"

	ID          = "id"
	Serial      = "serial"
	IsOpen      = "is_open"
	ExamID      = "exam_id"
	OrderNumber = "order_number"
	None        = "NONE"

	QueryParameterPage                 = "page"
	DefaultValueQueryParameterPage     = "1"
	QueryParameterPageSize             = "page_size"
	DefaultValueQueryParameterPageSize = "10"
	DefaultQueryPaginationPage         = 1
	DefaultQueryPaginationPageSize     = 10

	InsertionBatchSize = 100

	ExamSessionSubmissionCacheObjectKeyPrefix = "ExamSessionSubmissionCacheObject"
	UpdateAnswerQueueName                     = "updateAnswerQueue"
	UpdateAnswerConsumerName                  = "updateAnswerConsumer"

	DefaultRandomQuestionBlobFilenameLength = 64
)
