package clause

import (
	"fmt"
	"strings"
)

type Type int

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(name Type, vars ...any) {
	fmt.Println(vars)
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]any)
	}
	sqlFmt, vars := generators[name](vars)
	c.sql[name] = sqlFmt
	c.sqlVars[name] = vars
}

func (c *Clause) Build(keywords ...Type) (sqlFmt string, vals []any) {
	var sqls []string
	for _, key := range keywords {
		if _, ok := c.sql[key]; ok {
			sqls = append(sqls, c.sql[key])
			vals = append(vals, c.sqlVars[key]...)
		}
	}
	sqlFmt = strings.Join(sqls, " ")
	return
}
