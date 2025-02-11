package storage

type Storage []StorageEntry

type StorageEntry struct {
	ShortID string `json:"short_id"`
	FullURL string `json:"full_url"`
}

type StorageWriter interface {
	WriteEntry(entry *StorageEntry) error
	WriteEntries(entries []StorageEntry) error
	Close() error
}

type StorageReader interface {
	ReadEntry() (*StorageEntry, error)
	ReadAllEntries() (Storage, error)
	Close() error
}
