package misc

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "regexp"
    "strings"
)

// Cleanups and normalizes a path.
//  - If present removes the trailing EOL
//  - Ensures that the path finishes with a slash
//
// #unstable
func NormalizeBasePath(path string) string {
    path = strings.TrimSuffix(path, "\r")
    path = strings.TrimSuffix(path, "\n")
    path = strings.TrimSuffix(path, "\r")
    if strings.HasSuffix(path, "/") {
        return path
    } else {
        return path + "/"
    }
}

// Converts a space delimited text into a string->string map
//
// #unstable
func FieldsToMap(text string) []map[string]string {
    rows := PortableLineSplit(text)
    var ret []map[string]string
    var header []string
    for _, row := range rows {
        if header == nil {
            header = strings.Fields(row)
        } else {
            rec := strings.Fields(row)
            if len(rec) != len(header) {
                break
            }
            dict := map[string]string{}
            for i := range header {
                dict[header[i]] = rec[i]
            }
            ret = append(ret, dict)
        }
    }
    return ret
}

// Sprint a file in line by checking multiple EOL formats.
//
// #unstable
func PortableLineSplit(text string) []string {
    text = strings.Replace(text, "\r\n", "\n", -1)
    text = strings.Replace(text, "\n\r", "\n", -1)
    text = strings.Replace(text, "\r", "\n", -1)
    return strings.Split(text, "\n")
}

// Calculates the checksum of a file
//
// #unstable
func CalculateFileChecksum(file string) string {
    h := sha256.New()
    s, err := ioutil.ReadFile(file)
    h.Write(s)
    if err != nil {
        log.Fatal(err)
    }
    return hex.EncodeToString(h.Sum(nil))
}

// Appends a string content to a file
//
// #unstable
func AppendToFile(filename string, content string) error {
    return AppendToFileWithPerm(filename, content, 0644)
}

// Appends a string content to a file
// given also the file permission in case the files was not present
//
// #unstable
func AppendToFileWithPerm(filename string, content string, perm os.FileMode) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, perm)
    if err != nil {
        return err
    }
    defer f.Close()
    if _, err := f.WriteString(content); err != nil {
        return err
    }
    return nil
}

// Filter the candidates returning only the ones the have a match into the regexSelectors
//
// You can thing about it almost like a regexp inner join
//
// #unstable
func FilerViaMultiRegexp(candidates []string, regexSelectors ...string) (ret []string) {
    for _, selector := range regexSelectors {
        compiledSelector := regexp.MustCompile(selector)
        found := StringArrayFind(candidates, func(i int, e string) bool {
            return compiledSelector.MatchString(e)
        })
        if found != -1 {
            ret = append(ret, candidates[found])
        }
    }
    return ret
}

// Returns a wrapped error or nil if error is nil
//
// #unstable
func WrapError(err error, format string, args ...interface{}) error {
    if err == nil {
        return nil
    }
    msg := fmt.Sprintf(format, args...)
    //noinspection GoPrintFunctions
    return fmt.Errorf("%s: %w", msg, err)
}

// Returns the index of an item in an array
//
// #unstable
func StringArrayFind(arr []string, matcher func(i int, e string) bool) int {
    for i, e := range arr {
        if matcher(i, e) {
            return i
        }
    }
    return -1
}
