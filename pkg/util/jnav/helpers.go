package jnav

import (
    "fmt"
    "strings"
)

// Prints a table given a JNAV object as source
//  - arg           the jnav source, can be an Arr or an Obj
//  - fieldsPath    the path of the fields to print
//
func PrintTable(arg interface{}, fieldsPath ...string) {
    fieldMaxLen := [100]int{}
    rowOut := ""
    var arr Arr
    switch arg.(type) {
    case Arr:
        arr = arg.(Arr)
    case Obj:
        arr = Arr{
            Any: Any{[]interface{}{arg.(Obj).rawValue}},
            pos: -1,
        }

    }
    // Collect
    for _, c := range []int{1, 2} {
        arr.Rewind()
        for arr.Next() {
            row := arr.Current().AsObj()
            if c == 2 {
                rowOut = "| "
            }
            for i, fp := range fieldsPath {
                value := row.Nav(fp).ToString()
                valueLen := len(value)
                valueMaxLen := fieldMaxLen[i]
                if c == 1 {
                    if valueLen > fieldMaxLen[i] {
                        fieldMaxLen[i] = valueLen
                    }
                } else {
                    rowOut += fmt.Sprintf("%-*s | ", valueMaxLen, value)
                }
            }
            if c == 2 {
                fmt.Println(strings.TrimSuffix(rowOut, " "))
            }
        }
    }
}

func MkObjValueReturn(jnav Obj, key string, err error) (string, error) {
    if err != nil {
        return "", err
    }

    if ret, ok := jnav.Get(key).AsString(); ok {
        return ret, nil
    } else {
        return "", fmt.Errorf("error extracting key \"%s\"", key)
    }
}
