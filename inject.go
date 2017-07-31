// Package inject is the improved version of https://github.com/facebookgo/inject.
//
// Improvements are listed as below:
//
//  1. More intuitive declaration about dependency relations.
//  2. Support constructor.
//  3. Return the initialized dependency graph (retrieve object outside self scope).
//
// The usage example is in the `inject_test.go`, run `go test` and the output should be like as below
//  [:~ … /inject] $ go test
//  2017/07/31 22:36:21 Tillage 3 hours
//  2017/07/31 22:36:21 Got rice
//  2017/07/31 22:36:21 Boeing787 Fly from C++ to Go
//  2017/07/31 22:36:21 Tillage 5 hours
//  PASS
//  ok      github.com/browny/inject        0.018s
package inject

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/facebookgo/inject"
)

type iConstruct interface {
	Setup() error
}

// Weave sets up dependencies and returns the result graph.
//
// `depMap` is the map describing the dependency relations.
// The key of depMap is the reference to the dependency providing object.
// The value is the list of dependency requiring objects.
func Weave(depMap map[interface{}][]string) (map[reflect.Type]interface{}, error) {
	deps := buildObjects(depMap)
	return run(deps...)
}

func run(deps ...*inject.Object) (map[reflect.Type]interface{}, error) {
	var g inject.Graph
	err := g.Provide(deps...)
	if err != nil {
		return nil, err
	}

	err = g.Populate()
	if err != nil {
		return nil, err
	}

	// constructing
	c := newConstruct()
	initializedDeps, err := c.populate(g.Objects())
	if err != nil {
		return nil, err
	}

	return initializedDeps, nil
}

func newConstruct() *construct {
	return &construct{
		constructed:     make(map[reflect.Type]bool),
		initializedDeps: make(map[reflect.Type]interface{}),
	}
}

type construct struct {
	current         interface{}
	constructed     map[reflect.Type]bool
	initializedDeps map[reflect.Type]interface{}
}

func (c *construct) populate(objs []*inject.Object) (map[reflect.Type]interface{}, error) {
	objs = dedup(objs)

	var errs []error
	for _, obj := range objs {
		c.current = obj.Value
		if err := c.recurConstruct(obj); err != nil {
			errs = append(errs, err)
		}
	}

	return c.initializedDeps, combine(errs...)
}

func (c *construct) recurConstruct(obj *inject.Object) error {
	// has been constructed
	if c.constructed[reflect.TypeOf(obj.Value)] {
		return nil
	}

	// without constructor, look it as constructed
	if !withConstructor(obj.Value) {
		c.constructed[reflect.TypeOf(obj.Value)] = true
		c.initializedDeps[reflect.TypeOf(obj.Value)] = obj.Value
		return nil
	}

	objVal := reflect.ValueOf(obj.Value)
	objTyp := reflect.TypeOf(obj.Value)

	// make sure all members with `inject` tag have been constructed
	for i := 0; i < objVal.Elem().NumField(); i++ {
		fieldVal := objVal.Elem().Field(i)
		fieldTyp := objTyp.Elem().Field(i)

		_, ok := fieldTyp.Tag.Lookup("inject")
		if !ok {
			continue
		}

		if fieldVal.Interface() == c.current {
			return fmt.Errorf(
				"Dep loop: curr[%s], obj[%s]",
				reflect.TypeOf(c.current).String(), reflect.TypeOf(obj.Value).String(),
			)
		}

		err := c.recurConstruct(&inject.Object{Value: fieldVal.Interface()})
		if err != nil {
			return err
		}
	}

	ctor := obj.Value.(iConstruct)
	err := ctor.Setup()
	if err != nil {
		return err
	}

	c.constructed[reflect.TypeOf(obj.Value)] = true
	c.initializedDeps[reflect.TypeOf(obj.Value)] = obj.Value

	return nil
}

func buildObjects(depMap map[interface{}][]string) []*inject.Object {
	deps := []*inject.Object{}
	for value, names := range depMap {
		if len(names) == 0 {
			deps = append(deps, &inject.Object{
				Value: value,
			})
		}
		// named dependencies
		for _, name := range names {
			deps = append(deps, &inject.Object{
				Value: value,
				Name:  name,
			})
		}
	}
	return deps
}

func withConstructor(i interface{}) bool {
	_, ok := i.(iConstruct)
	if ok {
		return true
	}
	return false
}

// dedup removes duplication
func dedup(objs []*inject.Object) []*inject.Object {
	result := []*inject.Object{}
	bucket := map[reflect.Type]bool{}

	for _, obj := range objs {
		objType := reflect.TypeOf(obj.Value)

		if bucket[objType] {
			continue
		}

		bucket[objType] = true
		result = append(result, obj)
	}

	return result
}

type multiError struct {
	errs []error
}

func (e *multiError) Error() string {
	b := bytes.NewBuffer(nil)
	for i, err := range e.errs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

func combine(errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	nonNilErrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			nonNilErrs = append(nonNilErrs, err)
		}
	}

	if len(nonNilErrs) == 0 {
		return nil
	}
	if len(nonNilErrs) == 1 {
		return nonNilErrs[0]
	}
	return &multiError{
		errs: nonNilErrs,
	}
}
