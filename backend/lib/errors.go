package lib

import "errors"

var (
	// handler
	ErrFailedToParseRequest        = errors.New("failed to parse request")
	ErrUnknownError                = errors.New("unknown error")
	ErrInsufficientPermission      = errors.New("insufficient permission")
	ErrFailedToDecodeContent       = errors.New("failed to decode content")
	ErrFailedToProcessUploadedFile = errors.New("failed to process uploaded file")

	// handler.participant
	ErrExamAlreadySubmitted = errors.New("exam already submitted")
	ErrExamAlreadyStarted   = errors.New("exam already started")
	ErrExamNotStarted       = errors.New("exam not started")
	ErrSessionNotFound      = errors.New("session not found")

	// lib.jwt_claims
	ErrFailedToParseJWTClaimsInContext = errors.New("failed to parse jwt claims in context")
	ErrJWTClaimsNotFoundInContext      = errors.New("jwt claims not found in context")

	// lib.random
	ErrFailedToGenerateRandomString = errors.New("failed to generate random string")

	// lib.csv
	ErrCSVRecordNotMatchedWithHeader = errors.New("csv record not matched with header")

	// adminauth.service
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrSigningMethodInvalid = errors.New("signing method invalid")
	ErrUnauthorizedRequest  = errors.New("unauthorized request")

	// exam.repository
	ErrExamNotFound = errors.New("exam not found")

	// exam.service
	ErrFailedToCreateExam      = errors.New("failed to create exam")
	ErrFailedToGetExamBySerial = errors.New("failed to get exam by serial")
	ErrFailedToGetExam         = errors.New("failed to get exam")
	ErrFailedToGetExams        = errors.New("failed to get exams")
	ErrFailedToUpdateExam      = errors.New("failed to update exam")
	ErrFailedToDeleteExam      = errors.New("failed to delete exam")

	// question.repository
	ErrQuestionNotFound = errors.New("question not found")

	// question.service
	ErrFailedToCreateQuestion  = errors.New("failed to create question")
	ErrFailedToGetQuestionByID = errors.New("failed to get question by id")
	ErrFailedToGetQuestions    = errors.New("failed to get questions")
	ErrFailedToUpdateQuestion  = errors.New("failed to update question")
	ErrFailedToDeleteQuestion  = errors.New("failed to delete question")

	// mcqoption.repository
	ErrMcqOptionNotFound = errors.New("mcq option not found")

	// mcqoption.service
	ErrFailedToCreateMcqOption = errors.New("failed to create mcq option")
	ErrFailedToGetMcqOptions   = errors.New("failed to get mcq options")
	ErrFailedToGetMcqOption    = errors.New("failed to get mcq option")
	ErrFailedToUpdateMcqOption = errors.New("failed to update mcq option")
	ErrFailedToDeleteMcqOption = errors.New("failed to delete mcq option")

	// participant.repository
	ErrParticipantNotFound = errors.New("participant not found")

	// participant.service
	ErrFailedToCreateParticipants        = errors.New("failed to create participants")
	ErrFailedToGetParticipant            = errors.New("failed to get participant")
	ErrFailedToGetParticipants           = errors.New("failed to get participants")
	ErrFailedToUpdateParticipant         = errors.New("failed to update participant")
	ErrFailedToDeleteParticipant         = errors.New("failed to delete participant")
	ErrFailedToGetParticipantTotalPoints = errors.New("failed to get participant total points")

	// submission.repository
	ErrSubmissionNotFound = errors.New("failed to get submission")

	// submission.service
	ErrFailedToSaveAnswer = errors.New("failed to save answer")
	ErrAnswerNotFound     = errors.New("answer not found")
	ErrFailedToGetAnswer  = errors.New("failed to get answer")

	// storage.service
	ErrFailedToGetUploadURL = errors.New("failed to get upload url")

	// participantsession.repository
	ErrParticipantSessionNotFound = errors.New("failed to get participant session")

	// participantsession.service
	ErrFailedToCreateParticipantSession    = errors.New("failed to create participant session")
	ErrFailedToGetParticipantSession       = errors.New("failed to get participant session")
	ErrFailedToAuthorizeParticipantSession = errors.New("failed to authorize participant session")
)
