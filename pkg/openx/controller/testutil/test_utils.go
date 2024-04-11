package testutil

import (
	"k8s.io/client-go/tools/cache"
	"reflect"
	"testing"
)

var (
	keyFunc = cache.DeletionHandlingMetaNamespaceKeyFunc
)

// GetKey is a helper function used by controllers unit tests to get the
// key for a given kubernetes resource.
func GetKey(obj interface{}, t *testing.T) string {
	tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
	if ok {
		// if tombstone , try getting the value from tombstone.Obj
		obj = tombstone.Obj
	}
	val := reflect.ValueOf(obj).Elem()
	name := val.FieldByName("Name").String()
	kind := val.FieldByName("Kind").String()
	// Note kind is not always set in the tests, so ignoring that for now
	if len(name) == 0 || len(kind) == 0 {
		t.Errorf("Unexpected object %v", obj)
	}

	key, err := keyFunc(obj)
	if err != nil {
		t.Errorf("Unexpected error getting key for %v %v: %v", kind, name, err)
		return ""
	}
	return key
}
