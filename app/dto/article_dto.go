package dto

import "context"

type CreateArticleDTO struct{
	Title string `json:"title"`
	Content string `json:"content"`
	AuthorID int64 `json:"author_id"`
}

func (dto *CreateArticleDTO) Validate(ctx context.Context) error{
	return nil
}