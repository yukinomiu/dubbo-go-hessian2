/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package hessian

import (
	"reflect"
	"testing"
)

import (
	big "github.com/dubbogo/gost/math/big"
	"github.com/stretchr/testify/assert"
)

func TestEncUntypedMap(t *testing.T) {
	var (
		m   map[interface{}]interface{}
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	m = make(map[interface{}]interface{})
	m["hello"] = "world"
	m[100] = "100"
	m[100.1010] = 101910
	m[true] = true
	m[false] = true
	e.Encode(m)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", m, res, err)
}

func TestEncTypedMap(t *testing.T) {
	var (
		m   map[int]string
		err error
		e   *Encoder
		d   *Decoder
		res interface{}
	)

	e = NewEncoder()
	m = make(map[int]string)
	m[0] = "hello"
	m[1] = "golang"
	m[2] = "world"
	e.Encode(m)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", m, res, err)
}

func TestMap(t *testing.T) {
	testDecodeFramework(t, "replyTypedMap_0", map[interface{}]interface{}{})
	testDecodeFramework(t, "replyTypedMap_1", map[interface{}]interface{}{"a": int32(0)})
	testDecodeFramework(t, "replyTypedMap_2", map[interface{}]interface{}{int32(0): "a", int32(1): "b"})
	// testDecodeFramework(t, "replyTypedMap_3", []interface{}{})
	testDecodeFramework(t, "replyUntypedMap_0", map[interface{}]interface{}{})
	testDecodeFramework(t, "replyUntypedMap_1", map[interface{}]interface{}{"a": int32(0)})
	testDecodeFramework(t, "replyUntypedMap_2", map[interface{}]interface{}{int32(0): "a", int32(1): "b"})
	// testDecodeFramework(t, "replyTypedMap_3", []interface{}{})
}

func TestMapEncode(t *testing.T) {
	testJavaDecode(t, "argTypedMap_0", map[interface{}]interface{}{})
	testJavaDecode(t, "argTypedMap_1", map[interface{}]interface{}{"a": int32(0)})
	testJavaDecode(t, "argTypedMap_2", map[interface{}]interface{}{int32(0): "a", int32(1): "b"})
	testJavaDecode(t, "argUntypedMap_0", map[interface{}]interface{}{})
	testJavaDecode(t, "argUntypedMap_1", map[interface{}]interface{}{"a": int32(0)})
	testJavaDecode(t, "argUntypedMap_2", map[interface{}]interface{}{int32(0): "a", int32(1): "b"})
}

func TestCustomMapRefMap(t *testing.T) {
	r, e := decodeJavaResponse("customReplyMapRefMap", "", true)
	if e != nil {
		t.Errorf("%s: decode fail with error: %v", "customReplyMapRefMap", e)
		return
	}
	res := r.(map[interface{}]interface{})
	assert.Equal(t, int32(1), res["a"])
	assert.Equal(t, int32(2), res["b"])
	assert.Equal(t, res, res["self"])
}

type customMapObject struct {
	Int int32
	S   string
}

func TestCustomMap(t *testing.T) {
	testDecodeFramework(t, "customReplyMap", map[interface{}]interface{}{"a": int32(1), "b": int32(2)})

	mapInMap := map[interface{}]interface{}{
		"obj1": map[interface{}]interface{}{"a": int32(1)},
		"obj2": map[interface{}]interface{}{"b": int32(2)},
	}
	testDecodeFramework(t, "customReplyMapInMap", mapInMap)
	testDecodeFramework(t, "customReplyMapInMapJsonObject", mapInMap)

	b3 := &big.Decimal{}
	_ = b3.FromString("33.33")
	b3.Value = "33.33"

	b5 := &big.Decimal{}
	_ = b5.FromString("55.55")
	b5.Value = "55.55"

	multipleTypeMap := map[interface{}]interface{}{
		"m1": map[interface{}]interface{}{"a": int32(1), "b": int32(2)},
		"m2": map[interface{}]interface{}{int64(3): "c", int64(4): "d"},
		"m3": map[interface{}]interface{}{int32(3): b3, int32(5): b5},
	}

	testDecodeFramework(t, "customReplyMultipleTypeMap", multipleTypeMap)

	RegisterPOJOMapping("test.model.CustomMap", &customMapObject{})

	listMapListMap := []interface{}{

		map[interface{}]interface{}{
			"a": int32(1),
			"b": int32(2),
			"items": []interface{}{
				b5,
				"hello",
				int32(123),
				customMapObject{
					Int: 456,
					S:   "string",
				},
			},
		},
		customMapObject{
			Int: 789,
			S:   "string2",
		},
	}

	testDecodeFramework(t, "customReplyListMapListMap", listMapListMap)
}

type CustomMap map[string]interface{}

func (dict *CustomMap) JavaClassName() string {
	return "test.model.CustomMap"
}

func TestJavaMap(t *testing.T) {
	customMap := &CustomMap{"Name": "Test"}
	RegisterPOJO(customMap)
	testJavaDecode(t, "customArgTypedFixed_CustomMap", customMap)
}

type Obj struct {
	Map8    map[int8]int8
	Map8P   map[int8]*int8
	Map16   map[int16]int16
	Map16P  map[int16]*int16
	Map32   map[int32]int32
	Map32P  map[int32]*int32
	MapStr  map[string]string
	MapStrP map[string]*string
}

func (Obj) JavaClassName() string {
	return ""
}

func TestMapInObject(t *testing.T) {
	var (
		req *Obj
		e   *Encoder
		d   *Decoder
		err error
		res interface{}
	)

	int8Val := int8(0)
	int16Val := int16(0)
	int32Val := int32(0)
	strVal := "string"
	req = &Obj{
		Map8:    map[int8]int8{1: 2, 3: 4},
		Map8P:   map[int8]*int8{0: &int8Val},
		Map16:   map[int16]int16{1: 2, 3: 4},
		Map16P:  map[int16]*int16{0: &int16Val},
		Map32:   map[int32]int32{1: 2, 3: 4},
		Map32P:  map[int32]*int32{0: &int32Val},
		MapStr:  map[string]string{"a": "b", "c": "d"},
		MapStrP: map[string]*string{"non-nil": &strVal, "nil": nil},
	}

	e = NewEncoder()
	e.Encode(req)
	if len(e.Buffer()) == 0 {
		t.Fail()
	}

	d = NewDecoder(e.Buffer())
	res, err = d.Decode()
	if err != nil {
		t.Errorf("Decode() = %+v", err)
	}
	t.Logf("decode(%v) = %v, %v\n", req, res, err)

	if !reflect.DeepEqual(req, res) {
		t.Fatalf("req: %#v != res: %#v", req, res)
	}
}
