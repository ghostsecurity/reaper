package storage

type Store interface {
	Save(entry *Entry) error
	Get(id int64) (*Entry, error)
	List(limit, offset int) ([]*Entry, error)
	Search(params SearchParams) ([]*Entry, error)
	Close() error
}
