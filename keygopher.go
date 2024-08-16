package keygopher

type DB struct {
	Name   string
	Engine Engine
}
type Engine interface {
	Write(key, value string) error
	Get(key string) (string, error)
}
type Config struct {
	Name           string
	IndexingMethod func(string) (Engine, error)
}

func New(cfg *Config) (error, *DB) {
	db := DB{
		Name: cfg.Name,
	}

	e, err := InnitSimpleEngine(db.Name)
	db.Engine = e

	return err, &db
}

func (db *DB) Write(key, value string) error {
	return db.Engine.Write(key, value)
}
func (db *DB) Get(key string) (string, error) {
	return db.Engine.Get(key)
}
