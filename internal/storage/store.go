package storage

type Store interface {
	Save(entry *Entry) error
	Get(id int64) (*Entry, error)
	List(limit, offset int) ([]*Entry, error)
	ListAfter(afterID int64, limit int) ([]*Entry, error)
	Search(params SearchParams) ([]*Entry, error)
	Clear() error
	Close() error
}
