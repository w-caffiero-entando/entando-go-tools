package config

import (
    "fmt"
    "io/ioutil"
    "net/url"
    "os"
    "strings"

    "entando_go_tools/pkg/ent/misc"
    . "entando_go_tools/pkg/util/log/def"
    "entando_go_tools/pkg/util/validators"
)

func SaveValue(key, value string) error {
    return saveValue(key, value, Params.DbFile)
}
func LoadValue(key string) (string, error) {
    return loadValue(key, Params.DbFile)
}

type Config map[string]string

// Loads a configuration file.
//
// Configuration files syntax rules:
//  - Line format is key=value
//  - Keys must be standard identifiers
//  - IN Values the EOL chars are escaped via url.PathEscape
//  - Comments are possibles using "#" are first line character
//  - Blank/Space lines are ignored
//  - Order is not ensured
//
func LoadFile(filename string) (Config, error) {
    ret := Config{}
    wrapErr := func(err error) error {
        return fmt.Errorf("saving a config file: %w", err)
    }

    file, err := os.Open(filename)
    if err != nil {
        return Config{}, wrapErr(err)
    }
    data, err := ioutil.ReadAll(file)
    if err != nil {
        return Config{}, wrapErr(err)
    }
    errorsCount := 0
    for i, line := range misc.PortableLineSplit(string(data)) {
        if strings.HasPrefix(line, "#") {
            continue
        }
        if strings.TrimSpace(line) == "" {
            continue
        }
        pos := strings.IndexByte(line, '=')
        if pos == -1 {
            SysLog.Warn("line #%d ignored due to bad format (file \"%s\")", i+1, filename)
            errorsCount++
            continue
        }
        k := line[0:pos]
        if isKeyValid(k) == false {
            SysLog.Warn("line #%d ignored due to bad key format (file \"%s\")", i+1, filename)
            errorsCount++
            continue
        }
        v, _ := url.PathUnescape(line[pos+1:])
        ret[k] = v
    }
    if errorsCount > 0 && len(ret) == 0 {
        err = fmt.Errorf("file \"%s\" doesn't seem to contain valid config lines", filename)
        ss := err.Error()
        SysLog.Warn("%s", ss)
        return Config{}, wrapErr(err)
    }
    return ret, nil
}

func SaveFile(filename string, cfg Config) error {
    content := ""
    for k, v := range cfg {
        content += k + "=" + url.PathEscape(v) + "\n"
    }
    err := ioutil.WriteFile(filename, []byte(content), 0644)
    if err != nil {
        err = fmt.Errorf("saving the config file: %w", err)
        SysLog.Error(err.Error(), filename)
        return err
    }
    return nil
}

// PRIVATE
func loadValue(key string, filename string) (string, error) {
    if validators.Assert_Ext_Ic_Id("configuration key", key, "configuration key", false, "silent") == false {
        return "", fmt.Errorf("loading a config value: illegal key \"%s\" detected", key)
    }
    cfg, err := LoadFile(filename)
    return cfg[key], err
}

func saveValue(key, value string, filename string) error {
    if isKeyValid(key) == false {
        return fmt.Errorf("saving a config value: illegal key \"%s\" detected", key)
    }
    cfg := Config{}
    if _, err := os.Stat(filename); err == nil {
        cfg, err = LoadFile(filename)
        if err != nil {
            return fmt.Errorf("saving a config value (refresh phase): %w", err)
        }
    }
    cfg[key] = value
    return SaveFile(filename, cfg)
}

func isKeyValid(key string) bool {
    return validators.Assert_Ext_Ic_Id("configuration key", key, "configuration key", false, "silent")
}
