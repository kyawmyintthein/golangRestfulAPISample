package datastore

import (
	"github.com/mongodb/mongo-go-driver/mongo"
)

func IsDuplicateError(err error) bool {
	if err, ok := err.(mongo.WriteErrors); ok {
		if (err)[0].Code == 11000 {
			return true
		}
	}

	return false
}

func IsNotFoundError(err error) bool {
	return err.Error() == mongo.ErrNoDocuments.Error()
}
