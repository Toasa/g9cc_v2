package x86

import (
    "fmt"
    . "g9cc_v2/common"
    . "g9cc_v2/reg"
    . "g9cc_v2/utils"
)

var label_num int
func gen_label() string {
    buf := fmt.Sprintf(".L%d", label_num)
    label_num++
    return buf
}

func Gen_x86(irv []interface{}) {

    var ret string = gen_label()
    fmt.Printf("    push rbp\n")
    fmt.Printf("    mov rbp, rsp\n")

    for i := 0; i < len(irv); i++ {
        ir := irv[i].(*IR)
        switch ir.Op {
        case IR_IMM:
            fmt.Printf("    mov %s, %d\n", Regs[ir.Lhs], ir.Rhs)
        case IR_MOV:
            fmt.Printf("    mov %s, %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case '+':
            fmt.Printf("    add %s, %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case '-':
            fmt.Printf("    sub %s, %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case '*':
            fmt.Printf("    mov rax, %s\n", Regs[ir.Rhs])
            fmt.Printf("    mul %s\n", Regs[ir.Lhs])
            fmt.Printf("    mov %s, rax\n", Regs[ir.Lhs])
        case '/':
            fmt.Printf("    mov rax, %s\n", Regs[ir.Lhs])
            fmt.Printf("    cqo\n")
            fmt.Printf("    div %s\n", Regs[ir.Rhs])
            fmt.Printf("    mov %s, rax\n", Regs[ir.Lhs])
        case IR_RETURN:
            fmt.Printf("    mov rax, %s\n", Regs[ir.Lhs])
            fmt.Printf("    jmp %s\n", ret)
        case IR_ALLOCA:
            if ir.Rhs != 0 {
                fmt.Printf("    sub rsp, %d\n", ir.Rhs)
            }
            fmt.Printf("    mov %s, rsp\n", Regs[ir.Lhs])
        case IR_LOAD:
            fmt.Printf("    mov %s, [%s]\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case IR_STORE:
            fmt.Printf("    mov [%s], %s\n", Regs[ir.Lhs], Regs[ir.Rhs])
        case IR_NOP:

        default:
            Error("invalid operator")
        }
    }

    fmt.Printf("%s:\n", ret)
    fmt.Printf("    mov rsp, rbp\n")
    fmt.Printf("    pop rbp\n")
    fmt.Printf("    ret\n")
}
