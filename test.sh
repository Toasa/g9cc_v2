#!/bin/bash

try() {
    expected="$1"
    input="$2"

    go run main.go "$input" > tmp.s
    gcc -o tmp tmp.s
    ./tmp
    actual="$?"

    if [ "$actual" == "$expected" ]; then
        echo "$input => $actual"
    else
        echo "$input: $expected expected, but got $actual"
        exit 1
    fi
}

try 0 0
try 46 46
try 28 ' 20 - 5 + 10 - 3 + 6'
try 36 '1+2+3+4+5+6+7+8'
try 20 '4 * 5'
try 26 '4 * 5 + 6'
try 34 '4 + 5 * 6'
try 62 '4 * 5 + 6 * 7'
try 6 '12/2'
try 12 '10/5 + 2*5'
try 3 '1 + 10 / 5'
try 3 '10 / 5 + 1'

echo OK
