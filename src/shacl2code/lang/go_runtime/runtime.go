package v3_0

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"reflect"
)

// LDType is a special data holder property type for type-level linked data
type LDType struct{}

type LDContext map[string]*serializationContext

func (c LDContext) RegisterTypes(contextUrl string, types ...any) LDContext {
	ctx := c[contextUrl]
	if ctx == nil {
		ctx = &serializationContext{
			contextUrl: contextUrl,
			iriToType:  map[string]*typeContext{},
		}
		c[contextUrl] = ctx
	}
	for _, typ := range types {
		ctx.registerType(typ)
	}
	return c
}

func (c LDContext) ToJSON(writer io.Writer, value any) error {
	panic("not implemented")
}

func (c LDContext) FromJSON(reader io.Reader) ([]any, error) {
	vals := map[string]any{}
	dec := json.NewDecoder(reader)
	err := dec.Decode(&vals)
	if err != nil {
		return nil, err
	}
	return c.graphFromMaps(vals)
}

func (c LDContext) FromMaps(values map[string]any) (any, error) {
	return c.graphFromMaps(values)
}

func (c LDContext) ToMaps(o ...any) (values map[string]any, errors error) {
	panic("not implemented")
}

func (c LDContext) graphFromMaps(values map[string]any) ([]any, error) {
	instances := map[string]reflect.Value{}

	var errs error
	var graph []any

	context, _ := values["@context"].(string)
	currentContext := c[context]
	if currentContext == nil {
		return nil, fmt.Errorf("unknown document @context type: %v", context)
	}

	nodes, _ := values["@graph"].([]any)
	if nodes == nil {
		return nil, fmt.Errorf("@graph array not present in root object")
	}

	// one pass to create all the instances
	for _, node := range nodes {
		_, err := c.getOrCreateInstance(currentContext, instances, node)
		errs = appendErr(errs, err)
	}

	// second pass to fill in all refs
	for _, node := range nodes {
		got, err := c.getOrCreateInstance(currentContext, instances, node)
		errs = appendErr(errs, err)
		if err == nil {
			graph = append(graph, got.Interface())
		}
	}

	return graph, errs
}

func (c LDContext) getOrCreateInstance(currentContext *serializationContext, instances map[string]reflect.Value, incoming any) (reflect.Value, error) {
	switch incoming := incoming.(type) {
	case string:
		inst := c.findById(currentContext, instances, incoming)
		return inst, nil
	case map[string]any:
		return c.getOrCreateFromMap(currentContext, instances, incoming)
	}
	return emptyValue, fmt.Errorf("unexpected data type: %#v", incoming)
}

func (c LDContext) findById(_ *serializationContext, instances map[string]reflect.Value, incoming string) reflect.Value {
	inst, ok := instances[incoming]
	if ok {
		return inst
	}
	return emptyValue
}

func (c LDContext) getOrCreateFromMap(currentContext *serializationContext, instances map[string]reflect.Value, incoming map[string]any) (reflect.Value, error) {
	typ, ok := incoming["type"].(string)
	if !ok {
		return emptyValue, fmt.Errorf("not a string")
	}

	t, ok := currentContext.iriToType[typ]
	if !ok {
		return emptyValue, fmt.Errorf("don't have type: %v", typ)
	}

	id, _ := incoming["@id"].(string)
	if id == "" {
		// FIXME inspect type for field with @id to get the right name
		id, _ = incoming["spdxId"].(string)
	}
	inst, ok := instances[id]
	if !ok {
		inst = reflect.New(baseType(t.typ)) // New(T) returns *T
		instances[id] = inst
	}

	// valid type, make a new one and fill it from the incoming maps
	return inst, c.fill(currentContext, instances, inst, incoming)
}

func (c LDContext) fill(currentContext *serializationContext, instances map[string]reflect.Value, instance reflect.Value, incoming any) error {
	switch incoming := incoming.(type) {
	case string:
		inst := c.findById(currentContext, instances, incoming)
		if inst != emptyValue {
			return c.setValue(currentContext, instances, instance, inst)
		}
		return nil
	case map[string]any:
		return c.setStructProps(currentContext, instances, instance, incoming)
	}
	return fmt.Errorf("unsupported incoming data type: %#v attempting to set instance: %#v", incoming, instance.Interface())
}

func (c LDContext) setValue(currentContext *serializationContext, instances map[string]reflect.Value, target reflect.Value, incoming any) error {
	var errs error
	typ := target.Type()
	switch typ.Kind() {
	case reflect.Slice:
		switch incoming := incoming.(type) {
		case []any:
			return c.setSliceValue(currentContext, instances, target, incoming)
		}
		return fmt.Errorf("expected slice for %#v, got %#v", target, incoming)
	case reflect.Struct:
		switch incoming := incoming.(type) {
		case map[string]any:
			return c.setStructProps(currentContext, instances, target, incoming)
		}
	case reflect.Interface, reflect.Pointer:
		switch incoming := incoming.(type) {
		case string, map[string]any:
			inst, err := c.getOrCreateInstance(currentContext, instances, incoming)
			errs = appendErr(errs, err)
			if inst != emptyValue {
				target.Set(inst)
				return nil
			}
		}
	default:
		incomingValue := reflect.ValueOf(incoming)
		if incomingValue.CanConvert(typ) {
			newVal := incomingValue.Convert(typ)
			target.Set(newVal)
		} else {
			errs = appendErr(errs, fmt.Errorf("unable to convert %#v to %s, dropping", incoming, typeName(typ)))
		}
	}
	return nil
}

