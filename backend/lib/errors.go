package lib

import "errors"

var (
	// handler
	ErrFailedToParseRequest   = errors.New("failed to parse request")
	ErrUnknownError           = errors.New("unknown error")
	ErrInsufficientPermission = errors.New("insufficient permission")

	// lib.jwt_claims
	ErrFailedToParseJWTClaimsInContext = errors.New("failed to parse jwt claims in context")
	ErrJWTClaimsNotFoundInContext      = errors.New("jwt claims not found in context")

	// lib.random
	ErrFailedToGenerateRandomString = errors.New("failed to generate random string")

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
	ErrFailedToUpdateMcqOption = errors.New("failed to update mcq option")
	ErrFailedToDeleteMcqOption = errors.New("failed to delete mcq option")

	// participant.repository
	ErrParticipantNotFound = errors.New("participant not found")

	// participant.service
	ErrFailedToCreateParticipants = errors.New("failed to create participants")
	ErrFailedToGetParticipants    = errors.New("failed to get participants")
	ErrFailedToUpdateParticipant  = errors.New("failed to update participant")
	ErrFailedToDeleteParticipant  = errors.New("failed to delete participant")
)
