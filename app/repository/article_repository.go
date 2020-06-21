package repository

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/model"
)

type ArticlesRepository interface {
	Create(context.Context, *model.Article) (*model.Article, error)
	GetByURL(context.Context, string) (*model.Article, error)
}

