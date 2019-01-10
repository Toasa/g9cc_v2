package regalloc

import (
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
    . "g9cc_v2/reg"
)

var used [8]bool

// gen_ir()の段階で、中間表現のlhs, rhsには暫定的にレジスタ(のインデックス, regno)が割り当てられている
// regnoはインクリメンタルに増えていくので、reallocateする必要がある
// reg_mapはそのregnoをインデックスに、8つのレジスタの中から実際に割り当てるレジスタのインデックスを
// 値として格納する
var reg_map[1000]int

var irv []interface{}

// ir_regは暫定的に決められたレジスタのインデックス
func alloc(ir_reg int) int {
    if reg_map[ir_reg] != -1 {
        return reg_map[ir_reg]
    }

    for i := 0; i < len(Regs); i++ {
        if !used[i] {
            used[i] = true
            reg_map[ir_reg] = i
            return i
        }
    }

    Error("exhausted register")
    return 0
}

func kill(r int) {
    Assert(used[r], "register already used")
    used[r] = false
}

func Alloc_regs(irv []interface{}) []interface{} {
    irv = irv

    for i := 0; i < len(reg_map); i++ {
        reg_map[i] = -1
    }

    for i := 0; i < len(irv); i++ {
        ir := irv[i].(*IR)
        switch ir.Op {
        case IR_IMM:
            ir.Lhs = alloc(ir.Lhs)
        case '+', '-':
            ir.Lhs = alloc(ir.Lhs)
            ir.Rhs = alloc(ir.Rhs)
        case IR_KILL:
            kill(reg_map[ir.Lhs])
            ir.Op = IR_NOP
        case IR_RETURN:
            kill(reg_map[ir.Lhs])
        default:
            Error("invalid operator")
        }
    }

    return irv
}
