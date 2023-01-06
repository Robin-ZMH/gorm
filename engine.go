package gorm

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/Robin-ZMH/gorm/clause"
	"github.com/Robin-ZMH/gorm/dialect"
	"github.com/Robin-ZMH/gorm/log"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
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

	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dialect}
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
func (e *Engine) NewSession() *Session {
	return NewSession(e.db, e.dialect)
}

func Test() {
	var c clause.Clause
	c.Set(clause.LIMIT, 3)
	c.Set(clause.SELECT, "User", []string{"*"})
	c.Set(clause.WHERE, "Name = ?", "Tom")
	c.Set(clause.ORDERBY, "Age ASC")
	sql, vars := c.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	log.Info(fmt.Sprintf(sql, vars))
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		log.Error("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		log.Error("failed to build SQLVars")
	}
}
