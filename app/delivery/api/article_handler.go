package api

import (
	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/service"
	"net/http"
)

type ArticleHandler struct {
	*BaseHandler
	ArticleService service.ArticleService
}

func ProvideArticleHandler(baseHandler *BaseHandler, articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{BaseHandler: baseHandler, ArticleService: articleService}
}

func (handler *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var createArticleDTO dto.CreateArticleDTO
	err := handler.DecodeAndValidate(r, &createArticleDTO)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}

	article, err := handler.ArticleService.CreateNewArticle(r.Context(), &createArticleDTO)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}

	handler.RenderJSON(r, w, http.StatusCreated, article)
	return
}

func (handler *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	url := handler.URLParam(r, "url")
	article, err := handler.ArticleService.GetArticleByURL(r.Context(), url)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}
	handler.RenderJSON(r, w, http.StatusOK, article)
	return
}
