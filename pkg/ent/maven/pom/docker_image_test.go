package pom

import (
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExtractFqImageName(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }

    _ = os.Chdir("../../../_test_resources")
    res := ExtractDockerImageLocation()
    assert.Equal(t, res.ImageFqname, "organization/my-bundle:0.0.1-SNAPSHOT")
}
