package config

import (
    "os"

    "github.com/spf13/cobra"

    "entando_go_tools/pkg/util/validators"
)

func ParseParams(cmd *cobra.Command) {
    tmp := []*string{
        cmd.PersistentFlags().String("config", "", "kubeconfig file"),
        cmd.PersistentFlags().String("config-global", "", "global configuration directory"),
    }

    cobra.OnInitialize(func() {
        // Global Config
        userHomeDir, err := os.UserHomeDir()
        if err != nil {
            validators.FATAL("Unable to determine the user home dir due to error \"%s\"", err.Error())
        }
        globalConfigFile := *tmp[1]
        if globalConfigFile == "" {
            globalConfigFile = userHomeDir + ".entando/.global-cfg"
        }
        Params.GlobalDbFile = globalConfigFile
        // Config
        configFile := *tmp[1]
        if configFile == "" {
            config, _ := LoadFile(globalConfigFile)
            configFile = config["DESIGNATED_PROFILE_HOME"]
        }
        if configFile == "" {
            validators.FATAL("Unable to determine config file")
        }
        Params.DbFile = configFile
    })
}

func ParseTestParams() {
    Params = ParamsType{
        GlobalDbFile: os.Getenv("TEST_globalDbFile"),
        DbFile:       os.Getenv("TEST_dbFile"),
    }
}

type ParamsType struct {
    GlobalDbFile string
    DbFile       string
}

var Params ParamsType
