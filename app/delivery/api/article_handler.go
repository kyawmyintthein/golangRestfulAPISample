package api

import "github.com/kyawmyintthein/golangRestfulAPISample/app/service"

type ArticleHandler struct{
	ArticleService service.ArticleService
}

func ProvideArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{ArticleService: articleService}
}

