package ir

import (
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
)

var irv []interface{}

var regno int

func new_ir(op int, lhs int, rhs int) *IR {
    return &IR{Op: op, Lhs: lhs, Rhs: rhs}
}

func gen(node *Node) int {
    op := node.Ty
    switch op {
    case ND_NUM:
        r := regno
        regno++
        irv = append(irv, new_ir(IR_IMM, r, node.Val))
        return r
    case '+', '-', '*', '/':
        lhs := gen(node.Lhs)
        rhs := gen(node.Rhs)

        irv = append(irv, new_ir(op, lhs, rhs))
        irv = append(irv, new_ir(IR_KILL, rhs, 0))
        return lhs
    default:
        Error("invalid operator")
        return 0
    }
}

func Gen_ir(node *Node) []interface{} {
    r := gen(node)
    irv = append(irv, new_ir(IR_RETURN, r, 0))
    return irv
}
