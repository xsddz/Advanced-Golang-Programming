package storage

type FileCache struct{}

func NewFileCache(path string) (*FileCache, error) {

	return &FileCache{}, nil
}
