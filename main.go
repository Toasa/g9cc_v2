package main

import (
    "fmt"
    "os"
)

const (
    TK_NUM = iota + 256
    TK_EOF
)

type Token struct {
    Ty int // token type
    Val int // number literal
    Input string // token string
}

func add_token(ty int, input string) *Token {
    t := &Token{Ty: ty, Input: input}
    tokens = append(tokens, t)
    return t
}

func tokenize(s string) {
    //cur_token := 0
    i_input := 0

    for s[i_input] != '\000' {

        // while space
        if Isspace(s[i_input]) {
            i_input++
            continue
        }

        // + or -
        if s[i_input] == '+' || s[i_input] == '-' {
            add_token(int(s[i_input]), string(s[i_input]))
            i_input++
            continue
        }

        // number
        if Isdigit(s[i_input]) {
            var num int = int(s[i_input] - '0')
            i_input++
            for ; Isdigit(s[i_input]); i_input++ {
                num = num * 10 + int(s[i_input] - '0')
            }

            t := add_token(TK_NUM, string(num))
            t.Val = num
            continue
        }

        fmt.Printf("cannot tokenize: %c\n", s[i_input]);
        os.Exit(1)
    }

    add_token(TK_EOF, "")

}

// Recursive-descent parser

var pos int = 0
var tokens []interface{}

const (
    ND_NUM = iota + 256
    ND_EOF
)

type Node struct {
    Ty int // node type
    Lhs *Node // left-hand side
    Rhs *Node // right-hand side
    Val int // number
}

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

func error(msgs ...string) {
    for _, msg := range msgs {
        fmt.Println(msg)
    }
    os.Exit(1)
}

func num() *Node {
    t := tokens[pos].(*Token)
    if t.Ty == TK_NUM {
        pos++
        return new_node_num(t.Val)
    }
    error(fmt.Sprintf("number expected, but got %s", t.Input))
    return nil
}

func expr() *Node {
    var lhs *Node = num()
    for {
        t := tokens[pos].(*Token)
        op := t.Ty
        if op == '+' || op == '-' {
            pos++
            lhs = new_node(op, lhs, num())
        } else {
            break
        }
    }
    t := tokens[pos].(*Token)
    if t.Ty != TK_EOF {
        error("stray token")
    }
    return lhs
}

// Intermediate representation
const (
    IR_IMM = iota
    IR_MOV
    IR_RETURN
    IR_KILL
    IR_NOP
)

type IR struct {
    Op int
    Lhs int
    Rhs int
}

func new_ir(op int, lhs int, rhs int) *IR {
    return &IR{Op: op, Lhs: lhs, Rhs: rhs}
}

// var ins [1000]*IR
// //index of ins
// var inp int
var irv []interface{}

var regno int

func assert(b bool, msg string) {
    if !b {
        error(msg)
    }
}

func gen_ir_sub(node *Node) int {
    op := node.Ty
    if op == ND_NUM {
        r := regno
        regno++
        irv = append(irv, new_ir(IR_IMM, r, node.Val))
        return r
    } else if op == '+' || op == '-' {

        lhs := gen_ir_sub(node.Lhs)
        rhs := gen_ir_sub(node.Rhs)

        irv = append(irv, new_ir(op, lhs, rhs))
        irv = append(irv, new_ir(IR_KILL, rhs, 0))
        return lhs
    } else {
        error("invalid operator")
        return 0
    }
}

func gen_ir(node *Node) {
    r := gen_ir_sub(node)
    irv = append(irv, new_ir(IR_RETURN, r, 0))
}

var regs []string = []string{"rdi", "rsi", "r10", "r11", "r12", "r13", "r14", "r15"}
var used [8]bool

// gen_ir()の段階で、中間表現のlhs, rhsには暫定的にレジスタ(のインデックス, regno)が割り当てられている
// regnoはインクリメンタルに増えていくので、reallocateする必要がある
// reg_mapはそのregnoをインデックスに、8つのレジスタの中から実際に割り当てるレジスタのインデックスを
// 値として格納する
var reg_map[1000]int

// ir_regは暫定的に決められたレジスタのインデックス
func alloc(ir_reg int) int {
    if reg_map[ir_reg] != -1 {
        return reg_map[ir_reg]
    }

    for i := 0; i < len(regs); i++ {
        if !used[i] {
            used[i] = true
            reg_map[ir_reg] = i
            return i
        }
    }

    error("exhausted register")
    return 0
}

func kill(r int) {
    assert(used[r], "register already used")
    used[r] = false
}

func alloc_regs() {
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
            error("invalid operator")
        }
    }
}

func gen_x86() {
    for i := 0; i < len(irv); i++ {
        ir := irv[i].(*IR)
        switch ir.Op {
        case IR_IMM:
            fmt.Printf("    mov %s, %d\n", regs[ir.Lhs], ir.Rhs)
        case '+':
            fmt.Printf("    add %s, %s\n", regs[ir.Lhs], regs[ir.Rhs])
        case '-':
            fmt.Printf("    sub %s, %s\n", regs[ir.Lhs], regs[ir.Rhs])
        case IR_RETURN:
            fmt.Printf("    mov rax, %s\n", regs[ir.Lhs])
            fmt.Printf("    ret\n")
        case IR_NOP:

        default:
            error("invalid operator")
        }
    }
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: g9cc <code>")
        return
    }

    for i := 0; i < len(reg_map); i++ {
        reg_map[i] = -1
    }

    input := os.Args[1] + "\000"

    tokenize(input)

    var node *Node = expr()

    gen_ir(node)

    alloc_regs()

    fmt.Printf("    .intel_syntax noprefix\n")
    fmt.Printf("    .globl _main\n")
    fmt.Printf("_main:\n")

    gen_x86()
}

func Isdigit(c uint8) bool {
    return '0' <= c && c <= '9'
}

func Isspace(c uint8) bool {
    return c == ' ' || c == '\n' || c == '\t'
}
