package echoswagger

import (
	"encoding/xml"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaSwaggerTags(t *testing.T) {
	type Spot struct {
		Address string   `swagger:"desc(Address of Spot)"`
		Matrix  [][]bool `swagger:"default(true)"`
	}

	type User struct {
		Age    int       `swagger:"min(0),max(99)"`
		Gender string    `swagger:"enum(male|female|other),required"`
		CarNos string    `swagger:"minLen(5),maxLen(8)"`
		Spots  []*Spot   `swagger:"required"`
		Money  **float64 `swagger:"default(0),readOnly"`
	}

	a := prepareApi()
	a.AddParamBody(&User{}, "Body", "", true)
	sapi := a.(*api)
	assert.Len(t, sapi.operation.Parameters, 1)
	assert.Len(t, *sapi.defs, 2)

	su := (*sapi.defs)["User"].Schema
	pu := su.Properties
	assert.NotNil(t, su)
	assert.NotNil(t, pu)
	assert.Len(t, su.Required, 2)
	assert.ElementsMatch(t, su.Required, []string{"Spots", "Gender"})
	assert.Equal(t, *pu["Age"].Minimum, float64(0))
	assert.Equal(t, *pu["Age"].Maximum, float64(99))
	assert.Len(t, pu["Gender"].Enum, 3)
	assert.ElementsMatch(t, pu["Gender"].Enum, []string{"male", "female", "other"})
	assert.Equal(t, *pu["CarNos"].MinLength, int(5))
	assert.Equal(t, *pu["CarNos"].MaxLength, int(8))
	assert.Equal(t, pu["Money"].DefaultValue, float64(0))
	assert.Equal(t, pu["Money"].ReadOnly, true)

	ss := (*sapi.defs)["Spot"].Schema
	ps := ss.Properties
	assert.NotNil(t, ss)
	assert.NotNil(t, ps)
	assert.Equal(t, ps["Address"].Description, "Address of Spot")
	assert.Equal(t, ps["Matrix"].Items.Items.DefaultValue, true)
}

func TestParamSwaggerTags(t *testing.T) {
	type SearchInput struct {
		Q              string     `query:"q" swagger:"minLen(5),maxLen(8)"`
		BrandIds       string     `query:"brandIds" swagger:"allowEmpty"`
		Sortby         [][]string `query:"sortby" swagger:"default(id),allowEmpty"`
		Order          []int      `query:"order" swagger:"enum(0|1|n)"`
		SkipCount      int        `query:"skipCount" swagger:"min(0),max(999)"`
		MaxResultCount int        `query:"maxResultCount" swagger:"desc(items count in one page)"`
	}

	a := prepareApi()
	a.AddParamQueryNested(SearchInput{})
	o := a.(*api).operation
	assert.Len(t, o.Parameters, 6)
	assert.Equal(t, *o.Parameters[0].MinLength, 5)
	assert.Equal(t, *o.Parameters[0].MaxLength, 8)
	assert.Equal(t, o.Parameters[1].AllowEmptyValue, true)
	assert.Equal(t, o.Parameters[2].AllowEmptyValue, true)
	assert.Equal(t, o.Parameters[2].Items.Items.Default, "id")
	assert.Equal(t, o.Parameters[2].Items.CollectionFormat, "multi")
	assert.ElementsMatch(t, o.Parameters[3].Items.Enum, []int{0, 1})
	assert.Equal(t, o.Parameters[3].CollectionFormat, "multi")
	assert.Equal(t, *o.Parameters[4].Minimum, float64(0))
	assert.Equal(t, *o.Parameters[4].Maximum, float64(999))
	assert.Equal(t, o.Parameters[5].Description, "items count in one page")
}

func TestHeaderSwaggerTags(t *testing.T) {
	type SearchInput struct {
		Q              string     `json:"q" swagger:"minLen(5),maxLen(8)"`
		Enable         bool       `json:"-"`
		Sortby         [][]string `json:"sortby" swagger:"default(id)"`
		Order          []int      `json:"order" swagger:"enum(0|1|n)"`
		SkipCount      int        `json:"skipCount" swagger:"min(0),max(999)"`
		MaxResultCount int        `json:"maxResultCount" swagger:"desc(items count in one page)"`
	}

	a := prepareApi()
	a.AddResponse(http.StatusOK, "Resp", nil, SearchInput{})
	o := a.(*api).operation
	c := strconv.Itoa(http.StatusOK)
	h := o.Responses[c].Headers
	assert.Len(t, h, 5)
	assert.Equal(t, *h["q"].MinLength, 5)
	assert.Equal(t, *h["q"].MaxLength, 8)
	assert.Equal(t, h["sortby"].Items.Items.Default, "id")
	assert.Equal(t, h["sortby"].Items.CollectionFormat, "multi")
	assert.ElementsMatch(t, h["order"].Items.Enum, []int{0, 1})
	assert.Equal(t, h["order"].CollectionFormat, "multi")
	assert.Equal(t, *h["skipCount"].Minimum, float64(0))
	assert.Equal(t, *h["skipCount"].Maximum, float64(999))
	assert.Equal(t, h["maxResultCount"].Description, "items count in one page")
}

func TestXMLTags(t *testing.T) {
	type Spot struct {
		Id      int64  `xml:",attr"`
		Comment string `xml:",comment"`
		Address string `xml:"AddressDetail"`
		Enable  bool   `xml:"-"`
	}

	type User struct {
		X     xml.Name `xml:"Users"`
		Spots []*Spot  `xml:"Spots>Spot"`
	}

	a := prepareApi()
	a.AddParamBody(&User{}, "Body", "", true)
	sapi, ok := a.(*api)
	assert.Equal(t, ok, true)
	assert.Len(t, sapi.operation.Parameters, 1)
	assert.Len(t, *sapi.defs, 2)

	su := (*sapi.defs)["User"].Schema
	pu := su.Properties
	assert.NotNil(t, su)
	assert.NotNil(t, pu)
	assert.Equal(t, su.XML.Name, "Users")
	assert.NotNil(t, pu["Spots"].XML)
	assert.Equal(t, pu["Spots"].XML.Name, "Spots")
	assert.Equal(t, pu["Spots"].XML.Wrapped, true)

	ss := (*sapi.defs)["Spot"].Schema
	ps := ss.Properties
	assert.NotNil(t, ss)
	assert.NotNil(t, ps)
	assert.Equal(t, ss.XML.Name, "Spot")
	assert.Equal(t, ps["Id"].XML.Attribute, "Id")
	assert.Nil(t, ps["Comment"].XML)
	assert.Equal(t, ps["Address"].XML.Name, "AddressDetail")
	assert.Nil(t, ps["Enable"].XML)
}

func TestEnumInSchema(t *testing.T) {
	type User struct {
		Id      int64   `swagger:"enum(0|-1|200000|9.9)"`
		Age     int     `swagger:"enum(0|-1|200000|9.9)"`
		Status  string  `swagger:"enum(normal|stop)"`
		Amount  float64 `swagger:"enum(0|-0.1|ok|200.555)"`
		Grade   float32 `swagger:"enum(0|-0.5|ok|200.5)"`
		Deleted bool    `swagger:"enum(t|F),default(True)"`
	}

	a := prepareApi()
	a.AddParamBody(&User{}, "Body", "", true)
	sapi, ok := a.(*api)
	assert.Equal(t, ok, true)
	assert.Len(t, sapi.operation.Parameters, 1)
	assert.Len(t, *sapi.defs, 1)

	s := (*sapi.defs)["User"].Schema
	assert.NotNil(t, s)

	p := s.Properties

	assert.Len(t, p["Id"].Enum, 3)
	assert.ElementsMatch(t, p["Id"].Enum, []interface{}{int64(0), int64(-1), int64(200000)})

	assert.Len(t, p["Age"].Enum, 3)
	assert.ElementsMatch(t, p["Age"].Enum, []interface{}{0, -1, 200000})

	assert.Len(t, p["Status"].Enum, 2)
	assert.ElementsMatch(t, p["Status"].Enum, []interface{}{"normal", "stop"})

	assert.Len(t, p["Amount"].Enum, 3)
	assert.ElementsMatch(t, p["Amount"].Enum, []interface{}{float64(0), float64(-0.1), float64(200.555)})

	assert.Len(t, p["Grade"].Enum, 3)
	assert.ElementsMatch(t, p["Grade"].Enum, []interface{}{float32(0), float32(-0.5), float32(200.5)})

	assert.Len(t, p["Deleted"].Enum, 2)
	assert.ElementsMatch(t, p["Deleted"].Enum, []interface{}{true, false})
	assert.Equal(t, p["Deleted"].DefaultValue, true)
}

func TestExampleInSchema(t *testing.T) {
	u := struct {
		Id      int64
		Age     int
		Status  string
		Amount  float64
		Grade   float32
		Deleted bool
	}{
		Id:      10000000001,
		Age:     18,
		Status:  "normal",
		Amount:  195.50,
		Grade:   5.5,
		Deleted: true,
	}

	a := prepareApi()
	a.AddParamBody(u, "Body", "", true)
	sapi, ok := a.(*api)
	assert.Equal(t, ok, true)
	assert.Len(t, sapi.operation.Parameters, 1)
	assert.Len(t, *sapi.defs, 1)

	s := (*sapi.defs)[""].Schema
	assert.NotNil(t, s)

	p := s.Properties

	assert.Equal(t, p["Id"].Example, u.Id)
	assert.Equal(t, p["Age"].Example, u.Age)
	assert.Equal(t, p["Status"].Example, u.Status)
	assert.Equal(t, p["Amount"].Example, u.Amount)
	assert.Equal(t, p["Grade"].Example, u.Grade)
	assert.Equal(t, p["Deleted"].Example, u.Deleted)
}
