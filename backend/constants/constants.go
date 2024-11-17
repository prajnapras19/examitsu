package constants

const (
	Examitsu  = "examitsu"
	Success   = "success"
	JWTClaims = "jwt_claims"
	Error     = "error"
	Worker    = "worker"

	SystemUser = "SYSTEM"

	ID          = "id"
	Serial      = "serial"
	ExamID      = "exam_id"
	OrderNumber = "order_number"

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
)
