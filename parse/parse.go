package parse

import (
    "fmt"
    "os"
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
)

// Recursive-descent parser

var pos int = 0
var tokens []interface{}

func new_node(op int, lhs *Node, rhs *Node) *Node {
    return &Node{Ty: op, Lhs: lhs, Rhs: rhs}
}

func new_node_num(val int) *Node {
    return &Node{Ty: ND_NUM, Val: val}
}

func fail(i int) {
    t := tokens[pos].(*Token)
    fmt.Printf("unexpected token: %s\n", t.Input)
    os.Exit(1)
}

func num() *Node {
    t := tokens[pos].(*Token)
    if t.Ty == TK_NUM {
        pos++
        return new_node_num(t.Val)
    }
    Error(fmt.Sprintf("number expected, but got %s", t.Input))
    return nil
}

func expr() *Node {
    lhs := num()

    for {
        t := tokens[pos].(*Token)
        pos++
        if t.Ty == '+' || t.Ty == '-' {
            lhs = new_node(t.Ty, lhs, num())
        } else {
            break
        }
    }

    return lhs
}

func Parse(t []interface{}) *Node {
    tokens = t
    return expr()
}
