package entapp

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "entando_go_tools/pkg/support/kube"
)

func TestGetIngressesIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }

    ParseBaseTestExecutionParams()
    kube.ParseBaseTestExecutionParams()

    t.Run("QueryIngressesInfo", func(t *testing.T) {
        ii := QueryIngressesInfo()
        assert.NotEmpty(t, ii.Main)
    })
}
