package message_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kode4food/timebox/message"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string
	Age  int
}

func TestTypedInstantiators(t *testing.T) {
	as := assert.New(t)
	ti1 := message.TypedInstantiator{
		"test1": func() message.Payload {
			return new(testStruct)
		},
	}
	ti2 := message.TypedInstantiator{
		"test2": func() message.Payload {
			return new(testStruct)
		},
	}
	ti3 := ti1.Combine(ti2)
	i, ok := ti3["test1"]
	as.True(ok)
	as.Equal(reflect.ValueOf(ti1["test1"]), reflect.ValueOf(i))

	i, ok = ti3["test2"]
	as.True(ok)
	as.Equal(reflect.ValueOf(ti2["test2"]), reflect.ValueOf(i))
}

func TestTypedInstantiatorDecoder(t *testing.T) {
	as := assert.New(t)
	in := &message.Message{Type: "test",
		Payload: &testStruct{
			Name: "Bill",
			Age:  42,
		},
	}

	m, err := json.Marshal(in)
	as.NotNil(m)
	as.Nil(err)

	d := message.TypedInstantiator{
		"test": func() message.Payload {
			return new(testStruct)
		},
	}.Decoder()

	out, err := d(m)
	as.NotNil(out)
	as.Nil(err)

	p := out.Payload.(*testStruct)
	as.Equal(message.Type("test"), out.Type)
	as.Equal("Bill", p.Name)
	as.Equal(42, p.Age)
}
