package storage

import (
	"context"
	"io"
)

type ObjectStorage interface {
	Upload(
		ctx context.Context,
		bucket string,
		fileName string,
		body io.Reader,
	) (string, error)

	Delete(
		ctx context.Context,
		bucket string,
		key string,
	) error

	Presign(
		ctx context.Context,
		bucket string,
		key string,
	) (string, error)
}
