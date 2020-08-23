package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/kyawmyintthein/orange-contrib/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type MongodbConfig struct {
	DatabaseName  string `mapstructure:"db_name" json:"db_name"`
	DatabaseHosts string `mapstructure:"hosts" json:"hosts"`
	TimeOut       int    `mapstructure:"timeout" json:"timeout"`
	DialTimeOut   int64  `mapstructure:"dial_timeout" json:"dial_timeout"`
	PoolSize      int    `mapstructure:"pool_size" json:"pool_size"`
	Username      string `mapstructure:"username" json:"username"`
	Password      string `mapstructure:"password" json:"password"`
	ReplicaSet    string `mapstructure:"replica_set" json:"replica_set"`
	AuthSource    string `mapstructure:"auth_source" json:"auth_source"`
	URI           string `mapstructure:"uri" json:"uri"`
}

type MongodbConnector interface {
	DB(context.Context) *mongo.Database
	Client(context.Context) *mongo.Client
	Config() MongodbConfig
}

type mongodbConnector struct {
	cfg    *MongodbConfig
	db     *mongo.Database
	client *mongo.Client
}

func NewMongodbConnector(cfg *MongodbConfig) (MongodbConnector, error) {
	mongodbConnector := &mongodbConnector{
		cfg: cfg,
	}

	err := mongodbConnector.connect()
	if err != nil {
		return mongodbConnector, err
	}
	return mongodbConnector, nil
}

func (this *mongodbConnector) connect() error {
	var (
		connectOnce sync.Once
		err         error
		client      *mongo.Client
	)

	connectOnce.Do(func() {
		connStr := getConnectionString(this.cfg)
		client, err = mongo.NewClient(options.Client().ApplyURI(connStr))
		if err != nil {
			logx.Errorf(context.Background(), err, "Failed to connect to database: %s", this.cfg.DatabaseName)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(this.cfg.DialTimeOut))
		defer cancel()
		err = client.Connect(ctx)
		if err != nil {
			logx.Errorf(context.Background(), err, "Failed to connect to database: %s", this.cfg.DatabaseName)
			return
		}
	})
	if err != nil {
		return err
	}

	this.client = client
	this.db = this.client.Database(this.cfg.DatabaseName)
	return nil
}

func (this *mongodbConnector) DB(ctx context.Context) *mongo.Database {
	var rp readpref.ReadPref
	err := this.client.Ping(ctx, &rp)
	if err != nil {
		logx.Errorf(ctx, err, "fail to ping %s", this.cfg.DatabaseHosts)
	}
	return this.db
}
func (this *mongodbConnector) Client(ctx context.Context) *mongo.Client {
	return this.client
}

func (this *mongodbConnector) Config() MongodbConfig {
	return *this.cfg
}

func (this *mongodbConnector) EnsureIndex(collection *mongo.Collection, indexMap map[string]*options.IndexOptions) error {

	indexView := collection.Indexes()
	for k, index := range indexMap {
		if isCompositeKey(k) {
			doc := bsonx.Doc{}
			allKeys := strings.Split(k, "-")
			for i := 0; i < len(allKeys); i++ {
				elem := bsonx.Elem{allKeys[i], bsonx.Int32(int32(1))}
				doc = append(doc, elem)
				logx.InfoKVf(context.Background(), logx.KV{"Identifier": allKeys[i]}, "index check for key")
			}
			indexModel := mongo.IndexModel{Keys: doc, Options: index}
			_, err := indexView.CreateOne(context.Background(), indexModel)

			if err != nil {
				logx.Errorf(context.Background(), err, "fail to create %s", k)
			}
		} else {
			keys := bsonx.Doc{{Key: k, Value: bsonx.Int32(int32(1))}}
			indexModel := mongo.IndexModel{Keys: keys, Options: index}
			_, err := indexView.CreateOne(context.Background(), indexModel)
			logx.InfoKVf(context.Background(), logx.KV{"Identifier": k}, "index check for key")
			if err != nil {
				logx.Errorf(context.Background(), err, "fail to create %s", k)
			}
		}
	}

	return nil
}

func isCompositeKey(key string) bool {
	return len(strings.Split(key, "-")) > 1
}

func getConnectionString(config *MongodbConfig) string {
	if config.URI != "" {
		return config.URI
	}

	var b bytes.Buffer
	b.WriteString("mongodb://")
	if config.Username != "" {
		b.WriteString(config.Username)
		b.WriteString(":")
	}
	if config.Password != "" {
		b.WriteString(config.Password)
		b.WriteString("@")
	}
	b.WriteString(config.DatabaseHosts)
	b.WriteString("/")

	var urlQueryString []string

	if config.PoolSize != 0 {
		urlQueryString = append(urlQueryString, fmt.Sprintf("maxPoolSize=%d", config.PoolSize))
	}

	if config.ReplicaSet != "" {
		urlQueryString = append(urlQueryString, fmt.Sprintf("replicaSet=%s", config.ReplicaSet))
	}

	if config.AuthSource != "" {
		urlQueryString = append(urlQueryString, fmt.Sprintf("authSource=%s", config.AuthSource))
	}

	if len(urlQueryString) > 0 {
		s := strings.Join(urlQueryString, "&")
		s = "?" + s
		b.WriteString(s)
	}

	return b.String()
}

func IsDuplicateError(err error) bool {
	if err, ok := err.(mongo.WriteException); ok {
		if err.WriteErrors[0].Code == 11000 {
			return true
		}
	}
	return false
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == mongo.ErrNoDocuments.Error()
}

func UniqueIndex() *options.IndexOptions {
	uniqueIndex := options.Index()
	uniqueIndex.SetUnique(true)
	uniqueIndex.SetBackground(true)
	return uniqueIndex
}

func SparseIndex() *options.IndexOptions {
	sparseIndex := options.Index()
	sparseIndex.SetSparse(true)
	sparseIndex.SetBackground(true)
	return sparseIndex
}

func SparseUniqueIndex() *options.IndexOptions {
	sparseUniqueIndex := options.Index()
	sparseUniqueIndex.SetSparse(true)
	sparseUniqueIndex.SetUnique(true)
	sparseUniqueIndex.SetBackground(true)
	return sparseUniqueIndex
}
