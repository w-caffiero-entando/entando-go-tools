package pom

import (
    "os/exec"
    "strings"
    "time"

    "entando_go_tools/pkg/util/xmlnav"

    "entando_go_tools/pkg/util/sys/spawn"

    "entando_go_tools/pkg/support/maven"

    "entando_go_tools/pkg/util/sys"
)

type DockerImageInfo struct {
    ImageName   string
    ImageTag    string
    ImageOrg    string
    ImageFqname string
}

func ExtractDockerImageLocation() DockerImageInfo {
    ret := DockerImageInfo{}
    gocmd := &exec.Cmd{}

    mvnRes, _ := sys.ExecLongRunningTask("the POM analysis", 5*time.Second, 60*time.Second,
        func() (interface{}, error) {
            return maven.Maven(&gocmd,
                "-q",
                "--non-recursive",
                "exec:exec",
                spawn.SOA("D", "exec.executable", "echo"),
                spawn.SOA("D", "exec.args", "${project.artifactId}:${project.version}"),
            )
        }, func() {
            _ = gocmd.Process.Kill()
        },
    )

    var mvnResArr = strings.Split(mvnRes.Result.(spawn.Res).Stdout, ":")

    ret.ImageName = strings.TrimSuffix(mvnResArr[0], "\n")
    ret.ImageTag = strings.TrimSuffix(mvnResArr[1], "\n")

    node := &xmlnav.XmlNode{}
    _ = node.LoadFile("pom.xml")
    node = node.Nav(
        "build/pluginManagement/plugins/plugin/configuration[artifactId=jib-maven-plugin]/to/image",
    )

    imageTemplate := string(node.Content)
    ret.ImageFqname = strings.Replace(imageTemplate, "${project.artifactId}", ret.ImageName, -1)
    ret.ImageFqname = strings.Replace(ret.ImageFqname, "${project.version}", ret.ImageTag, -1)
    ret.ImageOrg = strings.Split(ret.ImageFqname, "/")[0]
    return ret
}
