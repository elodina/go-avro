/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package avro

import (
	"reflect"
	"strings"
	"sync"
)

func findField(where reflect.Value, name string) (reflect.Value, error) {
	if where.Kind() == reflect.Ptr {
		where = where.Elem()
	}
	rm := reflectEnsureRi(where.Type())
	if rf, ok := rm.names[name]; ok {
		return where.FieldByIndex(rf), nil
	}
	return reflect.Value{}, NewFieldDoesNotExistError(name)
}

func reflectEnsureRi(t reflect.Type) *reflectInfo {
	reflectMapLock.RLock()
	rm := reflectMap[t]
	reflectMapLock.RUnlock()
	if rm == nil {
		rm = reflectBuildRi(t)
	}
	return rm
}

func reflectBuildRi(t reflect.Type) *reflectInfo {
	rm := &reflectInfo{
		names: make(map[string][]int),
	}
	rm.fill(t, nil)

	reflectMapLock.Lock()
	reflectMap[t] = rm
	reflectMapLock.Unlock()
	return rm
}

var reflectMap = make(map[reflect.Type]*reflectInfo)
var reflectMapLock sync.RWMutex

type reflectInfo struct {
	names map[string][]int
}

// fill the given reflect info with the field names mapped.
//
// fill will recurse into anonymous structs incrementing the index prefix
// so that untagged anonymous structs can be used as the source of truth.
func (rm *reflectInfo) fill(t reflect.Type, indexPrefix []int) {
	// simple infinite recursion preventer: stop when we are >10 deep.
	if len(indexPrefix) > 10 {
		return
	}

	fillName := func(tag string, idx []int) {
		if _, ok := rm.names[tag]; !ok {
			rm.names[tag] = idx
		}
	}
	// these are anonymous structs to investigate (tail recursion)
	var toInvestigate [][]int
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("avro")
		idx := append(append([]int{}, indexPrefix...), f.Index...)

		if f.Anonymous && tag == "" && f.Type.Kind() == reflect.Struct {
			toInvestigate = append(toInvestigate, idx)
		} else if strings.ToLower(f.Name[:1]) != f.Name[:1] {
			if tag != "" {
				fillName(tag, idx)
			} else {
				fillName(f.Name, idx)
				fillName(strings.ToLower(f.Name[:1])+f.Name[1:], idx)
			}
		}
	}
	for _, idx := range toInvestigate {
		// recurse into anonymous structs now that we handled the base ones.
		rm.fill(t.Field(idx[len(idx)-1]).Type, idx)
	}
}
