package ir

import (
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
    "fmt"
)

var code []interface{}

var regno int

func add(op int, lhs int, rhs int) *IR {
    ir := &IR{Op: op, Lhs: lhs, Rhs: rhs}
    code = append(code, ir)
    return ir
}

func gen_expr(node *Node) int {
    op := node.Ty
    switch op {
    case ND_NUM:
        r := regno
        regno++
        add(IR_IMM, r, node.Val)
        return r
    case '+', '-', '*', '/':
        lhs := gen_expr(node.Lhs)
        rhs := gen_expr(node.Rhs)

        add(op, lhs, rhs)
        add(IR_KILL, rhs, 0)

        return lhs
    default:
        Error("invalid operator")
        return 0
    }
}

func gen_stmt(node *Node) {
    switch node.Ty {
    case ND_RETURN:
        r := gen_expr(node.Expr)
        add(IR_RETURN, r, 0)
        add(IR_KILL, r, 0)
        return
    case ND_EXPR_STMT:
        r := gen_expr(node.Expr)
        add(IR_KILL, r, 0)
        return
    case ND_COMP_STMT:
        for i := 0; i < len(node.Stmts); i++ {
            n := node.Stmts[i].(*Node)
            gen_stmt(n)
        }
        return
    }

    Error(fmt.Sprintf("unknown node: %d", node.Ty))

}

func Gen_ir(node *Node) []interface{} {
    Assert(node.Ty == ND_COMP_STMT, "type of root node is not ND_COMP_STMT")
    gen_stmt(node)
    return code
}
