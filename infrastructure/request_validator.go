package infrastructure

import "context"

type RequestValidator interface {
	Validate(ctx context.Context) error
}
