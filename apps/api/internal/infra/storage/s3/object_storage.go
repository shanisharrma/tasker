package s3

import (
	"context"
	"io"

	"github.com/shanisharrma/tasker/internal/domain/storage"
	"github.com/shanisharrma/tasker/internal/shared/aws"
)

type ObjectStorageRepository struct {
	s3 *aws.S3Client
}

func NewObjectStorageRepository(s3Client *aws.S3Client) storage.ObjectStorage {
	return &ObjectStorageRepository{
		s3: s3Client,
	}
}

func (r *ObjectStorageRepository) Upload(
	ctx context.Context,
	bucket string,
	fileName string,
	body io.Reader,
) (string, error) {
	return r.s3.UploadFile(ctx, bucket, fileName, body)
}

func (r *ObjectStorageRepository) Delete(
	ctx context.Context,
	bucket string,
	key string,
) error {
	return r.s3.DeleteObject(ctx, bucket, key)
}

func (r *ObjectStorageRepository) Presign(
	ctx context.Context,
	bucket string,
	key string,
) (string, error) {
	return r.s3.CreatePresignedUrl(ctx, bucket, key)
}
