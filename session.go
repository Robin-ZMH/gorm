package gorm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/Robin-ZMH/gorm/clause"
	"github.com/Robin-ZMH/gorm/dialect"
	"github.com/Robin-ZMH/gorm/log"
	"github.com/Robin-ZMH/gorm/table"
)

type Session struct {
	db      *sql.DB
	dialect dialect.Dialect
	Table   *table.Table
	clause  *clause.Clause
	sql     strings.Builder
	sqlVars []any
}

func NewSession(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
		clause:  &clause.Clause{},
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

// Raw appends sql and sqlVars
func (s *Session) Raw(sql string, vars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, vars...)
	return s
}

// Exec raw sql with vars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets one record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Model(obj interface{}) *Session {
	s.Table = table.NewTable(obj, s.dialect)
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

func (s *Session) Insert(objs ...any) (int64, error) {
	fieldVals := []any{}
	for _, obj := range objs {
		table := s.Model(obj).Table
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		fieldVals = append(fieldVals, table.FieldValues())
	}
	s.clause.Set(clause.VALUES, fieldVals...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) All(objs any) error {
	refSlice := reflect.Indirect(reflect.ValueOf(objs))
	objType := refSlice.Type().Elem()
	s.Model(reflect.New(objType).Elem().Interface())

	s.clause.Set(clause.SELECT, s.Table.Name, s.Table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT)
	res, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for res.Next() {
		fields := []any{}
		obj := reflect.New(objType).Elem()
		// obtain the address of fields
		for _, name := range s.Table.FieldNames {
			fields = append(fields, obj.FieldByName(name).Addr().Interface())
		}
		// fill in the fields with query data
		if err := res.Scan(fields...); err != nil {
			return err
		}
		refSlice.Set(reflect.Append(refSlice, obj))
	}
	return res.Close()
}
