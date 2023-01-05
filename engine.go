package gorm

import (
	"database/sql"

	"github.com/Robin-ZMH/gorm/log"
)

type Engine struct {
	db *sql.DB
}


// NewEngine create a instance of Engine
// connect database and ping it to test whether it's alive
func NewEngine(driver, dbName string) (e *Engine, err error) {
	db, err := sql.Open(driver, dbName)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

// Close database connection
func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// NewSession creates a new session for next operations
func (engine *Engine) NewSession() *Session {
	return NewSession(engine.db)
}

