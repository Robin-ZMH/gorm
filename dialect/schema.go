package dialect

import "reflect"

// Field represents a column of database
type Field struct {
	Name string
	Type string
	Tag  string
}

// Table represents a table of database
type Table struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field // [FieldName]Field
}

func (t *Table) GetField(name string) *Field {
	return t.fieldMap[name]
}

func Parse(dest any, dialect Dialect) *Table {
	modelTyp := reflect.Indirect(reflect.ValueOf(dest)).Type()
	table := &Table{
		Model:    dest,
		Name:     modelTyp.Name(),
		fieldMap: make(map[string]*Field),
	}
	for i := 0; i < modelTyp.NumField(); i++ {
		f := modelTyp.Field(i)
		if f.Anonymous || !f.IsExported() {
			continue
		}
		field := &Field{
			Name: f.Name,
			Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
		}
		if tag, ok := f.Tag.Lookup("db"); ok {
			field.Tag = tag
		}
		table.Fields = append(table.Fields, field)
		table.FieldNames = append(table.FieldNames, field.Name)
		table.fieldMap[field.Name] = field
	}
	return table
}
