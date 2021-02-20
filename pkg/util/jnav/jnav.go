package jnav

import (
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

type IAny interface {
    Navigate(jnav Any, path string) Any
}

type Any struct {
    rawValue interface{}
}

type Obj struct {
    Any
}

type Arr struct {
    Any
    pos int
}

// Navigates to a JNAV element through a dot-delimited navigation syntax.
//
// Examples:
//   - root.key.otherKey    # object navigation
//   - root.[{index}].key   # object and array ({index} is a number) navigation
//
func (j Any) Nav(path string) Any {
    return Nav(j, path)
}

// Navigates to a JNAV element through a dot-delimited navigation syntax.
//
// Examples:
//   - root.key.otherKey    # object navigation
//   - root.[{index}].key   # object and array ({index} is a number) navigation
//
func Nav(jnav Any, path string) Any {
    for _, subPath := range strings.Split(path, ".") {
        match := ArrayIndexRegexp.FindStringSubmatch(subPath)
        if match != nil && len(match) == 2 {
            arrayIdx, err := strconv.Atoi(match[1])
            if err == nil {
                jnav = jnav.AsArr().Get(arrayIdx)
                continue
            }
        }
        jnav = jnav.AsObj().Get(subPath)
    }
    return jnav
}

var ArrayIndexRegexp = regexp.MustCompile(`^\[([^]]+)\]$`)

// Moves forward the iterator
//
// Returns false when it reaches the EOF
//
func (jn *Arr) Next() bool {
    if jn.Eof() {
        return false
    }
    jn.next()
    return !jn.Eof()
}

func (jn *Arr) next() {
    jn.pos++
}

// Return the number of elements of a JNAV array
func (jn Arr) Len() int {
    a := jn.rawValue
    switch a.(type) {
    case []interface{}:
        return len(a.([]interface{}))
    case []map[string]interface{}:
        return len(a.([]map[string]interface{}))
    default:
        panic(fmt.Sprintf("Unable to cast the  provided jnav (%T) to an array", a))
    }
}

// Tells if the iterator is on one of the last "n" elements
func (jn Arr) Last(n int) bool {
    return jn.last(jn.pos, n)
}

func (jn Arr) last(index int, n int) bool {
    l := jn.Len()
    return index >= 0 && index >= l-n && index <= l
}

// Tells if the iterator is on one of the first "n" elements
func (jn Arr) First(n int) bool {
    return jn.first(jn.pos, n)
}

func (jn Arr) first(index int, n int) bool {
    return n >= -1 && index < n && index < jn.Len()
}

// Tells if the iterator is on the EOF
func (jn Arr) Eof() bool {
    return jn.Last(0)
}

// Tells if the iterator is on  the EOF
func (jn Arr) Bof() bool {
    return jn.First(0)
}

// Returns the elements at the current iterator
func (jn Arr) Current() Any {
    return jn.Get(jn.pos)
}

// Returns the elements at the given index
func (jn Arr) Get(index int) Any {
    if jn.first(index, 0) || jn.last(index, 0) {
        return Any{}
    }
    a := jn.rawValue
    switch a.(type) {
    case []interface{}:
        return Any{rawValue: a.([]interface{})[index]}
    case []map[string]interface{}:
        return Any{rawValue: a.([]map[string]interface{})[index]}
    default:
        panic("Unable to cast the jnav to an array")
    }
}

// Resets the position of the iterator on the BOM
func (jn *Arr) Rewind() {
    jn.pos = -1
}

// Builds a jnav from a go map
func FromMap(node map[string]interface{}) Obj {
    return Obj{Any{node}}
}

// Builds a jnav from an go array
func FromArr(node []interface{}) Arr {
    return Arr{
        Any: Any{node},
        pos: -1,
    }
}

// Translate a generic element to a JNav Array
func (jn Any) AsArr() Arr {
    return Arr{Any: Any{jn.rawValue}, pos: -1}
}

// Translate a generic element to a JNav Object
func (jn Any) AsObj() Obj {
    if casted, ok := jn.rawValue.(map[string]interface{}); ok {
        return Obj{Any: Any{casted}}
    } else {
        return Obj{}
    }
}

// Translates a generic element to a string
func (jn Any) AsString() (string, bool) {
    res, ok := jn.rawValue.(string)
    return res, ok
}

// Translate a generic element to a string
func (jn Any) ToString() string {
    return fmt.Sprintf("%v", jn.rawValue)
}

// Translate a generic element to a string or returns a fallback if not possible
func (jn Any) AsStringOr(fallback string) string {
    if res, ok := jn.rawValue.(string); ok {
        return res
    } else {
        return fallback
    }
}

// Translate a generic element to am int or returns
func (jn Any) AsInt() (int, bool) {
    res, ok := jn.rawValue.(int)
    return res, ok
}

// Translate a generic element to am int or returns a fallback if not possible
func (jn Any) AsIntOr(fallback int) int {
    if res, ok := jn.rawValue.(int); ok {
        return res
    } else {
        return fallback
    }
}

// Gets an the value of an object given the key
func (jn Obj) Get(key string) Any {
    if jn.rawValue == nil {
        return Any{}
    } else {
        return Any{rawValue: jn.rawValue.(map[string]interface{})[key]}
    }
}
