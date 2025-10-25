# Syntax

## Identifier
### Rules
Make sure your identifiers are not violating these rules:
- Must not start with digit;
- Must start with letter of an underscore;
- No spaces or special characters;
- Cannot be a keyword name;
- Case-sensitive;
- Limit of the identifier is 255 symbols.

## Variables
To initialize a variable you need to follow this syntax:
```
<identifier> <type> <expression>
```

Example:
```
val i32 10
```

## Functions
In `dywoqlang`, everything is an expression, so initializing functions is:
```dl
<identifier> <return type> (<arg>, <variadic arg...>) {
	<body, instructions>
}
```

Example:
```
printstring void (str string) {
	println str;
	ret;
}
```

## Instructions
in `dywoqlang`, instructions are the commands that tell the interpreter what to do.
The syntax:
```dl
<identifier> <args>...;
```
Notice that you need to paste `;` at the end of instruction call. 

This way, you can call functions:
```dl
printstring void (str string) {
	println str;
	ret;
}

main i32 () {
	printstring "Hi!"; 
	ret 0;
}
```

## Modifiers-conversions
[Modifiers-conversions](./modifiable-conversions.md)' arguments are always surrounded by the parens:
```
<identifier>(<args>)
```

The example:
```dl
val i32 const(10)
```
In this code, `const(10)` is a modifier-conversion.
