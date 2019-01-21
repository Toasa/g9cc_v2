package ir

import (
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
    "fmt"
)

var code []interface{}

var regno int
var basereg int

var vars map[string]int
var bpoff int

func add(op int, lhs int, rhs int) *IR {
    ir := &IR{Op: op, Lhs: lhs, Rhs: rhs}
    code = append(code, ir)
    return ir
}

func gen_lval(node *Node) int {
    if node.Ty != ND_IDENT {
        Error("not an lvalue")
    }

    // varsに識別子の登録がされていない場合、bpのオフセットを8上げる
    if _, ok := vars[node.Name]; !ok {
        vars[node.Name] = bpoff
        bpoff += 8
    }

    r1 := regno
    regno++
    add(IR_MOV, r1, basereg)

    r2 := regno
    regno++
    off := vars[node.Name]
    add(IR_IMM, r2, off)
    add('+', r1, r2)
    add(IR_KILL, r2, -1)
    return r1
}

func gen_expr(node *Node) int {
    op := node.Ty
    switch op {
    case ND_NUM:
        r := regno
        regno++
        add(IR_IMM, r, node.Val)
        return r
    case ND_IDENT:
        r := gen_lval(node)
        add(IR_LOAD, r, r)
        return r
    case '=':
        lhs := gen_lval(node.Lhs)
        rhs := gen_expr(node.Rhs)
        add(IR_STORE, lhs, rhs)
        add(IR_KILL, rhs, -1)
        return lhs
    case '+', '-', '*', '/':
        lhs := gen_expr(node.Lhs)
        rhs := gen_expr(node.Rhs)

        add(op, lhs, rhs)
        add(IR_KILL, rhs, -1)

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
        add(IR_RETURN, r, -1)
        add(IR_KILL, r, -1)
        return
    case ND_EXPR_STMT:
        r := gen_expr(node.Expr)
        add(IR_KILL, r, -1)
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

    regno = 1
    basereg = 0
    vars = make(map[string]int)
    bpoff = 0

    var alloca *IR = add(IR_ALLOCA, basereg, -1)
    gen_stmt(node)
    alloca.Rhs = bpoff
    add(IR_KILL, basereg, -1)

    return code
}
