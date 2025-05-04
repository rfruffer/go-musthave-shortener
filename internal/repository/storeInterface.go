package repository

type StoreRepositoryInterface interface {
	Save(shortID string, originalURL string) error
	Get(shortID string) (string, error)
}
