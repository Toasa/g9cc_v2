package utils

import (
    "fmt"
    "os"
)

func Assert(b bool, msg string) {
    if !b {
        Error(msg)
    }
}

func Error(msgs ...string) {
    for _, msg := range msgs {
        fmt.Println(msg)
    }
    os.Exit(1)
}

func Isdigit(c uint8) bool {
    return '0' <= c && c <= '9'
}

func Isspace(c uint8) bool {
    return c == ' ' || c == '\n' || c == '\t'
}

func Isalpha(c uint8) bool {
    return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
