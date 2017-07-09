package orm

/*
 * 	Custom ORM package
  * 	Based on gorm ORM.
   * 	Make sure the every database actions within the transaction.
    * 	Do begin, rollback and commit the db transaction automatically.
     * 	Return proper error
      *
*/

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golangRestfulAPISample/db/gorm"
	"reflect"
	"strings"
)

type (
	DBFunc func(tx *gorm.DB) error // func type which accept *gorm.DB and return error
)

// Create
// Helper function to insert gorm model to database by using 'WithinTransaction'
func Create(v interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if !gorm.DBManager().NewRecord(v) {
			return exception.CreateRecordIsNotNew.Throw(message.SomethingWentWrong) // throw error exception
		}
		if err = tx.Create(v).Error; err != nil {
			tx.Rollback()                                                                                                           // rollback
			return exception.FailedToCreateRecord.Throw(message.UnableToCreateDatabaseRecord, ReflectInterfaceName(v), err.Error()) // throw error exception
		}
		return err
	})
}

// Save
// Helper function to save gorm model to database by using 'WithinTransaction'
func Save(v interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if gorm.DBManager().NewRecord(v) {
			return exception.InvalidRecordIsUpdated.Throw(message.SomethingWentWrong) // throw error exception
		}
		if err = tx.Save(v).Error; err != nil {
			tx.Rollback()                                                                                                       // rollback
			return exception.FailedToSaveRecord.Throw(message.UnableToSaveDatabaseRecord, ReflectInterfaceName(v), err.Error()) // throw error exception
		}
		return err
	})
}

// Delete
// Helper function to save gorm model to database by using 'WithinTransaction'
func Delete(v interface{}) error {
	return WithinTransaction(func(tx *gorm.DB) (err error) {
		// check new object
		if err = tx.Delete(v).Error; err != nil {
			tx.Rollback()                                                                                                           // rollback
			return exception.FailedToDeleteRecord.Throw(message.UnableToDeleteDatabaseRecord, ReflectInterfaceName(v), err.Error()) // throw error exception
		}
		return err
	})
}

// FindOneByID
// Helper function to find a record by using 'WithinTransaction'
func FindOneByID(v interface{}, id uint64) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Last(v, id).Error; err != nil {
			tx.Rollback() // rollback db transaction
			if err.Error() != "record not found" {
				return exception.RecordNotFound.Throw(message.RecordDoesNotExist, ReflectInterfaceName(v))
			} else {
				return exception.CorruptedData.Throw(message.ResourceIsCorrupted, ReflectInterfaceName(v))
			}
		}
		return err
	})
}

// FindAll
// Helper function to find records by using 'WithinTransaction'
func FindAll(v interface{}) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Find(v).Error; err != nil {
			tx.Rollback() // rollback db transaction
			if err.Error() != "record not found" {
				return exception.RecordNotFound.Throw(message.RecordDoesNotExist, ReflectInterfaceName(v))
			} else {
				return exception.CorruptedData.Throw(message.ResourceIsCorrupted, ReflectInterfaceName(v))
			}
		}
		return err
	})
}

// FindOneByQuery
// Helper function to find a record by using 'WithinTransaction'
func FindOneByQuery(v interface{}, params map[string]interface{}) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Where(params).Last(v).Error; err != nil {
			tx.Rollback() // rollback db transaction
			if err.Error() != "record not found" {
				return exception.RecordNotFound.Throw(message.RecordDoesNotExist, ReflectInterfaceName(v))
			} else {
				return exception.CorruptedData.Throw(message.ResourceIsCorrupted, ReflectInterfaceName(v))
			}
		}
		return err
	})
}

// FindByQuery
// Helper function to find records by using 'WithinTransaction'
func FindByQuery(v interface{}, params map[string]interface{}) (err error) {
	return WithinTransaction(func(tx *gorm.DB) error {
		if err = tx.Where(params).Find(v).Error; err != nil {
			tx.Rollback() // rollback db transaction
			if err.Error() != "record not found" {
				return exception.RecordNotFound.Throw(message.RecordDoesNotExist, ReflectInterfaceName(v))
			} else {
				return exception.CorruptedData.Throw(message.ResourceIsCorrupted, ReflectInterfaceName(v))
			}
		}
		return err
	})
}

// WithinTransaction
// accept DBFunc as parameter
// call DBFunc function within transaction begin, and commit and return error from DBFunc
func WithinTransaction(fn DBFunc) (err error) {
	tx := gorm.DBManager().Begin() // start db transaction
	defer tx.Commit()
	err = fn(tx)
	// close db transaction
	return err
}

// ReflectInterfaceName
// reflect struct name from interface
func ReflectInterfaceName(i interface{}) string {
	value := reflect.ValueOf(i)
	fullName := fmt.Sprintf("%s", value.Type())
	names := strings.Split(fullName, ".")
	return names[len(names)-1]
}
