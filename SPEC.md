# Lune language specification

1. Introduction
2. Notation
3. Types
4. Blocks, declarations and scope
5. ...

## MVP language subset

Explicitly out of scope: 

* Interoperability with Go functions.
* Userdata (Go structures exposed to Lune).
* Meta-functions (operator overloading and prototypes).
* Coroutines (yield and resume).
* Variadic arguments, multiple return values.
* `for .. in` construct.

### Types

1. **nil** : The nil value (no value).
2. **bool** : Logical boolean type, which can be either `true` or `false`.
3. **number** : Any integer or floating-point number, represented internally as a `float64`.
4. **string** : A string of characters. Can be zero-length (empty string).
5. **function** : A function reference. Functions are closures, they close over their environment.
6. **table** : A map data structure, which holds a key and a value.

### Opcodes

MAYBE: Instead of LoadNil and LoadBool, use virtual Ks? For bools, to manage a = b > 5, use a CMP that stores result in a pseudo-register, and use a magic K index for LoadK? I.e. Nil = K(0), True = K(1), False = K(2), CmpResult = K(3) (or use higher values). Saves much instructions for the `var a = b > 5` case (a CMP - LT, LE, ... and a LoadK instead of CMP, JMP, 2 LOADBOOLs).

* MOVE    A B   R(A) := R(B) : Move registers around in the stack, i.e. to prepare arguments for a function call, or a for loop.
* LOADNIL A B   R(A) := ... := R(B) := nil : Initialize one or many registers to nil, since nil is not stored as a constant, something like `var a = nil` would be a LoadNil instruction.
* 