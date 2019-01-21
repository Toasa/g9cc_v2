package common

// token.go

const (
    TK_NUM = iota + 256
    TK_IDENT
    TK_RETURN
    TK_EOF
)

type Token struct {
    Ty int // token type
    Val int // number literal
    Name string // identifier
    Input string // token string
}

// parse.go

const (
    ND_NUM = iota + 256
    ND_IDENT
    ND_RETURN
    ND_COMP_STMT
    ND_EXPR_STMT
    ND_EOF
)

type Node struct {
    Ty int // node type
    Lhs *Node // left-hand side
    Rhs *Node // right-hand side
    Val int // number
    Name string // identifier
    Expr *Node
    Stmts []interface{} // Compound statement
}

// ir.go

const (
    IR_IMM = iota
    IR_MOV
    IR_RETURN
    IR_ALLOCA
    IR_LOAD
    IR_STORE
    IR_KILL
    IR_NOP
)

type IR struct {
    Op int
    Lhs int
    Rhs int
}
