package errors

import "errors"

var (
	ErrFailedToConvertInsertedIdToObjectId = errors.New("failed to convert inserted ID to ObjectID")
)
