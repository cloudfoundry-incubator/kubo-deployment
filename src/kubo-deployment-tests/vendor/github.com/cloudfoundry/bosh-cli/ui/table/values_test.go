package table_test

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	semver "github.com/cppforlife/go-semi-semantic/version"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-cli/ui/table"
)

var _ = Describe("ValueString", func() {
	It("returns string", func() {
		Expect(ValueString{S: "val"}.String()).To(Equal("val"))
	})

	It("returns itself", func() {
		Expect(ValueString{S: "val"}.Value()).To(Equal(ValueString{S: "val"}))
	})

	It("returns int based on string compare", func() {
		Expect(ValueString{S: "a"}.Compare(ValueString{S: "a"})).To(Equal(0))
		Expect(ValueString{S: "a"}.Compare(ValueString{S: "b"})).To(Equal(-1))
		Expect(ValueString{S: "b"}.Compare(ValueString{S: "a"})).To(Equal(1))
	})
})

var _ = Describe("ValueStrings", func() {
	It("returns new line joined strings", func() {
		Expect(ValueStrings{S: []string{"val1", "val2"}}.String()).To(Equal("val1\nval2"))
	})

	It("returns itself", func() {
		Expect(ValueStrings{S: []string{"val1"}}.Value()).To(Equal(ValueStrings{S: []string{"val1"}}))
	})

	It("returns int based on string compare", func() {
		Expect(ValueStrings{S: []string{"val1"}}.Compare(ValueStrings{S: []string{"val1"}})).To(Equal(0))
		Expect(ValueStrings{S: []string{"val1"}}.Compare(ValueStrings{S: []string{"val1", "val2"}})).To(Equal(-1))
		Expect(ValueStrings{S: []string{"val1", "val2"}}.Compare(ValueStrings{S: []string{"val1"}})).To(Equal(1))
	})
})

var _ = Describe("ValueInt", func() {
	It("returns string", func() {
		Expect(ValueInt{I: 1}.String()).To(Equal("1"))
	})

	It("returns itself", func() {
		Expect(ValueInt{I: 1}.Value()).To(Equal(ValueInt{I: 1}))
	})

	It("returns int based on int compare", func() {
		Expect(ValueInt{I: 1}.Compare(ValueInt{I: 1})).To(Equal(0))
		Expect(ValueInt{I: 1}.Compare(ValueInt{I: 2})).To(Equal(-1))
		Expect(ValueInt{I: 2}.Compare(ValueInt{I: 1})).To(Equal(1))
	})
})

var _ = Describe("ValueBytes", func() {
	It("returns formatted bytes", func() {
		Expect(ValueBytes{I: 1}.String()).To(Equal("1 B"))
	})

	It("returns formatted mebibytes", func() {
		Expect(NewValueMegaBytes(1).String()).To(Equal("1.0 MiB"))
	})

	It("returns formatted gibibytes", func() {
		Expect(NewValueMegaBytes(131072).String()).To(Equal("128 GiB"))
	})

	It("returns itself", func() {
		Expect(ValueBytes{I: 1}.Value()).To(Equal(ValueBytes{I: 1}))
	})

	It("returns int based on int compare", func() {
		Expect(ValueBytes{I: 1}.Compare(ValueBytes{I: 1})).To(Equal(0))
		Expect(ValueBytes{I: 1}.Compare(ValueBytes{I: 2})).To(Equal(-1))
		Expect(ValueBytes{I: 2}.Compare(ValueBytes{I: 1})).To(Equal(1))
	})
})

var _ = Describe("ValueTime", func() {
	t1 := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	t2 := time.Date(2009, time.November, 10, 23, 0, 0, 1, time.UTC)
	empty := time.Time{}

	It("returns formatted full time", func() {
		Expect(ValueTime{T: t1}.String()).To(Equal("Tue Nov 10 23:00:00 UTC 2009"))
	})

	It("returns empty", func() {
		Expect(ValueTime{T: empty}.String()).To(Equal(""))
	})

	It("returns itself", func() {
		Expect(ValueTime{T: t1}.Value()).To(Equal(ValueTime{T: t1}))
	})

	It("returns int based on time compare", func() {
		Expect(ValueTime{T: t1}.Compare(ValueTime{T: t1})).To(Equal(0))
		Expect(ValueTime{T: t1}.Compare(ValueTime{T: t2})).To(Equal(-1))
		Expect(ValueTime{T: t2}.Compare(ValueTime{T: t1})).To(Equal(1))
	})
})

var _ = Describe("ValueBool", func() {
	It("returns true/false as string", func() {
		Expect(ValueBool{B: true}.String()).To(Equal("true"))
		Expect(ValueBool{B: false}.String()).To(Equal("false"))
	})

	It("returns itself", func() {
		Expect(ValueBool{B: true}.Value()).To(Equal(ValueBool{B: true}))
	})

	It("returns int based on bool compare", func() {
		Expect(ValueBool{B: true}.Compare(ValueBool{B: true})).To(Equal(0))
		Expect(ValueBool{B: false}.Compare(ValueBool{B: true})).To(Equal(-1))
		Expect(ValueBool{B: true}.Compare(ValueBool{B: false})).To(Equal(1))
	})
})

