package token

import (
    "fmt"
    "os"
    "strings"
    . "g9cc_v2/common"
    . "g9cc_v2/utils"
)

var tokens []interface{}
var keywords map[string]int

func add_token(ty int, input string) *Token {
    t := &Token{Ty: ty, Input: input}
    tokens = append(tokens, t)
    return t
}

func scan(s string) []interface{} {
    //cur_token := 0
    i_input := 0

    for s[i_input] != '\000' {

        // while space
        if Isspace(s[i_input]) {
            i_input++
            continue
        }

        // single-letter token
        if strings.Contains("+-*/;=", string(s[i_input])) {
            add_token(int(s[i_input]), string(s[i_input]))
            i_input++
            continue
        }

        // Identifier
        if Isalpha(s[i_input]) || s[i_input] == '_' {
            len := 1
            for i := len + i_input; Isalpha(s[i]) || Isdigit(s[i]) || s[i] == '_'; {
                len++
                i = len + i_input
            }

            name := s[i_input : len + i_input]

            ty, ok := keywords[name]
            if !ok {
                ty = TK_IDENT
            }

            t := add_token(ty, name)
            t.Name = name
            i_input += len
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

    return tokens
}

func Tokenize(s string) []interface{} {
    keywords = make(map[string]int)
    keywords["return"] = TK_RETURN
    return scan(s)
}
