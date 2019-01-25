package datastore

import (
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/config"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logging"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"golang.org/x/net/context"
)

const(
	userCollectionName = `users`
)
type UserDatastoreInterface interface{
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, criteria interface{}, user *model.User) error
	Delete(ctx context.Context, criteria interface{}) error
	FindOne(ctx context.Context, criteria interface{}, user *model.User) error
	FindAll(ctx context.Context, criteria interface{}) ([]*model.User,error)
}

type UserDatastore struct{
	config              *config.GeneralConfig
	logging             logging.Logger
	mongoStore          infrastructure.MongoStore
	collection          *mongo.Collection
}

func NewUserDatastore(ctx context.Context, config *config.GeneralConfig, logger logging.Logger, mongoStore infrastructure.MongoStore) UserDatastoreInterface{
	userDatastore := &UserDatastore{
		config: config,
		logging: logger,
		mongoStore: mongoStore,
	}
	userDatastore.collection = userDatastore.mongoStore.DB().Collection(userCollectionName)
	return userDatastore
}

func (d *UserDatastore) Create(ctx context.Context, user *model.User) error{
	result, err := d.collection.InsertOne(ctx, bson.D{{"name", user.Name}})
	if err != nil{
		return err
	}
	objectID := result.InsertedID.(objectid.ObjectID)
	user.ID = objectID.Hex()
	return nil
}

func (d *UserDatastore) Update(ctx context.Context, criteria interface{}, user *model.User) error{
	doc := bson.D{{"$set", bson.D{
		{"name", user.Name},
	}}}
	_, err := d.collection.UpdateOne(ctx, criteria, doc)
	if err != nil{
		return err
	}
	return nil
}

func (d *UserDatastore) Delete(ctx context.Context, criteria interface{}) error{
	_, err := d.collection.DeleteMany(ctx, criteria)
	if err != nil{
		return err
	}
	return nil
}

func (d *UserDatastore) FindOne(ctx context.Context, criteria interface{}, user *model.User) error{
	err := d.collection.FindOne(ctx, criteria).Decode(user)
	fmt.Println(user)
	if err != nil{
		return err
	}
	return nil
}

func (d *UserDatastore) FindAll(ctx context.Context, criteria interface{}) ([]*model.User, error){
	var users []*model.User
	cur, err := d.collection.Find(ctx, criteria)
	if err != nil{
		return users, err
	}

	defer cur.Close(ctx)
	if cur.Next(ctx) {
		var user *model.User
		err = cur.Decode(&user)
		if err != nil{
			return users, err
		}
		users = append(users, user)
	}

	return users, nil
}