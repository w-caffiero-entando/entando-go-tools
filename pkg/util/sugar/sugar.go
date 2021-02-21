package sugar

import (
    "fmt"
)

var FATAL = func(format string, args ...interface{}) {
    panic(fmt.Sprintf(format, args...))
}

//  First Non Null/Empty
func FNN(args ...string) string {
    for _, a := range args {
        if a != "" {
            return a
        }
    }
    return ""
}

type SweetArr []string

func (arr SweetArr) At(index int) SweetResult {
    if len(arr) < index+1 {
        return NoSweet()
    }
    return SomeSweet(arr[index])
}

func SweetRes(value interface{}, err error) SweetResult {
    if err != nil {
        return NoSweetErr(err)
    }
    return SweetResult{value:value}
}
