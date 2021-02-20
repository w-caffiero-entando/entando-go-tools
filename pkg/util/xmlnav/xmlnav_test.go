package xmlnav

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "entando_go_tools/pkg/util/sys"
)

func TestFind(t *testing.T) {
    node := &XmlNode{}

    _ = node.LoadFile(sys.DetermineCurrentSourceFileDir() + "/example-pom.xml")

    node = node.Nav(
        "build",
    )

    assert.Equal(t, 3, len(node.Nodes))
    assert.Equal(t, "build", node.XMLName.Local)

    node = node.Nav(
        "pluginManagement/plugins/plugin/configuration[artifactId=jib-maven-plugin]/to/image",
    )

    assert.Equal(t, "organization/${project.artifactId}:${project.version}", string(node.Content))
}

func TestFindNothing(t *testing.T) {
    node := &XmlNode{}

    _ = node.LoadFile(sys.DetermineCurrentSourceFileDir() + "/example-pom.xml")

    node = node.Nav(
        "something-not-present",
    )

    assert.Nil(t, node)
}
