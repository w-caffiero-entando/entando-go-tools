package validators

import (
    "fmt"
    "regexp"
    "strings"

    . "entando_go_tools/pkg/util/log/def"

    "entando_go_tools/pkg/util/sugar"
)

// Literal object asserter
func assert_lit(expected, found interface{}, desc string) bool {
    if desc != "" {
        desc = " for " + desc
    }
    if expected != found {
        FATAL("Expected literal %v found %v%s", expected, found, desc)
    }
    return true
}

// Identifier
func Assert_Id(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-z][a-zA-Z0-9_]*$", sugar.FNN(desc, "identifier"), options)
}

// Like assert_id but case is ignored
func Assert_Ic_Id(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-zA-Z0-9_]*$", sugar.FNN(desc, "case insensitive identifier"), options)
}

// Extended Identifier
func Assert_Ext_Id(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-z][a-zA-Z0-9_-]*$", sugar.FNN(desc, "extended identifier"), options)
}

// Extended Identifier ignore case
func Assert_Ext_Ic_Id(varName, value, desc string, allowNulls bool, options string) bool {
    if value == "" && allowNulls {
        return true
    }
    return assertRegexNN(varName, value, "^[a-zA-Z0-9_-]*$",
        sugar.FNN(desc, "extended case insensitive identifier"), options)
}

// Extended Identifier ignore case with spaces
func Assert_Ext_Ic_Id_Spc(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-zA-Z0-9 _-]*$",
        sugar.FNN(desc, "extended identifier with spaces (ignore case)"), options)
}

// Identifier with spaces
func Assert_Id_Spc(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-zA-Z0-9_ ]*$", sugar.FNN(desc, "identifier with spaces"), options)
}

// Identifier with spaces
func Assert_Num(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[0-9]*$", sugar.FNN(desc, "number"), options)
}

// Strinct Identifier with spaces
func Assert_StrictId(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-z][a-zA-Z0-9]*$", sugar.FNN(desc, "strict identifier"), options)
}

// Single Domain Name
func Assert_DN(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, "^[a-z][a-z0-9_-]*$", sugar.FNN(desc, "single domain name"), options)
}

// Full Domain Name
func Assert_FDN(varName, value, desc string, options string) bool {
    regexps := []string{`\.`, `^([a-z0-9._-])+$`}
    return assertRegexNN(varName, value, regexps[0], sugar.FNN(desc, "full domain name"), options) &&
        assertRegexNN(varName, value, regexps[1], sugar.FNN(desc, "full domain name"), options)
}

// Single Domain Name
func Assert_URL(varName, value, desc string, options string) bool {
    regex := "^(https?|file)://[-A-Za-z0-9\\+&@///%?=~_|!:,.;]*[-A-Za-z0-9\\+&@///%=~_|]"
    return assertRegexNN(varName, value, regex, sugar.FNN(desc, "url"), options)
}

// Single Domain Name
func Assert_EMAIL(varName, value, desc string, options string) bool {
    desc = sugar.FNN(desc, "email")
    return assertRegexNN(varName, value, `^[^.]+.*$`, desc, options) &&
        assertMultiRegexNN(varName, value, []string{
            `^[a-z0-9._-]+@[a-z0-9._-]+\.[a-z0-9._-]+$`,
            `^[a-z0-9._-]+@[a-z0-9._-]+\.[a-z0-9._-]+\.[a-z0-9._-]+$`,
        }, desc, options) &&
        assertRegexNN(varName, value, `^[a-z0-9._-]+@[a-z0-9._-]+\.[a-z0-9._-]+$`, desc, options)
}

// Semantic Version
func Assert_Ver(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, `^v?[0-9.]+-?[a-zA-Z0-9-]+$`, sugar.FNN(desc, "semver"), options)
}

// IP
func Assert_IP(varName, value, desc string, options string) bool {
    return assertRegexNN(
        varName, value, `^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`, sugar.FNN(desc, "ip address"), options)
}

// GIGA
func Assert_Giga(varName, value, desc string, options string) bool {
    return assertRegexNN(varName, value, `^[0-9]*G$`, sugar.FNN(desc, "gigabyte value"), options)
}

func assertRegexNN(
    varName string,
    value string,
    regex string,
    typeDescription string,
    options string,
) bool {
    var fatal = strings.Contains(options, "fatal")
    var silent = strings.Contains(options, "silent")

    var valid = regexp.MustCompile(regex)
    if valid.MatchString(value) {
        return true
    }

    if fatal {
        FATAL("Value of %s (%v) is not a valid %s", varName, value, typeDescription)
    } else if ! silent {
        GenLog.Error("Value of %s (%v) is not a valid %s", varName, value, typeDescription)
        return false
    }
    return false
}

func assertMultiRegexNN(varName string, value string, regexps []string, desc string, options string) bool {
    last := len(regexps) - 1
    for _, r := range regexps[0:last] {
        if assertRegexNN(varName, value, r, desc, "silent") {
            return true
        }
    }
    return assertRegexNN(varName, value, regexps[last], desc, options)
}

func FATAL(format string, args ...interface{}) {
    GenLog.Error(format, args...)
    panic(fmt.Sprintf(format, args...))
}

func FATAL_IF(cond bool, format string, args ...interface{}) {
    if cond {
        FATAL(format, args...)
    }
}
