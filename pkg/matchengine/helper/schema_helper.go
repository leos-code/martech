package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"reflect"

	"github.com/tencentad/martech/api/proto/retrieval"
	"github.com/tencentad/martech/pkg/schema"
)

type Helper struct {
	fieldHelper map[string]*FieldHelper
	Version     int
}

func NewSchemaHelper(filepath string) (*Helper, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return NewSchemaHelperFromContent(content)
}

type FieldHelper struct {
	field       *schema.Field
	enumParents map[string][]string // 获取枚举值的父对象
}

func NewFieldHelper(field *schema.Field) (*FieldHelper, error) {
	fh := &FieldHelper{
		field: field,
	}

	if field.Type == schema.FieldTypeEnum {
		if err := fh.traverseEnumField(); err != nil {
			return nil, err
		}
	}

	return fh, nil
}

func NewSchemaHelperFromContent(content []byte) (*Helper, error) {
	h := &Helper{
		fieldHelper: make(map[string]*FieldHelper),
	}

	var s schema.Schema
	var err error
	if err = json.Unmarshal(content, &s); err != nil {
		return nil, err
	}

	if err = s.Check(); err != nil {
		return nil, err
	}

	h.Version = s.Version
	for _, f := range s.Fields {
		if err = h.traverseField(f); err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (h *Helper) traverseField(field *schema.Field) error {
	fh, err := NewFieldHelper(field)
	if err != nil {
		return err
	}
	h.fieldHelper[field.Name] = fh
	return nil
}

func (h *FieldHelper) traverseEnumField() error {
	field := h.field
	enumParents := make(map[string][]string)
	for _, e := range field.Enum {
		h.traverseEnum(e, nil, enumParents)
	}

	h.enumParents = enumParents
	return nil
}

func (h *FieldHelper) traverseEnum(e *schema.Enum, parents []string, m map[string][]string) {
	if len(parents) != 0 {
		m[e.Value] = parents
	}

	newParents := append(parents, e.Value)
	for _, child := range e.Children {
		h.traverseEnum(child, newParents, m)
	}
}

func (h *FieldHelper) getEnumParents(enum string) []string {
	return h.enumParents[enum]
}

func (h *Helper) GetFeatureValues(fieldName string, values ...interface{}) ([]*retrieval.Feature_Value, error) {
	fh, ok := h.fieldHelper[fieldName]
	if !ok {
		return nil, fmt.Errorf("field[%s] not in schema config", fieldName)
	}

	switch fh.field.Type {
	case schema.FieldTypeEnum:
		return h.getEnumFeatureValues(fh, values...)
	case schema.FieldTypeString:
		return h.getStringFeatureValues(fh.field, values...)
	case schema.FieldTypeInteger:
		return h.getIntegerFeatureValues(fh.field, values...)
	default:
		return nil, fmt.Errorf("filed[%s] not valid type", fh.field.Name)
	}
}

func (h *Helper) getEnumFeatureValues(fh *FieldHelper, values ...interface{}) ([]*retrieval.Feature_Value, error) {
	f := fh.field
	valueSet := make(map[string]struct{})
	enumParents := fh.enumParents
	for _, v := range values {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("field[%s] enum type, but value[%v] is not string", f.Name, v)
		}
		valueSet[str] = struct{}{}
		for _, p := range enumParents[str] {
			valueSet[p] = struct{}{}
		}
	}
	featureValues := make([]*retrieval.Feature_Value, 0, len(valueSet))
	for v := range valueSet {
		featureValues = append(featureValues, string2FeatureValue(v))
	}
	return featureValues, nil
}

func (h *Helper) getStringFeatureValues(f *schema.Field, values ...interface{}) ([]*retrieval.Feature_Value, error) {
	featureValues := make([]*retrieval.Feature_Value, 0, len(values))
	for _, v := range values {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("field[%s] string type, but value[%v] is not string", f.Name, v)
		}
		featureValues = append(featureValues, string2FeatureValue(str))
	}
	return featureValues, nil
}

func (h *Helper) getIntegerFeatureValues(f *schema.Field, values ...interface{}) ([]*retrieval.Feature_Value, error) {
	featureValues := make([]*retrieval.Feature_Value, 0, len(values))
	for _, v := range values {
		id, err := numberToValueID(f, v)
		if err != nil {
			return nil, err
		}
		featureValues = append(featureValues, id2FeatureValue(id))
	}
	return featureValues, nil
}

func string2FeatureValue(str string) *retrieval.Feature_Value {
	return &retrieval.Feature_Value{
		Type: retrieval.Feature_Value_String,
		Str:  str,
	}
}

func id2FeatureValue(id uint64) *retrieval.Feature_Value {
	return &retrieval.Feature_Value{
		Type: retrieval.Feature_Value_ID,
		Id:   id,
	}
}

func numberToValueID(f *schema.Field, value interface{}) (uint64, error) {
	v := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		v := v.Int()
		if v < 0 {
			return 0, fmt.Errorf("field[%s] value[%d] below zero", f.Name, value)
		}
		return uint64(v), nil
	case uint, uint8, uint16, uint32, uint64:
		return v.Uint(), nil
	case float32, float64:
		v := v.Float()
		if v < 0 {
			return 0, fmt.Errorf("field[%s] value[%d] below zero", f.Name, value)
		}
		if math.Floor(v) != v {
			return 0, fmt.Errorf("field[%s] value[%d] has decimal part", f.Name, value)
		}
		return uint64(v), nil
	default:
		return 0, fmt.Errorf("field[%s] value[%d] not valid integer", f.Name, value)
	}
}
