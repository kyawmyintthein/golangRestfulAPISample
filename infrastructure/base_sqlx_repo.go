package infrastructure

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-sql-driver/mysql"
	jsoniter "github.com/json-iterator/go"
	"github.com/kyawmyintthein/orange-contrib/logx"
	"github.com/kyawmyintthein/orange-contrib/optionx"
	"github.com/kyawmyintthein/orange-contrib/tracingx/newrelicx"
	newrelic "github.com/newrelic/go-agent"
)

type Criteria interface{}

type SqlRepositoryCfg struct {
	CustomIDPrefix      string `mapstructure:"custom_id_prefix" json:"custom_id_prefix"`
	EnabledNRTracer     bool   `mapstructure:"enabled_newrelic_tracer" json:"enabled_newrelic_tracer"`
	InsertMarshallerTag string `mapstructure:"insert_marshaller_tag" json:"insert_marshaller_tag"`
	UpdateMarshallerTag string `mapstructure:"update_marshaller_tag" json:"update_marshaller_tag"`
	DeleteMarshallerTag string `mapstructure:"delete_marshaller_tag" json:"delete_marshaller_tag"`
}

type BaseSqlRepository struct {
	Config           *SqlRepositoryCfg
	NRTracer         newrelicx.NewrelicTracer
	DBConnector      SqlDBConnector
	InsertMarshaller jsoniter.API
	UpdateMarshaller jsoniter.API
	DeleteMarshaller jsoniter.API
}

func NewBaseSqlRepository(cfg *SqlRepositoryCfg, dbConnector SqlDBConnector, opts ...optionx.Option) *BaseSqlRepository {
	options := optionx.NewOptions(opts...)
	base := &BaseSqlRepository{
		Config:      cfg,
		DBConnector: dbConnector,
		InsertMarshaller: jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			TagKey:                 cfg.InsertMarshallerTag,
			ValidateJsonRawMessage: true,
		}.Froze(),
		UpdateMarshaller: jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			TagKey:                 cfg.UpdateMarshallerTag,
			ValidateJsonRawMessage: true,
		}.Froze(),
		DeleteMarshaller: jsoniter.Config{
			EscapeHTML:             true,
			SortMapKeys:            true,
			TagKey:                 cfg.DeleteMarshallerTag,
			ValidateJsonRawMessage: true,
		}.Froze(),
	}

	// set newrelic
	newRelicTracer, ok := options.Context.Value(newrelicTracerKey{}).(newrelicx.NewrelicTracer)
	if newRelicTracer == nil && !ok {
		base.Config.EnabledNRTracer = false
	} else {
		base.NRTracer = newRelicTracer
		base.Config.EnabledNRTracer = true
	}
	return base
}

func (base *BaseSqlRepository) NrSegment(ctx context.Context, tableName string, opName string) newrelic.DatastoreSegment {
	txn := newrelic.FromContext(ctx)
	dbConfig := base.DBConnector.Config()
	var datastoreProduct newrelic.DatastoreProduct
	switch dbConfig.Driver {
	case Mysql:
		datastoreProduct = newrelic.DatastoreMySQL
		break
	case Postgres:
		datastoreProduct = newrelic.DatastorePostgres
		break
	case Sqlite3:
		datastoreProduct = newrelic.DatastoreSQLite
		break
	}
	return base.NRTracer.RecordDatastoreMetric(txn,
		datastoreProduct,
		fmt.Sprintf("Table.%s", tableName),
		opName,
		"",
		nil,
		dbConfig.DBHost,
		"",
		dbConfig.DBName)
}

func (base *BaseSqlRepository) FindAllQuery(ctx context.Context, criteria Criteria, table string) string {

	var queryString bytes.Buffer

	if criteria == "" {
		return fmt.Sprintf("SELECT * FROM %s;", table)
	}

	queryString.WriteString(fmt.Sprintf("%v", criteria))
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s;", table, queryString.String())
	logx.Debugf(ctx, "FindAllQuery %s, %v, %v \n", query, table, criteria)
	return query
}

