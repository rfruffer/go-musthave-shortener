package repository

type StoreRepositoryInterface interface {
	Save(shortID string, originalURL string, uuid string) error
	Get(shortID string) (string, error)

	SaveToFile(path string) error
	LoadFromFile(path string) error

	Ping() error
}
