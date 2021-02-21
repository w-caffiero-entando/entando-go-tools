package sugar

import "entando_go_tools/pkg/ent/misc"

type SweetResult struct {
    value interface{}
    err   error
}

func (sweetResult SweetResult) String() (string, error) {
    return sweetResult.value.(string), sweetResult.err
}

func (sweetResult SweetResult) Err() error {
    return sweetResult.err
}

func (sweetResult SweetResult) WrappedErr(format string, args ...interface{}) error {
    return misc.WrapError(sweetResult.err, format, args...)
}

func (sweetResult SweetResult) StringOrElse(fallback string) string {
    if sweetResult.err == nil {
        return sweetResult.value.(string)
    } else {
        return fallback
    }
}

func (sweetResult SweetResult) StringOrFail(failureDescription string) string {
    if sweetResult.err == nil {
        return sweetResult.value.(string)
    } else {
        FATAL("%s", failureDescription)
        return ""
    }
}

func SomeSweet(value string) SweetResult {
    return SweetResult{value, nil}
}

func NoSweet() SweetResult {
    return SweetResult{"", nil}
}

func NoSweetErr(err error) SweetResult {
    return SweetResult{"", err}
}
