package storage

// Storage записи
type Storage []StorageEntry

// StorageEntry запись
type StorageEntry struct {
	ShortID string `json:"short_id"`
	FullURL string `json:"full_url"`
}

// StorageWriter интерфейс для записи информации в файл
type StorageWriter interface {
	WriteEntry(entry *StorageEntry) error
	WriteEntries(entries []StorageEntry) error
	Close() error
}

// StorageReader интерфейс для чтения записей из файла
type StorageReader interface {
	ReadEntry() (*StorageEntry, error)
	ReadAllEntries() (Storage, error)
	Close() error
}
