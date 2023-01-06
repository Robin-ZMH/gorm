package gorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Robin-ZMH/gorm/dialect"
	"github.com/Robin-ZMH/gorm/log"
	"github.com/Robin-ZMH/gorm/table"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	vars    []any
	dialect dialect.Dialect
	Table   *table.Table
}

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
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

func (s *Session) Model(model interface{}) *Session {
	// nil or different model, update table
	if s.Table == nil || reflect.TypeOf(model) != reflect.TypeOf(s.Table.Model) {
		s.Table = table.NewTable(model, s.dialect)
	}
	return s
}

func (s *Session) CreateTable() error {
	if s.Table == nil {
		log.Error("Model is not set")
	}
	var columns []string
	for _, filed := range s.Table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s",
			filed.Name,
			filed.Type,
			filed.Tag))
	}

	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s(%s);", s.Table.Name, desc)).Exec()
	return err
}

func (s *Session) DropTable() error {
	if s.Table == nil {
		log.Error("Model is not set")
	}
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.Table.Name)).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, args := s.dialect.TableExistSQL(s.Table.Name)
	row := s.Raw(sql, args...).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.Table.Name
}