var _ = Describe("ValueVersion", func() {
	v1 := semver.MustNewVersionFromString("1.1")
	v2 := semver.MustNewVersionFromString("1.2")

	It("returns formatted version", func() {
		Expect(ValueVersion{V: v1}.String()).To(Equal("1.1"))
	})

	It("returns itself", func() {
		Expect(ValueVersion{V: v1}.Value()).To(Equal(ValueVersion{V: v1}))
	})

	It("returns int based on version compare", func() {
		Expect(ValueVersion{V: v1}.Compare(ValueVersion{V: v1})).To(Equal(0))
		Expect(ValueVersion{V: v2}.Compare(ValueVersion{V: v1})).To(Equal(1))
		Expect(ValueVersion{V: v1}.Compare(ValueVersion{V: v2})).To(Equal(-1))
	})
})

var _ = Describe("ValueError", func() {
	It("returns empty string or error description", func() {
		Expect(ValueError{}.String()).To(Equal(""))
		Expect(ValueError{E: errors.New("err")}.String()).To(Equal("err"))
	})

	It("returns itself", func() {
		Expect(ValueError{E: errors.New("err")}.Value()).To(Equal(ValueError{E: errors.New("err")}))
	})

	It("does not allow comparison", func() {
		f := func() { ValueError{}.Compare(ValueError{}) }
		Expect(f).To(Panic())
	})
})

var _ = Describe("ValueNone", func() {
	It("returns empty string", func() {
		Expect(ValueNone{}.String()).To(Equal(""))
	})

	It("returns itself", func() {
		Expect(ValueNone{}.Value()).To(Equal(ValueNone{}))
	})

	It("does not allow comparison", func() {
		f := func() { ValueNone{}.Compare(ValueNone{}) }
		Expect(f).To(Panic())
	})
})

var _ = Describe("ValueFmt", func() {
	fmtFunc := func(pattern string, vals ...interface{}) string {
		return fmt.Sprintf(">%s<", fmt.Sprintf(pattern, vals...))
	}

	It("returns plain string (not formatted with fmt func)", func() {
		Expect(ValueFmt{V: ValueInt{I: 1}, Func: fmtFunc}.String()).To(Equal("1"))
	})

	It("returns wrapped value", func() {
		Expect(ValueFmt{V: ValueInt{I: 1}, Func: fmtFunc}.Value()).To(Equal(ValueInt{I: 1}))
	})

	It("does not allow comparison", func() {
		f := func() { ValueFmt{V: ValueInt{I: 1}, Func: fmtFunc}.Compare(ValueFmt{}) }
		Expect(f).To(Panic())
	})

	It("writes out value using custom Fprintf", func() {
		buf := bytes.NewBufferString("")
		ValueFmt{V: ValueInt{I: 1}, Func: fmtFunc}.Fprintf(buf, "%s,%s", "val1", "val2")
		Expect(buf.String()).To(Equal(">val1,val2<"))
	})

	It("uses fmt.Fprintf if fmt func is not set", func() {
		buf := bytes.NewBufferString("")
		ValueFmt{V: ValueInt{I: 1}}.Fprintf(buf, "%s,%s", "val1", "val2")
		Expect(buf.String()).To(Equal("val1,val2"))
	})
})

type failsToYAMLMarshal struct{}

func (s failsToYAMLMarshal) MarshalYAML() (interface{}, error) {
	return nil, errors.New("marshal-err")
}

var _ = Describe("ValueInterface", func() {
	It("returns map as a string", func() {
		i := map[string]interface{}{"key": "value", "num": 123}
		Expect(ValueInterface{I: i}.String()).To(Equal("key: value\nnum: 123"))
	})

	It("returns nested items as a string", func() {
		i := map[string]interface{}{"key": map[string]interface{}{"nested_key": "nested_value"}}
		Expect(ValueInterface{I: i}.String()).To(Equal("key:\n  nested_key: nested_value"))
	})

	It("returns nested items as a string", func() {
		i := failsToYAMLMarshal{}
		Expect(ValueInterface{I: i}.String()).To(Equal(`<serialization error> : table_test.failsToYAMLMarshal{}`))
	})

	It("returns nil items as blank string", func() {
		Expect(ValueInterface{I: nil}.String()).To(Equal(""))
	})

	It("returns an empty map as blank string", func() {
		i := map[string]interface{}{}
		Expect(ValueInterface{I: i}.String()).To(Equal(""))
	})

	It("returns an empty slice as blank string", func() {
		i := []string{}
		Expect(ValueInterface{I: i}.String()).To(Equal(""))
	})
})

var _ = Describe("ValueSuffix", func() {
	It("returns formatted string with suffix", func() {
		Expect(ValueSuffix{V: ValueInt{I: 1}, Suffix: "*"}.String()).To(Equal("1*"))
		Expect(ValueSuffix{V: ValueString{S: "val"}, Suffix: "*"}.String()).To(Equal("val*"))
	})

	It("returns wrapped value", func() {
		Expect(ValueSuffix{V: ValueInt{I: 1}, Suffix: "*"}.Value()).To(Equal(ValueInt{I: 1}))
	})

	It("does not allow comparison", func() {
		f := func() { ValueSuffix{V: ValueInt{I: 1}, Suffix: ""}.Compare(ValueSuffix{}) }
		Expect(f).To(Panic())
	})
})
