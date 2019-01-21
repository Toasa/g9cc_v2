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

try 0 'return 0;'
try 46 'return 46;'
try 28 'return 20 - 5 + 10 - 3 + 6;'
try 36 'return 1+2+3+4+5+6+7+8;'
try 20 'return 4 * 5;'
try 26 'return 4 * 5 + 6;'
try 34 'return 4 + 5 * 6;'
try 62 'return 4 * 5 + 6 * 7;'
try 6 'return 12/2;'
try 12 'return 10/5 + 2*5;'
try 3 'return 1 + 10 / 5;'
try 3 'return 10 / 5 + 1;'
try 2 'a = 2; return a;'
try 10 'a=2; b=3+2; return a*b;'
try 10 'a=2; b=3+a; return a*b;'

echo OK
