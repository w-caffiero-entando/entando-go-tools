package config

import (
    "io/ioutil"
    "log"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"

    "entando_go_tools/pkg/util/sys"

    "entando_go_tools/pkg/ent/misc"

    . "entando_go_tools/pkg/util/sugar"
)

func TestLoadSaveCycle(t *testing.T) {
    prepare()
    defer func() { _ = os.Remove(Params.DbFile) }()

    var err error

    //~
    err = SaveValue("KEY", "VALUE")
    assert.Nil(t, err)
    value, err := LoadValue("KEY")
    assert.Nil(t, err)
    assert.Equal(t, "VALUE", value)
    //~
    err = SaveValue("KEY", "VALUE2")
    assert.Nil(t, err)
    value, err = LoadValue("KEY")
    assert.Nil(t, err)
    assert.Equal(t, "VALUE2", value)
    //~
    err = SaveValue("KEYX", "VALUEX")
    assert.Nil(t, err)
    value, err = LoadValue("KEY")
    assert.Nil(t, err)
    assert.Equal(t, "VALUE2", value)
    value, err = LoadValue("KEYX")
    assert.Nil(t, err)
    //~
    err = SaveValue("KEYLF", "VALUE\nVALUE")
    assert.Nil(t, err)
    value, err = LoadValue("KEYLF")
    assert.Nil(t, err)
    assert.Equal(t, "VALUE\nVALUE", value)
    //~
    err = SaveValue("KEYLF", "VALUE=VALUE")
    assert.Nil(t, err)
    value, err = LoadValue("KEYLF")
    assert.Nil(t, err)
    assert.Equal(t, "VALUE=VALUE", value)
    //~
}

func TestErrorAndSpecialCases(t *testing.T) {
    var tmp SweetResult

    sys.CaptureGoOutput(true, true, func() {
        // Bath config path
        badConfig := "/tmp/an-invalid-invalid-invalid-invalid-invalid-invalid-invalid-path/config"
        assert.NotNil(t, saveValue("key", "value", badConfig))
        assert.NotNil(t, SweetRes(loadValue("key", Params.DbFile)).Err())

        // Bad key
        badKey := "#KEY"
        assert.NotNil(t, saveValue(badKey, "value", badConfig))
        assert.NotNil(t, SweetRes(loadValue(badKey, "value")).Err())

        // Inaccessible files
        prepare()
        cleanup := func() {
            _ = os.Chmod(Params.DbFile, 0644)
            _ = os.Remove(Params.DbFile)
        }

        defer cleanup()

        _ = os.Chmod(Params.DbFile, 0000)
        assert.NotNil(t, saveValue("key", "value", badConfig))
        assert.NotNil(t, SweetRes(loadValue("key", Params.DbFile)).Err())
        cleanup()

        // Comment lines
        prepare()
        cleanup = func() { _ = os.Remove(Params.DbFile) }
        defer cleanup()

        _ = misc.AppendToFile(Params.DbFile, "#KEY=SOME_VALUE\n")
        assert.Equal(t, "", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))
        cleanup()
        assert.Nil(t, saveValue("KEY", "VALUE", Params.DbFile))
        _ = misc.AppendToFile(Params.DbFile, "#KEY=SOME_VALUE\n")
        assert.Equal(t, "VALUE", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))

        // Tolerable Corruption
        prepare()
        cleanup = func() { _ = os.Remove(Params.DbFile) }
        defer cleanup()

        _ = misc.AppendToFile(Params.DbFile, "@KEY=SOME_VALUE\n")
        _ = misc.AppendToFile(Params.DbFile, "CORRUPTION$/(KEY=VALUEUU£KNLDNSUC\n")
        _ = misc.AppendToFile(Params.DbFile, "AKEY=AVALUE\n")
        assert.Equal(t, "", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))
        assert.Nil(t, saveValue("KEY", "VALUE", Params.DbFile))
        _ = misc.AppendToFile(Params.DbFile, "CORRUPTION£/(KEY=VALUEUU£KNLDNSUC\n")
        tmp = SweetRes(LoadFile(Params.DbFile))
        assert.Nil(t, tmp.Err())
        assert.Equal(t, "VALUE", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))
        cleanup()

        // Intolerable Corruption
        //  - File cannot be loaded anymore
        //  - File cannot be updated anymore
        prepare()
        cleanup = func() { _ = os.Remove(Params.DbFile) }
        defer cleanup()

        _ = misc.AppendToFile(Params.DbFile, "@KEY=SOME_VALUE\n")
        _ = misc.AppendToFile(Params.DbFile, "CORRUPTION$/(KEY=VALUEUU£KNLDNSUC\n")
        assert.Equal(t, "", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))
        assert.NotNil(t, saveValue("KEY", "VALUE", Params.DbFile))
        tmp = SweetRes(LoadFile(Params.DbFile))
        assert.NotNil(t, tmp.Err())
        assert.Equal(t, "", SweetRes(loadValue("KEY", Params.DbFile)).StringOrElse(""))
        cleanup()
    })
}

// utils
func prepare() {
    testConfigFile, err := ioutil.TempFile("", "config-test")
    if err != nil {
        log.Fatal(err)
    }
    Params.DbFile = testConfigFile.Name()
}
