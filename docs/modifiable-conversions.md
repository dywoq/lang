# Modifiers-conversions

In `dywoqlang`, everything is an expression. So, you can't use like:
```dl
var const i64 10
```
Or
```dl
var consteval i64 10
```

Instead of this, the language provides **Modifiers-conversions**. 

## Description
**Modifiers-conversions** are something you can compare to type conversions, but these ones don't convert, for example, `i32` into `i64`. They modify how the symbol behaves.

## All available conversions

### `const`

#### General
`const` makes a constant, meaning its value can't be changed.
The only difference from `consteval` is `const` can make the value constant-evaluated, but only if it doesn't violate limitations (see [`consteval` - Functions](#functions-consteval-functions)).

#### Important regarding functions
This is recommended to use for functions as every symbol's value can be replaced by other developers:
```dl
# Module Foo:
# 
# Abc i64 (x i64) {
#	does something with x..	
# }
main i32 () {
	# foo.Abc is not a constant!
	mov foo.Abc, (x i64) {};
}
```
In the module `Foo`, we create a function `Abc` (not a constant).
And in `main` function, we replace the Abc function body to `(x i64)`: This is the problem, the language will use the assigned value in `main` instead of the one you put initially in the module `Foo`.

So use:
```dl
# Module Abc:
# 
# Abc i64 const((x i64) {
#	does something with x..	
# })
main i32 () {
	# foo.Abc is a constant, so an error!
	mov foo.Abc, (x i64) {};
}
```
Now, if you use `const(...)`, the functions can't be replaced.
You may not use `const(...)`, but if you're developing a library or safe-code - this is a mus-use. 

### `consteval`

#### General
`consteval` makes a variable constant-evaluated:
```dl
val i32 consteval(10)
```

#### Limitations
You can't use `const(...)` conversion, as the `consteval(10)` is already a constant and it's too redundant:
```dl
val i32 consteval(const(10)) # Not allowed
```

#### Bytecode
`consteval(...)` changes how the bytecode generator converts AST tree into the bytecode.
If you don't use `consteval`, the bytecode would be:

Original code:
```dl
val i32 20+20
```
Bytecode:

```
GLOBAL_VAR val i32 0
GLOBAL_VAR_ADD val 20
GLOBAL_VAR_ADD val 20
```

If you do, the bytecode would be:
```dl
val i32 consteval(20+20)
```

Bytecode:
```
GLOBAL_VAR val i32 40
```

In `consteval(...)` case, the bytecode automatically evaluates `20+20`.
You should always use this when you're

#### Functions

Functions are allowed to be used in `consteval(...)`, but with limitations.

First, you can't use system calls, as they're runtime only:
```dl
# An error!
write void consteval((msg str) {
	# system call...
})
```

Second, `consteval` functions must always return something to be evaluated and stored:
```dl
# An error! foo return type is void
foo void consteval(() {
	# do something...
})
```

**How does consteval functions present in the bytecode?**

They use the same instructions (like the early-mentioned `GLOBAL_VAR`), but the bytecode evaluates the functions, already inserting the ready results:

Original code:
```dl
do_hard i32 consteval(() {
	return 3 * 1000
})
```

Bytecode:
```dl
GLOBAL_VAR do_hard i32 3000
```

### `copy`

`copy` is a modifier-conversion that copies the value from another symbol.
It can be function, constant or constant-evaluated.

#### Limitations

You can't apply nested modifier-conversions. For example, this code will cause an error:

```dl
x i32 20
val i32 copy(consteval(x)) # Invalid
```

Correct:
```
x i32 20
val i32 copy(x)
```

#### Bytecode

The bytecode presentation is just a copy of another variable:
```
GLOBAL_VAR x i32 20
GLOBAL_VAR val i32 20
```
