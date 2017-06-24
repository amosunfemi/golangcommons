package utility

import (
	"fmt"
	"time"

	r "github.com/dancannon/gorethink"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)
import log "gopkg.in/inconshreveable/log15.v2"

// DataStore ...
type DataStore struct {
	StoreType string
	RDBMS     RDBMSImpl
	RethinkDB RethinkDBImpl
}

// InitSchema ...
func (datastore *DataStore) InitSchema(values []interface{}) {
	if datastore.StoreType == "rdbms" {
		//for model := range values {
		//fmt.Println(model)
		datastore.RDBMS.MigrateModels(values)
		//}
	}
}

//SaveDatabaseObject ...
func (datastore *DataStore) SaveDatabaseObject(p interface{}) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.SaveDatabaseObject(p)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

//UpdateDatabaseObject ...
func (datastore *DataStore) UpdateDatabaseObject(dbobj interface{}, changeVal map[string]interface{}) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.UpdateDatabaseObject(dbobj, changeVal)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

//SearchAnyGenericObject ... This method search through with any map given to search a table
func (datastore *DataStore) SearchAnyGenericObject(queryDict map[string]interface{}, entity interface{}) (interface{}, *DataStoreError) {
	var err *DataStoreError
	fmt.Println(datastore.StoreType)
	if datastore.StoreType == "rdbms" {
		_, errInt, rdbmsErr := datastore.RDBMS.SearchAnyGenericObject(queryDict, entity)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return entity, err
}

//SearchAnyGenericObjectList ... This method search through with any map given to search a table
func (datastore *DataStore) SearchAnyGenericObjectList(queryDict map[string]interface{}, entity interface{}, relatedent ...interface{}) *DataStoreError {
	var err *DataStoreError
	//entityArr := []interface{}{}
	var errInt int
	var rdbmsErr error
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr = datastore.RDBMS.SearchAnyGenericObjectList(queryDict, entity, relatedent)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

//SearchAnyGenericObjectOr ...
func (datastore *DataStore) SearchAnyGenericObjectOr(andQueryDict map[string]interface{}, orQueryDict map[string]interface{}, entity interface{}) (interface{}, *DataStoreError) {
	var err *DataStoreError
	fmt.Println(datastore.StoreType)
	if datastore.StoreType == "rdbms" {
		_, errInt, rdbmsErr := datastore.RDBMS.SearchAnyGenericObjectOr(andQueryDict, orQueryDict, entity)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return entity, err
}

//FetchAllGenericObject ... This method search through with any map given to search a table
func (datastore *DataStore) FetchAllGenericObject(entity interface{}) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.FetchAllGenericObject(entity)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

//FetchAllGenericObjectWithParam ... This method search through with any map given to search a table
func (datastore *DataStore) FetchAllGenericObjectWithParam(queryDict map[string]interface{}, entity interface{}) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.FetchAllGenericObjectWithParam(queryDict, entity)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

//FetchFirstGenericObject ...
func (datastore *DataStore) FetchFirstGenericObject(queryDict map[string]interface{}, entity interface{}) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.FetchFirstGenericObject(queryDict, entity)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
			return err
		}
	}
	return nil
}

//FetchRelatedObject ...
func (datastore *DataStore) FetchRelatedObject(queryDict map[string]interface{}, entity interface{}, relatedEntity interface{}, foreignKeys string) *DataStoreError {
	var err *DataStoreError
	if datastore.StoreType == "rdbms" {
		errInt, rdbmsErr := datastore.RDBMS.FetchRelatedObject(queryDict, entity, relatedEntity, foreignKeys)
		if rdbmsErr != nil {
			err = datastore.composeCustomError(rdbmsErr, errInt)
		}
	}
	return err
}

func (datastore *DataStore) composeCustomError(err error, errInt int) *DataStoreError {
	errCustom := new(DataStoreError)
	errCustom.ErrNo = errInt
	errCustom.Msg = err.Error()
	errCustom.When = time.Now()
	return errCustom
}

//DataStoreError ...
type DataStoreError struct {
	CustomErr // Err is an embedded struct - ErrWidget_A inherits it's data and behavior
}

// Error ...
func (derror DataStoreError) Error() string {
	return fmt.Sprintf("%s [%d] %s", derror.When, derror.ErrNo, derror.Msg)
}

// InitDataStore ...
func (datastore *DataStore) InitDataStore(storeType string, connDetails map[string]interface{}, log log.Logger) error {
	var err error
	datastore.StoreType = storeType
	if storeType == "rdbms" {
		rdbms := RDBMSImpl{}
		err = rdbms.Init(connDetails["dbtype"].(string), connDetails["host"].(string), connDetails["username"].(string), connDetails["password"].(string), int(connDetails["port"].(float64)), connDetails["database"].(string), log)
		if err != nil {
			log.Crit("Got error when connect database, the error is '%v'", err)
		}
		datastore.RDBMS.DB.LogMode(true)
		datastore.RDBMS = rdbms
	}
	return err
}

// RDBMSImpl ... take care of RDBMS Datastore
type RDBMSImpl struct {
	DB *gorm.DB
}

// Init ...
func (rdbms *RDBMSImpl) Init(dbtype string, dbhost string, username string, password string, port int, dbname string, log log.Logger) error {
	var err error
	connString := fmt.Sprintf("user=%s port=%d password=%s dbname=%s sslmode=disable", username, port, password, dbname)
	rdbms.DB, err = gorm.Open(dbtype, connString)
	if err != nil {
		log.Crit("Got error when connect database, the error is '%v'", err)
		return err
	}
	rdbms.DB.LogMode(true)
	fmt.Println("Connected")
	return nil
}

