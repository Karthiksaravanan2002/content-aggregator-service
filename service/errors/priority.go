package errors

// Lower number = higher priority
var priorityMap = map[int]int{
	400: 1, // BadRequest, Unsupported, Validation errors → client error wins
	401: 2,
	403: 3,
	404: 4,
	422: 5,
	429: 6,

	502: 7, // Provider failure
	503: 8,

	500: 9, // Internal server error (least important)
}

// errorPriority returns a sortable number for a given HTTP status.
func ErrorPriority(status int) int {
	if p, exists := priorityMap[status]; exists {
		return p
	}
	return 99 // unknown = lowest priority
}

// PickBetter returns the error with *higher priority*.
// If current is nil → return newErr.
// If newErr has higher priority → replace current.
// Else keep existing.
func PickBetter(current AppError, newErr AppError) AppError {
	if newErr == nil {
		return current
	}
	if current == nil {
		return newErr
	}

	if ErrorPriority(newErr.StatusCode()) < ErrorPriority(current.StatusCode()) {
		return newErr
	}
	return current
}
