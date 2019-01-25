package infrastructure

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/constant/ecodes"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"strings"
	"sync"
)

type MongoStore interface{
	DB() *mongo.Database
	Session() *mongo.Client
	DatabaseName() (string, error)
}

type mongoDatastore struct {
	db      *mongo.Database
	session *mongo.Client
	logging logging.Logger
}

func NewMongoStore(config *config.GeneralConfig, logging logging.Logger) (MongoStore, error) {
	db, session, err := connect(config)
	if err != nil{
		return nil, err
	}

	return &mongoDatastore{
		db: db,
		logging: logging,
		session: session,
	}, nil
}

func connect(generalConfig *config.GeneralConfig) (db *mongo.Database, client *mongo.Client, err error) {
	var connectOnce sync.Once
	connectOnce.Do(func() {
		db, client, err = connectToMongo(generalConfig)
	})
	return db, client, err
}

func connectToMongo(generalConfig *config.GeneralConfig) (db *mongo.Database, client *mongo.Client, err error) {
	client, err = mongo.NewClient(generalConfig.Mongodb.Host)
	if err != nil {
		return db, client, err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return db, client, err
	}

	db = client.Database(generalConfig.Mongodb.Database)

	return db, client, nil
}

func (mongoStore *mongoDatastore) DB() *mongo.Database{
	return mongoStore.db
}

func (mongoStore *mongoDatastore) DatabaseName() (string, error){
	if mongoStore.db == nil{
		return "", errors.New(ecodes.DatabaseConnnectionFailed, constant.DatabaseConnnectionFailedErr)
	}
	return mongoStore.db.Name(), nil
}

func (mongoStore *mongoDatastore) Session() *mongo.Client{
	return mongoStore.session
}

func (mongoStore *mongoDatastore) EnsureIndex(collection *mongo.Collection, indexMap map[string]bsonx.Doc) error {
	log := mongoStore.logging.GetLogger(context.Background())
	indexView := collection.Indexes()

	for k, index := range indexMap {
		if isCompositeKey(k) {
			var doc bsonx.Doc
			allKeys := strings.Split(k, "-")
			for i := 0; i < len(allKeys); i++ {
				doc = append(doc, bsonx.Elem{allKeys[i], bsonx.Int32(1)})
				log.WithField(constant.Identifier, allKeys[i]).Info("index check for key")
			}
			indexModel := mongo.IndexModel{Keys: doc, Options: index}
			_, err := indexView.CreateOne(context.Background(), indexModel)

			if err != nil {
				log.WithError(err).Fatal("fail to create")
			}
		} else {
			doc := bsonx.Doc{bsonx.Elem{k, bsonx.Int32(1)}}
			indexModel := mongo.IndexModel{Keys: doc, Options: index}
			_, err := indexView.CreateOne(context.Background(), indexModel)

			log.WithField(constant.Identifier, k).Info("index check for key")
			if err != nil {
				log.WithError(err).Fatal("fail to create")
			}
		}
	}

	return nil
}

func isCompositeKey(key string) bool {
	return len(strings.Split(key, "-")) > 1
}

func (mongoStore *mongoDatastore) UniqueIndex()  bsonx.Doc {
	indexBuilder := mongo.NewIndexOptionsBuilder()
	indexBuilder.Unique(true)
	indexBuilder.Background(true)
	return indexBuilder.Build()
}

func (mongoStore *mongoDatastore) SparseIndex() bsonx.Doc {
	indexBuilder := mongo.NewIndexOptionsBuilder()
	indexBuilder.Sparse(true)
	indexBuilder.Background(true)
	return indexBuilder.Build()
}

func (mongoStore *mongoDatastore) SparseUniqueIndex() bsonx.Doc {
	indexBuilder := mongo.NewIndexOptionsBuilder()
	indexBuilder.Sparse(true)
	indexBuilder.Unique(true)
	indexBuilder.Background(true)
	return indexBuilder.Build()
}
