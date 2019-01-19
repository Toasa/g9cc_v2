package main

import (
    "fmt"
    "os"
    . "g9cc_v2/token"
    . "g9cc_v2/parse"
    . "g9cc_v2/ir"
    . "g9cc_v2/regalloc"
    . "g9cc_v2/x86"
    // "github.com/k0kubun/pp"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: g9cc <code>")
        return
    }

    input := os.Args[1]

    tokens := Tokenize(input + "\000")

    node := Parse(tokens)

    irv := Gen_ir(node)

    irv = Alloc_regs(irv)

    fmt.Printf("    .intel_syntax noprefix\n")
    fmt.Printf("    .globl _main\n")
    fmt.Printf("_main:\n")

    Gen_x86(irv)
}
