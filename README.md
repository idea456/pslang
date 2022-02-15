# What is the PSU Language?

The PSU Language is an imperative programming language that takes inspiration from the pseudocode grammar to make writing code much more readable, easy and less complex.

The parser and interpreter are written in Golang, in which parsing is written using a top-down recursive descent approach.

## Variable declarations and usage

Variables are declared in the following syntax:
```
set x to 0;
set y to x + 1;
set word to "Hewwo";
set z to x + y;
```
Alternatively, PSUL also has dedicated increment/decrement keywords similar to `+=` and `-=`:
```
set x to 0;
increment x by 1;

set y to 1;
decrement y by x;
```

## Conditional statements

Conditional statements are written in the following manner:
```
set a to 0;
if a <= 0 then {
    increment a by 1;
} else {
    decrement a by 1;
}
```
