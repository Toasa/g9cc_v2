package common

// token.go

const (
    TK_NUM = iota + 256
    TK_EOF
)

type Token struct {
    Ty int // token type
    Val int // number literal
    Input string // token string
}

// parse.go

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

// ir.go

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
