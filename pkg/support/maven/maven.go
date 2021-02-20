package maven

import (
    "os/exec"

    "entando_go_tools/pkg/util/sys/spawn"

    "entando_go_tools/pkg/util/sys"
)

func Maven(gocmd **exec.Cmd, args ...interface{}) (spawn.Res, error) {
    cmd := "mvn"
    if sys.IsCommandAvailable("./mvnw") {
        cmd = "./mvnw"
    }
    return spawn.Spawn(gocmd,
        cmd,
        spawn.ComposeArgs(args),
        spawn.Environ{},
        spawn.Options{
            CaptureStdout: true,
            CaptureStderr: false,
            Interactive: true,
            WithSudo: false,
        },
    )
}
