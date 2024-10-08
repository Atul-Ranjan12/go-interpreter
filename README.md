# Lang Interpreter

Lang is a simple, dynamically-typed programming language with struct-based features. This interpreter is implemented in Go and provides a straightforward way to execute Lang programs.
Lang interpreter is implemented in go

## Language Features

- Dynamic typing
- Struct-based programming with methods
- Functions and closures
- Control structures (if-else, while, for)
- Basic arithmetic and logical operations

## Grammar

Lang's grammar is designed to be intuitive and easy to read. Here's a simplified version of the grammar:

```
program -> declaration* EOF

## Statements

statement -> ifStatement
           | exprStatement
           | printStatement
           | forStatement
           | whileStatement
           | returnStatement
           | block

returnStatement -> "return" expression ";"

forStatement -> "for" "(" (varDeclaration | expression) ";" expression? ";" expression? ")" statement

whileStatement -> "while" "(" expression ")" statement

ifStatement -> "if" "(" expression ")" statement ("else" statement)?

block -> "{" declaration* "}"

## Declarations

declaration -> funcDeclaration
             | classDeclaration
             | varDeclaration
             | statement
             | breakStatement

classDeclaration -> "class" IDENTIFIER "{" function* "}"

funcDeclaration -> "fun" function

function -> IDENTIFIER "(" parameters? ")" block

parameters -> IDENTIFIER ("," IDENTIFIER)*

breakStatement -> "break" ";"

varDeclaration -> "var" IDENTIFIER ("=" expression)? ";"

exprStatement -> expression ";"

printStatement -> "print" "(" expression ")" ";"

## Expressions

expression -> assignment

assignment -> (call ".")? IDENTIFIER "=" assignment
            | logic_or

logic_or -> logic_and ("or" logic_and)*

logic_and -> equality ("and" equality)*

equality -> comparison (("!=" | "==") comparison)*

comparison -> term ((">" | ">=" | "<" | "<=") term)*

term -> factor (("/" | "*") factor)*

factor -> unary (("+" | "-") unary)*

unary -> ("!" | "-") unary
       | call

call -> primary (("(" arguments? ")") | "." IDENTIFIER)*

arguments -> expression ("," expression)*

primary -> NUMBER
         | STRING
         | "true"
         | "false"
         | "nil"
         | "(" expression ")"
         | IDENTIFIER
```

## Example Syntax

Here's an example of Lang code demonstrating various language features:

```lang
struct Animal {
    construct() {
        this.type = "Lion";
        println "Constructed";
        return true;
    }

    makeSound() {
        println this.type;
    }
}

def add(a, b) {
    return a + b;
}

def fib(n) {
    if (n <= 1) {
        return 1;
    }

    return fib(n - 1) + fib(n - 2);
}

var a = Animal();
println "Animal Type " + a.type;
a.makeSound();

println "Addition";
var x = 1;
var y = 0;
println add(x, y);

println "Fib";
println fib(22);

for (var i = 0; i < 10; i = i + 1){
    println i;
}

```

This example demonstrates:
- Struct definition with methods
- Function definitions
- Variable declarations
- Conditional statements
- Recursion
- Struct instantiation and method calls
- Basic arithmetic operations

## Running Lang Programs

To run a Lang program, use the interpreter as follows:

```
go run cmd/main.go path/to/your/program.lang
```

## Implementation Details

The interpreter is implemented in Go and consists of several key components:

1. Lexer: Tokenizes the input source code
2. Parser: Builds an Abstract Syntax Tree (AST) from tokens
3. Interpreter: Executes the AST

The implementation follows a tree-walk interpreter pattern, which is suitable for learning and prototyping but may not be as performant for large-scale applications.

## Key Language Characteristics

- Structs: Lang uses structs as its primary mechanism for creating custom data types with associated methods.
- No Inheritance: Unlike traditional object-oriented languages, Lang does not support inheritance.
- Dynamic Typing: Variables in Lang are dynamically typed.
- First-Class Functions: Functions in Lang are first-class citizens and can be passed as arguments or returned from other functions.

## Future Enhancements

- Standard library with common functions
- Error handling and better error reporting
- REPL (Read-Eval-Print Loop) for interactive use
- Performance optimizations

## Contributing

Contributions to Lang are welcome! Please feel free to submit pull requests, create issues, or suggest new features.
