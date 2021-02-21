package spawn

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func Test_MkRawSpawnArgs(t *testing.T) {
    args := []interface{}{
        "a", WordArg{"b"}, SOA("o", "x", "y"),
    }

    assert.Equal(t, MkRawSpawnArgs(args), []string{"a", "b", "-o", "x=y"})

    type TheSpanishInquisition struct {
        isExpectedBy string
    }

    args = append(args, TheSpanishInquisition{isExpectedBy: "nobody"})
    assert.Panics(t, func() {
        MkRawSpawnArgs(args)
    })
}

func TestArgsAppendPrepend(t *testing.T) {
    args := []interface{}{"middle"}
    args = AppendToArgs(args, "last")
    args = PrependToArgs(args, "first")
    assert.Equal(t, 3, len(args))
    assert.Equal(t, []interface{}{"first", "middle", "last"}, args)
}
