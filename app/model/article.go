package model

type Article struct {
	BaseBSONModel
	Title    string    `json:"title" json:"title"`
	Url      string    `json:"url" json:"url"`
	Content  string    `json:"content" bson:"title"`
	AuthorID int64     `json:"author_id" bson:"title"`
	Comment  []Comment `json:"comments" bson:"comments"`
}
