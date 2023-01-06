package gorm

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id   int `db:"PRIMARY KEY"`
	Name string
	Age  int
}

var (
	user1 = &User{1, "Tom", 18}
	user2 = &User{2, "Sam", 25}
	user3 = &User{3, "Jack", 25}
)

func TestCreateTable(t *testing.T) {
	e, _ := NewEngine("sqlite3", "gee.db")
	defer e.Close()

	s := e.NewSession().Model(&User{})
	err := s.CreateTable()
	if err != nil {
		t.Fatal("failed create table")
	}
}

func TestInsert(t *testing.T) {
	e, _ := NewEngine("sqlite3", "gee.db")
	defer e.Close()

	RowsAffected, err := e.NewSession().Insert(user1, user2, user3)
	if err != nil {
		t.Fatal("failed insert test")
	}
	t.Log(RowsAffected)
}

func TestAll(t *testing.T) {
	e, _ := NewEngine("sqlite3", "gee.db")
	defer e.Close()

	users := []User{}
	if err := e.NewSession().All(&users); err != nil {
		t.Fatalf("failed select all data test, err is %v", err)
	}
	t.Log(users)
}

func TestDropTabe(t *testing.T) {
	e, _ := NewEngine("sqlite3", "gee.db")
	defer e.Close()
	err := e.NewSession().Model(&User{}).DropTable()
	if err != nil {
		t.Fatal("failed drop table test")
	}
}
