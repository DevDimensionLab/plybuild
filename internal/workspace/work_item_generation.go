package workspace

import (
	"os"
	"reflect"
)

// A generation belongs only to this lock period. External content is deliberately
// not retained: Snapshot validates it using the store's current scope each time.
// The stat observation also drops reuse on out-of-protocol physical changes.
func (s *systemWorkItemStoreSession) confirm(r WorkItemRegistry, observed os.FileInfo) {
	info, err := os.Lstat(workItemsPath(s.root))
	if err != nil || observed == nil || !os.SameFile(info, observed) || !sameRegistryFileState(info, observed) {
		s.confirmed, s.observed = nil, nil
		return
	}
	owned := cloneWorkItemRegistry(r)
	s.confirmed, s.observed = &owned, info
}

func cloneWorkItemRegistry(r WorkItemRegistry) WorkItemRegistry {
	// The private upgrade permit is an ephemeral publication capability, never a
	// part of a confirmed snapshot. All persisted reference-bearing fields are
	// recursively copied, including queue requests with nested maps/interfaces.
	r.upgradePermit = nil
	return cloneRegistryValue(reflect.ValueOf(r)).Interface().(WorkItemRegistry)
}

func cloneRegistryValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		out := reflect.New(v.Type()).Elem()
		if v.Kind() == reflect.Pointer {
			out.Set(reflect.New(v.Type().Elem()))
			out.Elem().Set(cloneRegistryValue(v.Elem()))
		} else {
			out.Set(cloneRegistryValue(v.Elem()))
		}
		return out
	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		out := reflect.MakeSlice(v.Type(), v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			out.Index(i).Set(cloneRegistryValue(v.Index(i)))
		}
		return out
	case reflect.Array:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			out.Index(i).Set(cloneRegistryValue(v.Index(i)))
		}
		return out
	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		out := reflect.MakeMapWithSize(v.Type(), v.Len())
		iter := v.MapRange()
		for iter.Next() {
			out.SetMapIndex(iter.Key(), cloneRegistryValue(iter.Value()))
		}
		return out
	case reflect.Struct:
		out := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if out.Field(i).CanSet() {
				out.Field(i).Set(cloneRegistryValue(v.Field(i)))
			}
		}
		return out
	default:
		return v
	}
}

// Access time may change when a guard or another reader observes the file.
// Compare the remaining native stat fields (including change time), so an
// in-place rewrite with a restored modification time still invalidates reuse.
func sameRegistryFileState(a, b os.FileInfo) bool {
	stat := func(info os.FileInfo) any {
		v := reflect.ValueOf(info.Sys())
		if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
			return info
		}
		out := reflect.New(v.Elem().Type()).Elem()
		out.Set(v.Elem())
		for _, name := range []string{"Atim", "Atimespec", "Atime", "Atimensec", "LastAccessTime"} {
			if f := out.FieldByName(name); f.IsValid() && f.CanSet() {
				f.SetZero()
			}
		}
		return out.Interface()
	}
	return reflect.DeepEqual(stat(a), stat(b))
}
