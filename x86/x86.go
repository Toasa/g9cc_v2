package x86

import (
    "fmt"
    . "g9cc_v2/common"
    . "g9cc_v2/reg"
    . "g9cc_v2/utils"
)

func Gen_x86(irv []interface{}) {
    for i := 0; i < len(irv); i++ {
        ir := irv[i].(*IR)
        switch ir.Op {
        case IR_IMM:
            fmt.Printf("    mov %s, %d\n", Regs[ir.Lhs], ir.Rhs)
        case '+':
            fmt.Printf("    add %s, %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case '-':
            fmt.Printf("    sub %s, %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case '*':
            fmt.Printf("    mov rax, %s\n", Regs[ir.Rhs])
            fmt.Printf("    mul %s\n", Regs[ir.Lhs])
            fmt.Printf("    mov %s, rax\n", Regs[ir.Lhs])
        case IR_RETURN:
            fmt.Printf("    mov rax, %s\n", Regs[ir.Lhs])
            fmt.Printf("    ret\n")
        case IR_NOP:

        default:
            Error("invalid operator")
        }
    }
}
