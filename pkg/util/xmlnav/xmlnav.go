package xmlnav

import (
    "bytes"
    "encoding/xml"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
)

const (
    EnterNode  = iota
    NextNode   = iota
    RepeatNode = iota
    ReturnNode = iota
    ExitNode   = iota
)

type XmlNode struct {
    XMLName xml.Name
    Attrs   []xml.Attr `xml:"-"`
    Content []byte     `xml:",innerxml"`
    Nodes   []XmlNode  `xml:",any"`
}

// Unmarshal a raw XML document
func (node *XmlNode) LoadXmlData(data []byte) error {
    buf := bytes.NewBuffer(data)
    dec := xml.NewDecoder(buf)
    return dec.Decode(&node)
}

// Unmarshal an XML file
func (node *XmlNode) LoadFile(filename string) error {
    file, _ := os.Open(filename)
    data, _ := ioutil.ReadAll(file)
    return node.LoadXmlData(data)
}

// Navigates to a specific node through a simple query string with these rules:
//  - It's a slash-delimited sequence of nodes names that indicates the path to take
//  - It supports a special syntax [*1] to select a node using the value of a sibling
//  - The root tag should not be included in the query
//
//  *1: node[sibling=value]  => selects the node that has a sibling of name "value"
//
func (node *XmlNode) Nav(expr string) *XmlNode {
    find := strings.Split(expr, "/")
    return Walk(node.Nodes, 0, func(n XmlNode, level int, repeatNode bool) int {
        tag := n.XMLName.Local
        findExpr := find[level]
        match := siblingSelectorRegex.FindStringSubmatch(findExpr)

        if len(match) == 0 {
            if tag == findExpr {
                return ret(level, find, EnterNode)
            }
        } else {
            var baseTag = match[1]
            var condTag = match[2]
            var condTagValue = match[3]

            if repeatNode {
                if tag == baseTag {
                    return ret(level, find, EnterNode)
                }
            } else
            if tag == condTag {
                if string(n.Content) == condTagValue {
                    return RepeatNode
                } else {
                    return ExitNode
                }
            }
        }
        return NextNode
    })
}

var siblingSelectorRegex = regexp.MustCompile(`^([^|]+)\[([^]=]+)=(.+)]$`)

func ret(level int, find []string, direction int) int {
    if level < len(find)-1 {
        return direction
    } else {
        return ReturnNode
    }
}

// Walks through a tree of XmlNode(s) according to the directions provided by the user function
//
// The possible directions are:
//  - EnterNode     walks into the current node
//  - ReturnNode    stops the walk and returns the current node
//  - RepeatNode    walks though all the sub-nodes of the current node
//  - ExitNode      goes back to the parent node and then walks to its next sibling
//
func Walk(nodes []XmlNode, level int, getDirections func(XmlNode, int, bool) int) *XmlNode {
    repeat := false
    for {
        for _, currentNode := range nodes {
            switch getDirections(currentNode, level, repeat) {
            case EnterNode:
                if ret := Walk(currentNode.Nodes, level+1, getDirections); ret != nil {
                    return ret
                }
            case ReturnNode:
                return &currentNode
            case RepeatNode:
                repeat = true
                break
            case ExitNode:
                return nil
            case NextNode:
            }
        }
        if !repeat {
            break
        }
    }
    return nil
}
