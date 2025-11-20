package errors

type AppError interface {
	error 
	// StatusCode will return the http status code
	StatusCode() int
	Code() string
	Message() string

}

type appError struct {
	error
	Status int    `json:"statusCode"`
	ErrorCode      string `json:"code"`
	ErrorMessage    string `json:"message"`
}

func (e *appError) StatusCode() int {
	return e.Status
}

func (e *appError) Code() string{
	return e.ErrorCode
}

func (e *appError) Message() string{
	return e.ErrorMessage
}

func NewErr(err error) AppError{
	return &appError{
		error: err,
	}
}

func BadRequest(err error) AppError {
   errMessage := err.Error()
	 return &appError{
		error: err,
		Status: 400,
		ErrorCode: "BAD_REQUEST",
		ErrorMessage: errMessage,
	 }
}

func InternalError(err error) AppError {
   errMessage := err.Error()
	 return &appError{
		error: err,
		Status: 500,
		ErrorCode: "INTERNAL_SERVICE_ERROR",
		ErrorMessage: errMessage,
	 }
}

func BadGateway(err error) AppError {
	   errMessage := err.Error()
	 return &appError{
		error: err,
		Status: 502,
		ErrorCode: "EXTERNAL_SERVICE_ERROR",
		ErrorMessage: errMessage,
	 }
}

func ProviderError(statusCode int, err error) AppError{
		   errMessage := err.Error()
	 return &appError{
		error: err,
		Status: statusCode,
		ErrorCode: "PROVIDER_ERROR",
		ErrorMessage: errMessage,
	 }
}

// func WriteError(w http.ResponseWriter, logger *zap.Logger, err error) {
//     // Not an AppError? → Internal error wrapper
//     appErr, ok := err.(*AppError)
//     if !ok {
//         logger.Error("unhandled error type", zap.Error(err))

//         internal := InternalError(err)

//         writeJSON(w, logger, internal.StatusCode, internal)
//         return
//     }

//     // Log the error (structured)
//     logger.Error("request failed",
//         zap.String("code", appErr.Code),
//         zap.String("message", appErr.Message),
//     )

//     writeJSON(w, logger, appErr.StatusCode, appErr)
// }
