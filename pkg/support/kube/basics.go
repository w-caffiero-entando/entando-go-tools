package kube

type KNamespace string

const AllNamespaces = KNamespace("*")
const NoNamespace = KNamespace("")

type KConfig string

const DefaultKubeconfig = KConfig("")

type KContext string

const DefaultKubectx = KConfig("")

type KKubectlCommand string

const DefaultCommand = KConfig("kubectl")

type Header struct {
    Name  string
    Value string
}
