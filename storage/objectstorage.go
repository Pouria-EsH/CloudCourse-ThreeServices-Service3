package storage

import "io"

type ImageStorage interface {
	Upload(imageFile io.ReadSeeker, size int64, key string) (string, error)
}
