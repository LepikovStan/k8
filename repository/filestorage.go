package repository

type FileStorage interface {
	Save(key string, data []byte) error
	Get(key string) ([]byte, error)
}
