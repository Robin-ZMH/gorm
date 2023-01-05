package gorm

import (
	"database/sql"
	"github.com/Robin-ZMH/gorm/log"
	"strings"
)

type Session struct {
	db   *sql.DB
	sql  strings.Builder
	vars []any
}

func NewSession(db *sql.DB) *Session {
	return &Session{db: db}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.vars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw appends sql and sqlVars
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.vars = append(s.vars, values...)
	return s
}

// Exec raw sql with vars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.vars)
	if result, err = s.db.Exec(s.sql.String(), s.vars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets one record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.vars)
	return s.db.QueryRow(s.sql.String(), s.vars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.vars)
	if rows, err = s.db.Query(s.sql.String(), s.vars...); err != nil {
		log.Error(err)
	}
	return
}