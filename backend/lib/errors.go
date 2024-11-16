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

	// adminauth.service
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrSigningMethodInvalid = errors.New("signing method invalid")
	ErrUnauthorizedRequest  = errors.New("unauthorized request")

	// exam.repository
	ErrExamNotFound = errors.New("exam not found")

	// exam.service
	ErrFailedToCreateExam      = errors.New("failed to create exam")
	ErrFailedToGetExamBySerial = errors.New("failed to get exam by serial")
	ErrFailedToGetExams        = errors.New("failed to get exams")
	ErrFailedToUpdateExam      = errors.New("failed to update exam")
	ErrFailedToDeleteExam      = errors.New("failed to delete exam")

	// question.repository
	ErrQuestionNotFound = errors.New("question not found")

	// question.service
	ErrFailedToCreateQuestion = errors.New("failed to create question")
	ErrFailedToGetQuestions   = errors.New("failed to get questions")
	ErrFailedToUpdateQuestion = errors.New("failed to update question")
	ErrFailedToDeleteQuestion = errors.New("failed to delete question")
)
