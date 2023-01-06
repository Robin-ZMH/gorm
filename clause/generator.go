package clause

import (
	"fmt"
	"strings"
)

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = insert
	generators[VALUES] = values
	generators[SELECT] = _select
	generators[LIMIT] = limit
	generators[WHERE] = where
	generators[ORDERBY] = orderBy
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type generator func(vals ...any) (sqlFmt string, vars []any)

var generators map[Type]generator

func genBindVars(num int) string {
	var varSlice []string
	for i := 0; i < num; i++ {
		varSlice = append(varSlice, "?")
	}
	return strings.Join(varSlice, ", ")
}

func insert(vals ...any) (sqlFmt string, vars []any) {
	table := vals[0]
	fields := strings.Join(vals[1].([]string), ",")
	sqlFmt = fmt.Sprintf("INSERT INTO %s (%v)", table, fields)
	return
}

func values(vals ...any) (sqlFmt string, vars []any) {
	var bindStr string
	var sqls []string
	for _, v := range vals {
		val := v.([]any)
		if bindStr == "" {
			bindStr = genBindVars(len(val))
		}
		sqls = append(sqls, fmt.Sprintf("(%v)", bindStr))
		vars = append(vars, val...)
	}
	sqlFmt = strings.Join(sqls, ", ")
	return
}

func _select(vals ...any) (sqlFmt string, vars []any) {
	fields := strings.Join(vals[1].([]string), ",")
	table := vals[0].(string)
	sqlFmt = fmt.Sprintf("SELECT %s FROM %s", fields, table)
	return
}

func limit(vals ...interface{}) (sqlFmt string, vars []any) {
	return "LIMIT ?", vals
}

func where(values ...interface{}) (sqlFmt string, vars []any) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	sqlFmt = fmt.Sprintf("WHERE %s", desc)
	return
}

func orderBy(values ...interface{}) (sqlFmt string, vars []any) {
	return fmt.Sprintf("ORDER BY %s", values[0]), vars
}