func (c LDContext) setStructProps(currentContext *serializationContext, instances map[string]reflect.Value, instance reflect.Value, incoming map[string]any) error {
	var errs error
	typ := instance.Type()
	for typ.Kind() == reflect.Pointer {
		instance = instance.Elem()
		typ = instance.Type()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("unable to set struct properties on non-struct type: %#v", instance.Interface())
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if skipField(field) {
			continue
		}
		fieldVal := instance.Field(i)

		propName := field.Tag.Get(propIriCompactTagName)
		if propName == "" {
			propName = field.Tag.Get(propIriTagName)
		}
		if propName != "" {
			incomingVal, ok := incoming[propName]
			if ok {
				errs = appendErr(errs, c.setValue(currentContext, instances, fieldVal, incomingVal))
			}
		}
	}
	return errs
}

func (c LDContext) setSliceValue(currentContext *serializationContext, instances map[string]reflect.Value, target reflect.Value, incoming []any) error {
	var errs error
	sliceType := target.Type()
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("expected slice, got: %#v", target)
	}
	sz := len(incoming)
	if sz > 0 {
		elemType := sliceType.Elem()
		newSlice := reflect.MakeSlice(sliceType, 0, sz)
		for i := 0; i < sz; i++ {
			incomingValue := incoming[i]
			if incomingValue == nil {
				continue // don't allow null values
			}
			newItemValue, err := c.getOrCreateInstance(currentContext, instances, incomingValue)
			errs = appendErr(errs, err)
			if newItemValue != emptyValue {
				// validate we can actually set the type
				if newItemValue.Type().AssignableTo(elemType) {
					newSlice = reflect.Append(newSlice, newItemValue)
				}
			}
		}
		target.Set(newSlice)
	}
	return errs
}

func skipField(field reflect.StructField) bool {
	return field.Type.Size() == 0
}

func typeName(t reflect.Type) string {
	switch {
	case isPointer(t):
		return "*" + typeName(t.Elem())
	case isSlice(t):
		return "[]" + typeName(t.Elem())
	case isMap(t):
		return "map[" + typeName(t.Key()) + "]" + typeName(t.Elem())
	case isPrimitive(t):
		return t.Name()
	}
	return path.Base(t.PkgPath()) + "." + t.Name()
}

func isSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

func isMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map
}

func isPointer(t reflect.Type) bool {
	return t.Kind() == reflect.Pointer
}

func isPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Bool:
		return true
	default:
		return false
	}
}

const (
	typeIriTagName        = "iri"
	typeIriCompactTagName = "iri-compact"
	propIriTagName        = "iri"
	propIriCompactTagName = "iri-compact"
	typeIdPropTagName     = "id-prop"
)

var emptyValue reflect.Value
var errNotFound = fmt.Errorf("not found")

type typeContext struct {
	typ     reflect.Type
	iri     string
	compact string
	idProp  string
}

type serializationContext struct {
	contextUrl string
	iriToType  map[string]*typeContext
}

func fieldByType[T any](t reflect.Type) (reflect.StructField, bool) {
	var v T
	typ := reflect.TypeOf(v)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Type == typ {
			return f, true
		}
	}
	return reflect.StructField{}, false
}

func (m *serializationContext) registerType(instancePointer any) {
	t := reflect.TypeOf(instancePointer)
	t = baseType(t) // types may be passed as pointers; we want the base types
	tc := &typeContext{
		typ: t,
	}
	meta, ok := fieldByType[LDType](t)
	if ok {
		tc.iri = meta.Tag.Get(typeIriTagName)
		tc.compact = meta.Tag.Get(typeIriCompactTagName)
		tc.idProp = meta.Tag.Get(typeIdPropTagName)
	}
	m.iriToType[tc.iri] = tc
	m.iriToType[tc.compact] = tc
	//m.typeToType[tc.typ] = tc
}

func isNil(from reflect.Value) bool {
	return from.IsZero()
}

func appendErr(err error, errs ...error) error {
	if joined, ok := err.(interface{ Unwrap() []error }); ok {
		return errors.Join(append(joined.Unwrap(), errs...)...)
	}
	if err == nil {
		return errors.Join(errs...)
	}
	return errors.Join(append([]error{err}, errs...)...)
}

func baseType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Pointer:
		return baseType(t.Elem())
	default:
		return t
	}
}

//var skipTypes = []reflect.Type{
//	reflect.TypeOf(LDType{}),
//	//reflect.TypeOf(ShaclIri("")),
//}
//
//func skipField(tc *typeContext, f reflect.StructField) bool {
//	if tc.idProp == f.Tag.Get(propIriCompactTagName) || tc.idProp == f.Name {
//		return false
//	}
//	return f.Name == "_" ||
//		f.Tag.Get(propIriTagName) == "" ||
//		slices.Contains(skipTypes, f.Type)
//}
