package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	pb "github.com/kahlys/battools/blaze/proto/basic"
)

func TestProtoText(t *testing.T) {
	data := pb.Hero{
		Name:      "batman",
		Age:       30,
		Gender:    pb.Gender_MALE,
		Abilities: []string{"money", "tech"},
		Attributes: map[string]string{
			"wealth": "infinite",
			"power":  "none",
		},
	}

	tests := map[string]struct {
		expected string
		marshal  func(proto.Message) ([]byte, error)
		equal    func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		"text": {
			expected: `name:"batman"  age:30  abilities:"money"  abilities:"tech"  attributes:{key:"power"  value:"none"}  attributes:{key:"wealth"  value:"infinite"}`,
			marshal:  prototext.Marshal,
			equal: func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
				return assert.Equal(t, expected.(string), actual.(string), msgAndArgs...)
			},
		},
		"json": {
			expected: `{"name":"batman","age":30,"abilities":["money","tech"],"attributes":{"power":"none","wealth":"infinite"}}`,
			marshal:  protojson.Marshal,
			equal: func(t assert.TestingT, expected, actual interface{}, msgAndArgs ...interface{}) bool {
				return assert.JSONEq(t, expected.(string), actual.(string), msgAndArgs...)
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual, err := tt.marshal(&data)
			require.NoError(t, err)
			tt.equal(t, tt.expected, string(actual))
		})
	}
}
