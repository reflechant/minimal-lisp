# miniLISP

A barebone LISP interpeter described in Paul Graham's paper ["The
Roots of LISP"](https://paulgraham.com/rootsoflisp.html) implemented
in Go.

It's not intended to be practically usable, only to demonstrate the
"Maxwell equations" nature of LISP by implementing the Eval function
in itself using 7 axiomatic operators.

The idea is that in `core.lisp` we define an `eval.` function that is
an *interpreter for this language written in itself*.
It's considered a "rite of passage" moment for compiled languages to write
the compiler in itself. For interpreted languages it's similar.

The whole point of the paper was to show the beauty and mathematical
profundity of LISP which achieves this goal by using just 7 primitive operators!

LISP is like the coordinate system of programming.
If you have X, Y and Z you don't need any more axes for a 3-dimensional space.
ALl other languages choose to contain lots of non-orthogonal primitives
and be bloated with syntactic sugar. Once you really get what Paul Graham
was trying to show (after McCarthy) you will not be able to enjoy this tower
of unnecessary complexity our industry so much likes to be building.

## What can I do with it?

Just run main.go and read the article and the code to understand.

You can try to evaluate these expressions in the REPL:

``` common-lisp
    (print 'x)
```

``` common-lisp
    (quote a)
```

``` common-lisp
    'a
```

``` common-lisp
    (atom 'x)
```

``` common-lisp
    (eq 'a 'a)
```

Example from p.9 of "The Roots of LISP":
``` common-lisp
    (eval. 'x '((x a) (y b)) )
```
