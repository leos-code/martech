package schema

import (
	"fmt"
)

type Schema struct {
	Version int      `json:"version,omitempty"`
	Fields  []*Field `json:"fields,omitempty"`
}

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeEnum    FieldType = "enum"
)

type Field struct {
	Name        string    `json:"name,omitempty"`
	DisplayName string    `json:"display_name,omitempty"` // 前端展示的标签名称
	Type        FieldType `json:"type,omitempty"`
	Range       *Range    `json:"range,omitempty"`
	Enum        []*Enum   `json:"enum,omitempty"`
}

type Range struct {
	Min uint64 `json:"min,omitempty"`
	Max uint64 `json:"max,omitempty"`
}

type Enum struct {
	Value    string  `json:"value,omitempty"`
	Children []*Enum `json:"children,omitempty"`
}

func (s *Schema) Check() error {
	var err error
	for _, f := range s.Fields {
		if err = f.check(); err != nil {
			return err
		}
	}

	return nil
}

func (f *Field) check() error {
	if f.Name == "" {
		return fmt.Errorf("field 'name' empty")
	}

	var err error

	switch f.Type {
	case FieldTypeString:
	case FieldTypeInteger:
		if f.Range != nil {
			if err = f.Range.check(); err != nil {
				return err
			}
		}
	case FieldTypeEnum:
		if len(f.Enum) == 0 {
			return fmt.Errorf("field[%s] no enum", f.Name)
		}
		for _, e := range f.Enum {
			if err := e.check(); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("field[%s] not valid field type[%s]", f.Name, f.Type)
	}

	return nil
}

func (f *Field) GetAllEnum() []string {
	ret := make([]string, 0)
	for _, e := range f.Enum  {
		f.traverseEnum(e, &ret)
	}

	return ret
}

func (f *Field) traverseEnum(e *Enum, ret *[]string) {
	*ret = append(*ret, e.Value)
	for _, child := range e.Children {
		f.traverseEnum(child, ret)
	}
}

func (r *Range) check() error {
	if r.Min > r.Max {
		return fmt.Errorf("range min bigger than max")
	}

	return nil
}

func (e *Enum) check() error {
	if e.Value == "" {
		return fmt.Errorf("enum can not be empty")
	}
	for _, child := range e.Children {
		if err := child.check(); err != nil {
			return err
		}
	}
	return nil
}
