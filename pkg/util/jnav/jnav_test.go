package jnav

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestSimpleCasting(t *testing.T) {
    t.Run("Jnav generic", func(t *testing.T) {
        m := []interface{}{
            1,
            2,
            3,
        }
        j := FromArr(m)
        j.Next()
        resInt, ok := j.Current().AsInt()
        assert.True(t, ok)
        assert.Equal(t, 1, resInt)
        resStr, ok := j.Current().AsString()
        assert.False(t, ok)
        assert.Empty(t, resStr)
    })
}

func TestNav(t *testing.T) {
    t.Run("Jnav Nav", func(t *testing.T) {
        m := map[string]interface{}{
            "root": []interface{}{
                1,
                map[string]interface{}{
                    "a": 1,
                    "b": map[string]interface{}{
                        "X": "x",
                    },
                    "c": 3,
                },
                3,
            },
        }

        j1 := FromMap(m).Nav("root").AsArr()
        assert.Equal(t, 3, j1.Len())
        assert.Equal(t, 1, j1.Get(0).AsIntOr(-1))
        j2 := FromMap(m).Nav("root.[1].a").AsIntOr(-1)
        assert.Equal(t, 1, j2)
    })
}

func TestRewind(t *testing.T) {
    t.Run("Jnav generic", func(t *testing.T) {
        m := []interface{}{
            1,
            2,
            3,
        }
        j := FromArr(m)
        assert.Equal(t, 3, j.Len())

        assert.True(t, j.Next())
        assert.Equal(t, 1, j.Current().AsIntOr(-1))
        assert.True(t, j.Next())
        assert.Equal(t, 2, j.Current().AsIntOr(-1))
        assert.True(t, j.Next())
        assert.Equal(t, 3, j.Current().AsIntOr(-1))
        assert.False(t, j.Next())
        j.Rewind()
        assert.Equal(t, -1, j.Current().AsIntOr(-1))
        assert.True(t, j.Next())
        assert.Equal(t, 1, j.Current().AsIntOr(-1))
    })
}

func TestJnav(t *testing.T) {
    t.Run("Jnav generic", func(t *testing.T) {
        m := map[string]interface{}{
            "a": 1,
            "b": 2,
            "c": 3,
            "sub": []map[string]interface{}{
                {
                    "xval": "x1",
                },
                {
                    "xval": "x2",
                },
            },
        }

        j := FromMap(m)

        assert.Equal(t, 1, j.Get("a").AsIntOr(-1))
        assert.Equal(t, 2, j.Get("b").AsIntOr(-1))
        assert.Equal(t, 3, j.Get("c").AsIntOr(-1))

        // ARRAY "sub"
        // ..SOF
        ja := j.Nav("sub").AsArr()
        multiAssert(t, ja, 0, true, false, "")

        // ..ELEMENT #1
        assert.True(t, ja.Next())
        multiAssert(t, ja, 1, false, false, "x1")

        // ..ELEMENT #2
        assert.True(t, ja.Next())
        multiAssert(t, ja, 2, false, false, "x2")

        // ..EOF
        assert.False(t, ja.Next())
        multiAssert(t, ja, 3, false, true, "")

        // ..OVER THE EOF?
        assert.False(t, ja.Next())
        multiAssert(t, ja, 3, false, true, "")
    })
}

func multiAssert(t *testing.T, ja Arr, n int, bof bool, eof bool, xval string) {
    assert.Equal(t, bof, ja.Bof())
    assert.Equal(t, bof, ja.First(0))
    assert.Equal(t, eof, ja.Eof())
    assert.Equal(t, eof, ja.Last(0))
    assert.Equal(t, n <= 0, ja.First(0))
    assert.Equal(t, n <= 1, ja.First(1))
    assert.Equal(t, n <= 2, ja.First(2))
    assert.Equal(t, n < 3, ja.First(3)) // If iterator out-of-bounds "First" always returns false
    assert.Equal(t, n >= 3-2, ja.Last(2))
    assert.Equal(t, n >= 3-1, ja.Last(1))
    assert.Equal(t, xval, ja.Current().AsObj().Get("xval").AsStringOr(""))
}

func TestStrangeCases(t *testing.T) {
    type Custom struct {
        val int
    }

    m := map[string]interface{}{
        "custom": []Custom{
            {1},
            {2},
            {3},
        },
    }

    t.Run("Jnav generic", func(t *testing.T) {
        defer func() {
            if r := recover(); r == nil {
                t.Errorf("The test did not panic as expected")
            }
        }()

        j := FromMap(m).Get("custom").AsArr()
        j.next()
        resInt, ok := j.Current().AsInt()
        assert.True(t, ok)
        assert.Equal(t, 1, resInt)
    })

    t.Run("Jnav generic", func(t *testing.T) {
        defer func() {
            if r := recover(); r == nil {
                t.Errorf("The test did not panic as expected")
            }
        }()

        j := FromMap(m).Get("custom").AsArr()

        j.next()
        resStr, ok := j.Current().AsString()
        assert.False(t, ok)
        assert.Empty(t, resStr)
    })
}
