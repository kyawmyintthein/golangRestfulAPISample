package api

import (
	"net/http"

	"github.com/kyawmyintthein/golangRestfulAPISample/app/dto"
	"github.com/kyawmyintthein/golangRestfulAPISample/app/usecase"
)

type ArticleHandler struct {
	*BaseHandler
	ArticleUsecase usecase.ArticleUsecase
}

func ProvideArticleHandler(baseHandler *BaseHandler, articleUsecase usecase.ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{BaseHandler: baseHandler, ArticleUsecase: articleUsecase}
}

func (handler *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var createArticleDTO dto.CreateArticleDTO
	err := handler.DecodeAndValidate(r, &createArticleDTO)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}

	article, err := handler.ArticleUsecase.CreateNewArticle(r.Context(), &createArticleDTO)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}

	handler.RenderJSON(r, w, http.StatusCreated, article)
	return
}

func (handler *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	url := handler.URLParam(r, "url")
	article, err := handler.ArticleUsecase.GetArticleByURL(r.Context(), url)
	if err != nil {
		handler.RenderErrorAsJSON(r, w, err)
		return
	}
	handler.RenderJSON(r, w, http.StatusOK, article)
	return
}
