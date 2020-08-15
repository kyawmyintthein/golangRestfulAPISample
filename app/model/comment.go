package model

type Comment struct {
	BaseBSONModel
	Content     string `json:"content" bson:"title"`
	CommentedBy uint64 `json:"commented_by" bson:"commented_by"`
}

