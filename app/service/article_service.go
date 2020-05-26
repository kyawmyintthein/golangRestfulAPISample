package service

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/repository"
)

type ArticleService interface{
	CreateNewArticle(context.Context, int64, *dto.CreateArticleDTO) (*model.Article, error)
}

type articleService struct {
	userRepository repository.UserRepository
	articleRepository repository.ArticlesRepository
}

func ProvideArticleService(userRepository repository.UserRepository, articleRepository repository.ArticlesRepository) ArticleService{
	return &articleService{
		userRepository: userRepository,
		articleRepository: articleRepository,
	}
}

func (articleService *articleService) CreateNewArticle(ctx context.Context, authorID int64, createArticleDTO *dto.CreateArticleDTO) (*model.Article, error){
	newArticle := model.Article{
		Title:         createArticleDTO.Title,
		Content:       createArticleDTO.Content,
		AuthorID:      authorID,
	}
	article, err := articleService.articleRepository.Create(ctx, &newArticle)
	if err != nil{
		return article, err
	}
	return article, err
}
