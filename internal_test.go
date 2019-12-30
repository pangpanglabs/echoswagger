package echoswagger

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParamTypes(t *testing.T) {
	var pa interface{}
	var pb *int64
	var pc map[string]string
	var pd [][]float64
	type Parent struct {
		Child struct {
			Name string
		}
	}
	var pe Parent
	tests := []struct {
		p     interface{}
		panic bool
		name  string
	}{
		{
			p:     pa,
			panic: true,
			name:  "Interface type",
		},
		{
			p:     &pa,
			panic: true,
			name:  "Interface pointer type",
		},
		{
			p:     &pb,
			panic: false,
			name:  "Int type",
		},
		{
			p:     &pc,
			panic: true,
			name:  "Map type",
		},
		{
			p:     nil,
			panic: true,
			name:  "Nil type",
		},
		{
			p:     0,
			panic: false,
			name:  "Int type",
		},
		{
			p:     &pd,
			panic: false,
			name:  "Array float64 type",
		},
		{
			p:     &pe,
			panic: true,
			name:  "Struct type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := prepareApi()
			if tt.panic {
				assert.Panics(t, func() {
					a.AddParamPath(tt.p, tt.name, "")
				})
			} else {
				a.AddParamPath(tt.p, tt.name, "")
				sapi, ok := a.(*api)
				assert.Equal(t, ok, true)
				assert.Equal(t, len(sapi.operation.Parameters), 1)
				assert.Equal(t, tt.name, sapi.operation.Parameters[0].Name)
			}
		})
	}
}

func TestNestedParamTypes(t *testing.T) {
	var pa struct {
		ExpiredAt time.Time
	}
	type User struct {
		ExpiredAt time.Time
	}
	var pb struct {
		User User
	}
	type Org struct {
		Id      int64  `json:"id"`
		Address string `json:"address"`
	}
	var pc struct {
		User
		Org
	}
	var pd struct {
		User
		Org `json:"org"` // this tag would be ignored
	}

	tests := []struct {
		p      interface{}
		panic  bool
		name   string
		params []string
	}{
		{
			p:     0,
			panic: true,
			name:  "Basic type",
		},
		{
			p:      pa,
			panic:  false,
			name:   "Struct type",
			params: []string{"ExpiredAt"},
		},
		{
			p:     pb,
			panic: true,
			name:  "Nested struct type",
		},
		{
			p:      pc,
			panic:  false,
			name:   "Embedded struct type",
			params: []string{"ExpiredAt", "id", "address"},
		},
		{
			p:      pd,
			panic:  false,
			name:   "Embedded struct type with tag",
			params: []string{"ExpiredAt", "id", "address"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := prepareApi()
			if tt.panic {
				assert.Panics(t, func() {
					a.AddParamPathNested(tt.p)
				})
			} else {
				a.AddParamPathNested(tt.p)
				sapi, ok := a.(*api)
				assert.Equal(t, ok, true)
				assert.Equal(t, len(sapi.operation.Parameters), len(tt.params))
				var params []string
				for _, p := range sapi.operation.Parameters {
					params = append(params, p.Name)
				}
				assert.ElementsMatch(t, params, tt.params)
			}
		})
	}
}

func TestSchemaTypes(t *testing.T) {
	var pa interface{}
	var pb map[string]string
	type PT struct {
		Name      string
		ExpiredAt time.Time
	}
	var pc map[PT]string
	var pd PT
	var pe map[time.Time]string
	var pf map[*int]string
	type PU struct {
		Any interface{}
	}
	var pg PU
	var ph map[string]interface{}
	tests := []struct {
		p     interface{}
		panic bool
		name  string
	}{
		{
			p:     pa,
			panic: true,
			name:  "Interface type",
		},
		{
			p:     nil,
			panic: true,
			name:  "Nil type",
		},
		{
			p:     "",
			panic: false,
			name:  "String type",
		},
		{
			p:     &pb,
			panic: false,
			name:  "Map type",
		},
		{
			p:     &pc,
			panic: true,
			name:  "Map struct type",
		},
		{
			p:     pd,
			panic: false,
			name:  "Struct type",
		},
		{
			p:     &pd,
			panic: false,
			name:  "Struct pointer type",
		},
		{
			p:     &pe,
			panic: false,
			name:  "Map time.Time key type",
		},
		{
			p:     &pf,
			panic: false,
			name:  "Map pointer key type",
		},
		{
			p:     &pg,
			panic: false,
			name:  "Struct interface field type",
		},
		{
			p:     &ph,
			panic: false,
			name:  "Map interface value type",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := prepareApi()
			if tt.panic {
				assert.Panics(t, func() {
					a.AddParamBody(tt.p, tt.name, "", true)
				})
			} else {
				a.AddParamBody(tt.p, tt.name, "", true)
				sapi, ok := a.(*api)
				assert.Equal(t, ok, true)
				assert.Equal(t, len(sapi.operation.Parameters), 1)
				assert.Equal(t, tt.name, sapi.operation.Parameters[0].Name)
			}
		})
	}
}

type testUser struct {
	Id   int64
	Name string
	Pets []testPet
}

type testPet struct {
	Id      int64
	Masters []testUser
}

func TestSchemaRecursiveStruct(t *testing.T) {
	tests := []struct {
		p    interface{}
		name string
	}{
		{
			p:    &testUser{},
			name: "User",
		},
		{
			p:    &testPet{},
			name: "Pet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := prepareApi()
			a.AddParamBody(tt.p, tt.name, "", true)
			sapi, ok := a.(*api)
			assert.Equal(t, ok, true)
			assert.Equal(t, len(sapi.operation.Parameters), 1)
			assert.Equal(t, len(*sapi.defs), 2)
			assert.Equal(t, tt.name, sapi.operation.Parameters[0].Name)
		})
	}
}

func TestSchemaNestedStruct(t *testing.T) {
	type User struct {
		ExpiredAt time.Time
	}
	type Org struct {
		Id      int64  `json:"id"`
		Address string `json:"address"`
	}
	var pa struct {
		User `json:"user"`
		Org
	}
	a := prepareApi()
	a.AddParamBody(pa, "pa", "", true)
	sapi, ok := a.(*api)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(sapi.operation.Parameters), 1)
	assert.NotNil(t, (*sapi.defs)[""])
	assert.NotNil(t, (*sapi.defs)[""].Schema.Properties["address"])
	assert.NotNil(t, (*sapi.defs)[""].Schema.Properties["id"])
	assert.NotNil(t, (*sapi.defs)["User"])
	assert.NotNil(t, (*sapi.defs)["User"].Schema.Properties["ExpiredAt"])
}
