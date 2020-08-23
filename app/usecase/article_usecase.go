package usecase

import (
	"context"

	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/errors"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
	"github.com/kyawmyintthein/golangRestfulAPISample/infrastructure"
)

type ArticleUsecase interface {
	CreateNewArticle(context.Context, *dto.CreateArticleDTO) (*model.Article, error)
	GetArticleByURL(context.Context, string) (*model.Article, error)
}

type articleUsecase struct {
	stringHelper      infrastructure.StringHelper
	userRepository    repository.UserRepository
	articleRepository repository.ArticlesRepository
}

func ProvideArticleUsecase(userRepository repository.UserRepository,
	articleRepository repository.ArticlesRepository,
	stringHelper infrastructure.StringHelper) ArticleUsecase {
	return &articleUsecase{
		stringHelper:      stringHelper,
		userRepository:    userRepository,
		articleRepository: articleRepository,
	}
}

func (service *articleUsecase) CreateNewArticle(ctx context.Context, createArticleDTO *dto.CreateArticleDTO) (*model.Article, error) {
	newArticle := model.Article{
		Title:    createArticleDTO.Title,
		Url:      service.stringHelper.StringToURL(createArticleDTO.Title),
		Content:  createArticleDTO.Content,
		AuthorID: createArticleDTO.AuthorID,
	}
	article, err := service.articleRepository.Create(ctx, &newArticle)
	if err != nil {
		if infrastructure.IsDuplicateError(err) {
			return article, errors.NewDuplicateResourceError("article").Wrap(err)
		}
		return article, errors.NewUnknownError().Wrap(err)
	}
	return article, err
}

func (service *articleUsecase) GetArticleByURL(ctx context.Context, url string) (*model.Article, error) {
	article, err := service.articleRepository.GetByURL(ctx, url)
	if err != nil {
		return article, errors.NewUnknownError().Wrap(err)
	}
	return article, err
}
