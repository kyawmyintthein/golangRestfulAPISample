package model

type BaseModel struct {
	ID        int64 `json:"id" db:"id"`
	CreatedAt int64 `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt int64 `json:"updated_at,omitempty" db:"updated_at""`
}
