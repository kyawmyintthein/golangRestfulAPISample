package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/api_errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/mongo"
)

type ArticleService interface {
	CreateNewArticle(context.Context, *dto.CreateArticleDTO) (*model.Article, error)
	GetArticleByURL(context.Context, string) (*model.Article, error)
}

type articleService struct {
	stringHelper      infrastructure.StringHelper
	userRepository    repository.UserRepository
	articleRepository repository.ArticlesRepository
}

func ProvideArticleService(userRepository repository.UserRepository,
	articleRepository repository.ArticlesRepository,
	stringHelper infrastructure.StringHelper) ArticleService {
	return &articleService{
		stringHelper:      stringHelper,
		userRepository:    userRepository,
		articleRepository: articleRepository,
	}
}

func (service *articleService) CreateNewArticle(ctx context.Context, createArticleDTO *dto.CreateArticleDTO) (*model.Article, error) {
	newArticle := model.Article{
		Title:    createArticleDTO.Title,
		Url:      service.stringHelper.StringToURL(createArticleDTO.Title),
		Content:  createArticleDTO.Content,
		AuthorID: createArticleDTO.AuthorID,
	}
	article, err := service.articleRepository.Create(ctx, &newArticle)
	if err != nil {
		if mongo.IsDuplicateError(err){
			return article, api_errors.NewDuplicateResourceError("article").Wrap(err)
		}
		return article, api_errors.NewUnknownError().Wrap(err)
	}
	return article, err
}


func (service *articleService) GetArticleByURL(ctx context.Context, url string) (*model.Article, error) {
	article, err := service.articleRepository.GetByURL(ctx, url)
	if err != nil {
		return article, api_errors.NewUnknownError().Wrap(err)
	}
	return article, err
}
