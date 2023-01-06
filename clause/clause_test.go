package clause

import (
	"reflect"
	"testing"
)

func TestSet(t *testing.T) {
	c := &Clause{}
	c.Set(LIMIT, 3)
	c.Set(SELECT, "User", []string{"*"})
	c.Set(WHERE, "Name = ?", "Tom")
	c.Set(ORDERBY, "Age ASC")
	sql, vars := c.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Error("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Error("failed to build SQLVars")
	}
}

func Test_insert(t *testing.T) {
	c := &Clause{}
	c.Set(INSERT, "User", []string{"Name", "Age"})
	sql := c.sql[INSERT]
	vars := c.sqlVars[INSERT]
	t.Log(sql, vars)
	if sql != "INSERT INTO User (Name,Age)" || len(vars) != 0 {
		t.Fatal("failed to get clause")
	}
}

// func Test_values(t *testing.T) {
// 	c := &Clause{}
// 	c.Set(VALUES, "User", "Robin", 25)
// }