package subcmds

import (
    "fmt"
    "os"
    "strings"

    . "entando_go_tools/pkg/util/sugar"

    "entando_go_tools/pkg/util/validators"

    "entando_go_tools/pkg/ent/misc"

    "entando_go_tools/pkg/ent/entapp"

    "entando_go_tools/pkg/support/kube"
    "entando_go_tools/pkg/util/sys/spawn"

    "github.com/spf13/cobra"

    . "entando_go_tools/pkg/util/log/def"

    "entando_go_tools/pkg/ent/maven/pom"
)

func init() {
    moduleCmd.AddCommand(&cobra.Command{
        Use:   "purge",
        Short: "uninstalls a bundle and removes all its objects",
        Long:  `uninstalls a bundle and removes all its objects`,
        Run: func(cmd *cobra.Command, args []string) {
            _ = os.Chdir("/home/wrt/work/prj/entando/_test_area/external-test-area/ent/myBundle/testdir")
            ecrPurge()
        },
    })
    ParsePurgeParams(moduleCmd)
}

func ecrPurge() {
    dockerImage := pom.ExtractDockerImageLocation()
    baseBundleName := dockerImage.ImageOrg + "-" + dockerImage.ImageName
    var normalizedBundleName = baseBundleName
    if len(normalizedBundleName) > 31 {
        normalizedBundleName = normalizedBundleName[0:31]
    }
    // Delete the main resources
    delResourcesByStrategy(ByName, "EntandoPlugin", "", normalizedBundleName)
    delResourcesByStrategy(ByLabel, "deployment", "", normalizedBundleName)
    delResourcesByStrategy(ByLabel, "pod", "(main)", normalizedBundleName)
    // Delete the related resources
    if PurgeParams.Related {
        BizLog.Info("Purging also the related resources as requested")
        delResourcesByStrategy(ByLinkLabel, "pod", "(aux)", string(entapp.Params.TheAppname)+"-"+normalizedBundleName+"-link")
        delResourcesByStrategy(ByLabel, "service", "", normalizedBundleName)
        delResourcesByStrategy(ByLabel, "ingress", "", normalizedBundleName)
        delResourcesByStrategy(ByLabel, "secret", "(main)", normalizedBundleName)
        delResourcesByStrategy(ByNameFilters, "secret", "(aux)",
            normalizedBundleName,
            "^"+baseBundleName+"-server-secret$",
            "^"+baseBundleName+"-sidecar-secret$",
            "^"+string(entapp.Params.TheAppname)+"-"+baseBundleName+"-link-controller-ca-cert-secret$")
    }
    // Delete the volumes resources
    if PurgeParams.Volumes {
        BizLog.Info("Purging also the volumes as requested")
        delResourcesByStrategy(ByLabel, "pvc", "", normalizedBundleName)
    }
    //~
    BizLog.Info("done.")
}

func delResourcesByStrategy(strategy int, resourceType string, desc string, resourceSelector ...string) {
    var err error

    BizLog.Info("Deleting the %s(s) %s", resourceType, desc)
    resourcePattern := resourceSelector[0]

    switch strategy {
    case ByName:
        err = delResourceDirectly(resourceType, []interface{}{resourcePattern})
    case ByLabel:
        err = delResourceDirectly(resourceType, []interface{}{"-l", "EntandoPlugin=" + resourcePattern})
    case ByLinkLabel:
        err = delResourceDirectly(resourceType, []interface{}{"-l", "EntandoAppPluginLink=" + resourcePattern})
    case ByNameFilters:
        err = delResourceByNameFilters(resourceType, resourceSelector)
    default:
        validators.FATAL("Unknown mode %d", strategy)
    }

    if err == skipErr || err == notFound {
        BizLog.Info("%s", err.Error())
        return
    }
    if err != nil {
        BizLog.Error("error deleting %s: %s", resourceType, err.Error())
        return
    }

    BizLog.Error("resource %s deleted", resourceType)
}

// Deletes the resource obtained by filtering, with the given filters, the downloaded list of resource
//
func delResourceByNameFilters(resourceType string, regexpFilters []string) error {
    // Downloads the list by type
    res, err := kube.Kubectl(nil,
        "get",
        resourceType,
        "--no-headers",
        spawn.SOA("o", "custom-columns", "NAME:.metadata.name"),
    )

    if err != nil {
        return misc.WrapError(err, "deleting %s resource(s)", resourceType)
    }

    // Filter the list
    resourcesToDelete := misc.FilerViaMultiRegexp(
        misc.PortableLineSplit(res.Stdout),
        regexpFilters...,
    )

    if len(resourcesToDelete) == 0 {
        return notFound
    }

    BizLog.Debug("deleting %s(s): %s", resourceType, strings.Join(resourcesToDelete, ", "))

    if PurgeParams.DryRun {
        return skipErr
    }

    // Deletes
    return SweetRes(

        kube.Kubectl(nil, "delete", resourceType, resourcesToDelete),

    ).WrappedErr("deleting %s resource(s)", resourceType)
}

// Deletes resources by using a native kubernetes selection method
//
func delResourceDirectly(resourceType string, selector []interface{}) error {
    res, err := kube.Kubectl(nil,
        "get",
        resourceType,
        selector,
        "--no-headers",
        spawn.SOA("o", "custom-columns", "NAME:.metadata.name"),
    )

    if err != nil {
        return misc.WrapError(err, "deleting %s resource(s)", resourceType)
    }

    rawListOfResourcesToDelete := strings.TrimSpace(res.Stdout)
    if rawListOfResourcesToDelete == "" {
        return notFound
    }

    resourcesToDelete := misc.PortableLineSplit(rawListOfResourcesToDelete)
    BizLog.Debug("deleting %s(s): %s", resourceType, strings.Join(resourcesToDelete, ", "))

    if PurgeParams.DryRun {
        return skipErr
    }

    if resourceType == "pvc" {
        _, err = kube.Kubectl(nil,
            "patch",
            resourceType,
            resourcesToDelete,
            "--no-headers",
            "p", `{"metadata":{"finalizers":null}}`,
        )
        if err != nil {
            return misc.WrapError(err, "deleting %s resource(s)", resourceType)
        }
    }

    _, err = kube.Kubectl(nil, "delete", resourceType, resourcesToDelete)

    return misc.WrapError(err, "deleting %s resource(s)", resourceType)
}

func ParsePurgeParams(cmd *cobra.Command) {
    cmd.PersistentFlags().BoolVar(&PurgeParams.Related, "purge-related", false, "version to install")
    cmd.PersistentFlags().BoolVar(&PurgeParams.Volumes, "purge-volumes", false, "version to install")
    cmd.PersistentFlags().BoolVar(&PurgeParams.DryRun, "dry-run", false, "version to install")
}

type PurgeParamsType struct {
    Related bool
    Volumes bool
    DryRun  bool
}

var PurgeParams PurgeParamsType

var skipErr = fmt.Errorf("skipped due to dry-run")
var notFound = fmt.Errorf("no suitable resource found")

const (
    ByName        = iota
    ByLabel       = iota
    ByLinkLabel   = iota
    ByNameFilters = iota
)
