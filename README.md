# ikou

* `ikou` is a Lisp interpreter written in Go.

## Tokens

* Identifier `[a-zA-Z+\-*/=][a-zA-Z0-9+\-*/=]*`
  * A name for a function, constant, let binding, etc.
  * Keywords: `fn`, `lambda`, `if`, `let`, `define`
  
* Open bracket `(`
* Close bracket `)`
* Colon `:`
* Integer `[0-9]+`
* Float `[0-9]+\.[0-9]+`

## Language

### Comments

```
; Comments are defined like this

; A comment begins with a semi-colon and lasts until the end of the line
```

### If-Then-Else Expressions

* Syntax: `(if <bool expr> <expr of type T> <expr of type T>)`

```
(# if a number is even, return it, else return that number add 1 #)
(if (= (% n 2) 0) n (+ n 1))
```

### Let Bindings

* Syntax: `(let ((<binding> <expr>) (<binding> <expr>) ...) <exprs>)`
* The final expression in `<exprs>` determines the value returned.
* `<binding>` refers to either an identifier or can be used to separate the head and tail of a list when in the form `<ident>:<ident>`.

```
(let ((my_int 25) (head:tail my_list))
    (println my_int)
    (println head))
```

### Named Functions

* Syntax: `(fn <ident> <args> <exprs>)`
* `<args>` can be single identifier `<ident>` or a list of identifiers `(<ident> <ident> ...)`.
* The final expression in the function body determines the value returned.
* Functions can be partially applied.

```
(fn add_ints (x y) (+ x y))

(fn factorial (n)
    (if (< n 2)
        1
        (* n (factorial (- n 1)))))
```

### Anonymous Functions

* Syntax: `(lambda <args> <exprs>)`
