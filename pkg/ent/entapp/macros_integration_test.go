package entapp

import (
    "testing"

    "github.com/stretchr/testify/assert"

    "entando_go_tools/pkg/support/kube"
)

func TestPrepareResourceActionIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }

    ParseBaseTestExecutionParams()
    kube.ParseBaseTestExecutionParams()

    t.Run("QueryIngressesInfo", func(t *testing.T) {
        resp, _ := ResourceAction{
            Verb:         "GET",
            ResourceName: "components",
        }.Authenticate().Exec()
        v, _ := resp.Nav("payload.[0].code").AsString()
        assert.Equal(t, "entando-app-engine-v7-ecommerce-demo-bundle", v)
    })
}
