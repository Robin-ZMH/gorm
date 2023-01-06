package clause

import (
	"strings"
)

type Type int

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(name Type, vars ...any) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]any)
	}
	sqlFmt, sqlVars := generators[name](vars...)
	c.sql[name] = sqlFmt
	c.sqlVars[name] = sqlVars
}

func (c *Clause) Build(keywords ...Type) (sqlFmt string, vals []any) {
	sqls := []string{}
	for _, key := range keywords {
		if _, ok := c.sql[key]; ok {
			sqls = append(sqls, c.sql[key])
			vals = append(vals, c.sqlVars[key]...)
		}
	}
	sqlFmt = strings.Join(sqls, " ")
	return
}