//MigrateModels ...
func (rdbms *RDBMSImpl) MigrateModels(values []interface{}) {
	for _, model := range values {
		rdbms.DB.CreateTable(model)
	}
}

//SaveDatabaseObject ...
func (rdbms *RDBMSImpl) SaveDatabaseObject(p interface{}) (int, error) {
	var err error
	var errInt = 1000
	tx := rdbms.DB.Begin()
	txerr := tx.Create(p).Error
	if txerr != nil {
		tx.Rollback()
		err = txerr
		errInt = 1002
		fmt.Println(err.Error())
		return errInt, err
	} else {
		tx.Commit()
		return errInt, err
	}

}

//UpdateDatabaseObject ...
func (rdbms *RDBMSImpl) UpdateDatabaseObject(dbobj interface{}, changeVal map[string]interface{}) (int, error) {
	var err error
	var errInt = 1000
	tx := rdbms.DB.Begin()
	txerr := tx.Save(dbobj).Error
	if txerr != nil {
		fmt.Println(txerr)
		tx.Rollback()
		err = txerr
		errInt = 1002
	}
	tx.Commit()
	return errInt, err
}

//SearchAnyGenericObject ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) SearchAnyGenericObject(queryDict map[string]interface{}, entity interface{}) (interface{}, int, error) {
	query := rdbms.DB.Where(queryDict).Find(entity)
	if rdbms.DB.Where(queryDict).Find(entity).RecordNotFound() {
		return entity, 1001, query.Error
	} else if query.Error != nil {
		return entity, 1002, query.Error
	}
	return entity, 1000, nil
}

//SearchAnyGenericObjectList ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) SearchAnyGenericObjectList(queryDict map[string]interface{}, entity interface{}, relatedent ...interface{}) (int, error) {
	//entityArr := []interface{}{}
	query := rdbms.DB.Where(queryDict).Find(entity)

	if rdbms.DB.Where(queryDict).Find(entity).RecordNotFound() {
		return 1001, query.Error
	} else if query.Error != nil {
		return 1002, query.Error
	}

	for _, rentity := range relatedent {
		rdbms.DB.Model(&entity).Related(&rentity)
	}

	return 1000, nil
}

//SearchAnyGenericObjectOr ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) SearchAnyGenericObjectOr(andQueryDict map[string]interface{}, orQueryDict map[string]interface{}, entity interface{}) (interface{}, int, error) {
	query := rdbms.DB.Where(andQueryDict).Or(orQueryDict).Find(entity)
	if query.RecordNotFound() {
		return entity, 1001, query.Error
	} else if query.Error != nil {
		return entity, 1002, query.Error
	}
	return entity, 1000, nil
}

//FetchAllGenericObject ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) FetchAllGenericObject(entity interface{}) (int, error) {
	query := rdbms.DB.Find(entity)
	if query.RecordNotFound() {
		return 1001, query.Error
	} else if query.Error != nil {
		return 1002, query.Error
	}
	return 1000, nil
}

//FetchAllGenericObjectWithParam ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) FetchAllGenericObjectWithParam(queryDict map[string]interface{}, entity interface{}) (int, error) {
	query := rdbms.DB.Where(queryDict).Find(entity)
	if query.RecordNotFound() {
		return 1001, query.Error
	} else if query.Error != nil {
		return 1002, query.Error
	}
	return 1000, nil
}

//FetchFirstGenericObject ... This method search through with any map given to search a table
func (rdbms *RDBMSImpl) FetchFirstGenericObject(queryDict map[string]interface{}, entity interface{}) (int, error) {
	query := rdbms.DB.Where(queryDict).First(entity)
	if query.RecordNotFound() {
		return 1001, query.Error
	} else if query.Error != nil {
		return 1002, query.Error
	}
	return 1000, nil
}

//FetchRelatedObject ... Fetch Related Objects via the primary key
func (rdbms *RDBMSImpl) FetchRelatedObject(queryDict map[string]interface{}, entity interface{}, relatedEntity interface{}, foreignKeys string) (int, error) {

	query := rdbms.DB.Where(queryDict).Find(entity)
	rdbms.DB.Model(entity).Related(relatedEntity, foreignKeys)
	if query.RecordNotFound() {
		return 1001, query.Error
	} else if query.Error != nil {
		return 1002, query.Error
	}
	return 1000, nil
}

//CreateDatabaseObjects ...
func (rdbms *RDBMSImpl) CreateDatabaseObjects(models ...interface{}) {
	rdbms.DB.AutoMigrate(models)
}

// RethinkDBImpl ... take care of RethinkDB DataStore
type RethinkDBImpl struct {
	Session *r.Session
}

// Init ...
func (rthink *RethinkDBImpl) Init(conntype string, dbhost string, username string, password string, port int, dbname string, log log.Logger, hosts ...string) error {
	var retErr error
	if conntype == "cluster" {
		session, err := r.Connect(r.ConnectOpts{
			Addresses: hosts,
			Database:  dbname,
			AuthKey:   password,
		})
		retErr = err
		if err != nil {
			log.Crit("Got error when connect database, the error is '%v'", err)
			return retErr
		}
		rthink.Session = session
	} else {
		session, err := r.Connect(r.ConnectOpts{
			Address:  dbhost,
			Database: dbname,
			AuthKey:  password,
		})
		retErr = err
		if err != nil {
			log.Crit("Got error when connect database, the error is '%v'", err)
			return retErr
		}
		rthink.Session = session
	}
	fmt.Println("Rethink Connected")
	return nil
}
