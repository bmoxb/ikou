# ikou

*ikou* is a Lisp language with a bytecode compiler and virtual machine written in Go.

## Roadmap

* [ ] Lexer
  * [#] Define all necessary token types
  * [#] Handle single character tokens (open close brackets, quotes, etc.)
  * [#] Handle number tokens (integers and floats)
  * [#] Handle character tokens
  * [#] Handle identifier tokens
  * [#] Handle string tokens
  * [#] Handle comments
  * [#] Appropriately track positioning of tokens (line number and horizontal position)
  * [#] Produce useful error messages
  * [#] Thorough testing
* [ ] Parser
* [ ] Bytecode generation
* [ ] Virtual machine (bytecode execution)

## Language

### Tokens

* Identifier `[a-zA-Z+\-*/=][a-zA-Z0-9+\-*/=]*`
  * A name for a function, constant, let binding, etc.
  * Keywords: `lambda`, `if`, `let`, `define`
* Open bracket `(`
* Close bracket `)`
* Open square bracket `[`
* Close square bracket `]`
* Colon `:`
* Quote `'`
* Backquote `,`
* Integer `~?[0-9]+`
* Float `~?[0-9]+\.[0-9]+`
* String `"[^"]*"`
  * Escape sequences `\"`, `\n`, `\t` can be used inside string literals.
* Character `(\\.|\\space|\\newline|\\tab)`

### Literals

* Integer: `0`, `12`, `5746782678`, `~12`
  * Note that negative numbers are indicated with `~` and not `-`. This is to prevent confusion with identifiers which can contain both `-` and numeral characters.
* Float: `5.12`, `123.456`, `.5`, `1.`, `~12.2`
* List: `[<expr> <expr> ...]`
  * Equivalent to `(list <expr> <expr> ...)` in most other Lisp languages.
* Character: `\a`, `\!`, `\夢`, `\space`, `\newline`, `\tab`
* String: `"abc"`, `""`, `"Hello, 世界！"`, `"line 1\nline 2"`
  * String literals are just syntactic sugar for lists of characters (i.e., `"abc"` is equivalent to `[\a \b \c]`).
* Boolean; `true`, `false`

### Comments

```
; Comments are defined like this

; A comment begins with a semi-colon and lasts until the end of the line
```

### If-Then-Else Expressions

* Expression that yields some value when a condition is true, and a different value otherwise.
* Syntax: `(if <bool expr> <expr of type T> <expr of type T>)`

```
(# if a number is even, return it, else return that number add 1 #)
(if (= (% n 2) 0) n (+ n 1))
```

### Let Bindings

* Create some variables within a scoped block.
* Syntax: `(let ((<binding> <expr>) (<binding> <expr>) ...) <exprs>)`
* The final expression in `<exprs>` determines the value returned.
* `<binding>` refers to either an identifier or can be used to separate the head and tail of a list when in the form `<ident>:<ident>`.

```
(let ((my_int 25) (head:tail my_list))
    (println my_int)
    (println head))
```

### Define

* Create a global constant.
* Syntax: `(define <ident> <expr>)`

```
(define my_favourite_number 35)
(define my_function (lambda (x) (* x 2)))
```

### Anonymous Functions

* Syntax: `(lambda <args> <exprs>)`
* `<args>` can be single identifier or a list of identifiers.
* The final expression in the function body determines the value returned.
* Functions can be partially applied.

### Named Functions

* `fn` is a macro that equivalent to using `define` to give a name to a `lambda` function.
* Syntax: `(fn <ident> <args> <exprs>)`

```
(fn add_ints (x y) (+ x y))

(fn factorial (n)
    (if (< n 2)
        1
        (* n (factorial (- n 1)))))
```

### Quote

* Specify that an expression should not be evaluated and instead be treated as a literal list.
* Syntax: `'(<expr> <expr> ...)`

### Backquote

* Selectively evaluate some of the expressions inside a quoted list.
* Syntax: `'(,<expr> ,<expr> ...)`