func (base *BaseSqlRepository) FindAllQueryWithOrder(ctx context.Context, criteria Criteria, table string, orderBy string, isAsc bool) string {
	var queryString bytes.Buffer

	if criteria == "" {
		return fmt.Sprintf("SELECT * FROM %s;", table)
	}

	queryString.WriteString(fmt.Sprintf("%v", criteria))

	order := "ASC"
	if !isAsc {
		order = "DESC"
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s ORDER BY %s %s;", table, queryString.String(), orderBy, order)
	logx.Debugf(ctx, "FindAllQueryWithOrder %s, %v, %v , %v \n", query, table, orderBy, isAsc)
	return query
}

func (base *BaseSqlRepository) FindOneQuery(ctx context.Context, criteria Criteria, table string) string {
	var queryString bytes.Buffer

	if criteria == "" {
		return fmt.Sprintf("SELECT * FROM %s limit 1;", table)
	}
	queryString.WriteString(fmt.Sprintf("%v", criteria))

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s limit 1;", table, queryString.String())
	logx.Debugf(ctx, "FindOneQuery %s, %v, %v \n", query, table, criteria)
	return query
}

func (base *BaseSqlRepository) InsertStatement(ctx context.Context, table string, modelStructPtr interface{}) (string, error) {
	valueOfIStructPointer := reflect.ValueOf(modelStructPtr)

	if k := valueOfIStructPointer.Kind(); k != reflect.Ptr {
		return "", fmt.Errorf("model should be pointer type.")
	}

	valueOfIStructPointerElem := valueOfIStructPointer.Elem()

	// Below is a further (and definitive) check regarding settability in addition to checking whether it is a pointer earlier.
	if !valueOfIStructPointerElem.CanSet() {
		return "", fmt.Errorf("unable to set value to type!")
	}

	var modelValues map[string]interface{}
	var values []interface{}
	dataBytes, err := base.InsertMarshaller.Marshal(modelStructPtr)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(dataBytes, &modelValues)
	if err != nil {
		return "", err
	}

	var columns bytes.Buffer
	var valuePlaceHolders bytes.Buffer
	valueLen := len(modelValues)
	i := 0
	for k, val := range modelValues {
		columns.WriteString(k)
		valuePlaceHolders.WriteString(fmt.Sprintf(":%s", k))
		if valueLen-1 != i {
			valuePlaceHolders.WriteString(",")
			columns.WriteString(",")
		}
		values = append(values, val)
		i++
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, table, columns.String(), valuePlaceHolders.String())
	logx.Debugf(ctx, "InsertStatement %s, %v, %v \n", q, table, modelStructPtr)
	return q, nil
}

func (base *BaseSqlRepository) UpdateStatement(ctx context.Context, criteria Criteria, table string, modelStructPtr interface{}) (string, error) {

	valueOfIStructPointer := reflect.ValueOf(modelStructPtr)

	if k := valueOfIStructPointer.Kind(); k != reflect.Ptr {
		return "", fmt.Errorf("model should be pointer type.")
	}

	valueOfIStructPointerElem := valueOfIStructPointer.Elem()

	// Below is a further (and definitive) check regarding settability in addition to checking whether it is a pointer earlier.
	if !valueOfIStructPointerElem.CanSet() {
		return "", fmt.Errorf("unable to set value to type!")
	}

	var modelValues map[string]interface{}
	dataBytes, err := base.UpdateMarshaller.Marshal(modelStructPtr)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(dataBytes, &modelValues)
	if err != nil {
		return "", err
	}

	var colValues bytes.Buffer
	valueLen := len(modelValues)
	i := 0
	for k, _ := range modelValues {
		colValues.WriteString(k)
		colValues.WriteString("=")
		colValues.WriteString(fmt.Sprintf(":%s", k))
		if valueLen-1 != i {
			colValues.WriteString(",")
		}
		i++
	}

	var queryString bytes.Buffer
	queryString.WriteString(fmt.Sprintf("%v", criteria))

	q := fmt.Sprintf(`UPDATE %s SET %s WHERE %s`, table, colValues.String(), queryString.String())
	logx.Debugf(ctx, "UpdateStatement %s, %v , %+v \n", q, table, modelStructPtr)
	return q, nil
}

func (base *BaseSqlRepository) DeleteStatement(ctx context.Context, criteria Criteria, table string) (string, error) {
	var queryString bytes.Buffer
	queryString.WriteString(fmt.Sprintf("%v", criteria))
	q := fmt.Sprintf(`DELETE  FROM %s  WHERE %s`, table, criteria)
	logx.Debugf(ctx, "Delete statement %s, %v , %+v \n", q)
	return q, nil
}

// TODO:: complete this for other drivers Postgres, Sqlite3
// Note:: only support for mysql
func (base *BaseSqlRepository) IsDuplicateError(ctx context.Context, err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	if mysqlErr.Number == 1062 {
		return true
	}
	return false
}

// TODO:: complete this for other drivers Postgres, Sqlite3
// Note:: only support for mysql
func (base *BaseSqlRepository) IsNotFoundError(ctx context.Context, err error) bool {
	if sql.ErrNoRows.Error() == err.Error() {
		return true
	}
	return false
}
