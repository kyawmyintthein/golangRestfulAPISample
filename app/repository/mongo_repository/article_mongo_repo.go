package mongo_repository

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	base_repository "github.com/kyawmyintthein/golangRestfulAPISample/internal/base-repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type articleMongoRepo struct{
	*base_repository.BaseMongoRepo
	collection string
}

func ProvideArticleRepository(baseMongoRepo *base_repository.BaseMongoRepo) repository.ArticlesRepository{
	return &articleMongoRepo{
		baseMongoRepo,
		"articles",
	}
}


func (repo *articleMongoRepo) Create(ctx context.Context, article *model.Article) (*model.Article, error){
	collection := repo.MongodbConnector.DB(ctx).Collection(repo.collection)
	result, err := collection.InsertOne(ctx, article)
	if err != nil{
		return article, err
	}

	article.RawID, _ = result.InsertedID.(primitive.ObjectID)
	return article, nil
}


func (repo *articleMongoRepo) GetByURL(ctx context.Context, url string) (*model.Article, error){
	var (
		err error
		article model.Article
	)

	collection := repo.MongodbConnector.DB(ctx).Collection(repo.collection)
	err = collection.FindOne(ctx, bson.M{"url": url}).Decode(&article)
	if err != nil{
		return &article, err
	}

	return &article, nil
}