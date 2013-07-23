# Lune

A pure Go implementation of the Lua virtual machine (v5.2).

## Goals / Features

* Run Lua code. Make every pure Lua code (without C calls outside of the standard libs) work as expected. Every lua binary chunk compiled using `luac` should work on Lune, unless it relies on the C API.
* Implement the Lua standard libraries in Go. Make it transparent to Lua code.
* Embeddable. This is a Go package, it can be embedded in any Go application.
* Go-friendly. Just like Lua is the dynamic companion to C, Lune tries to bring this Batman and Robin *camaraderie* to Go. This means registering Go functions to be callable from Lua-on-Lune. This may be a port of the C API to Go, or something else that does more or less the same thing.

## Current status

Dormant. Unstable. Ugly. Unsafe. Unfast.

A few things work, though, like ummm... loading and deserializing the binary chunks. On 64-bit little-endian architectures at least. And running some trivial programs (see ./vm/testdata). Closures that actually use the closed-over environment currently don't work (*upvalues* in Lua literature). Tail calls and variadic arguments and return values don't work. Metamethods are not there yet.

## License

The [BSD 3-Clause license][bsd3], the same as the [Go language][golic]. [Lua itself is licensed under the MIT License][lua].

[lua]: http://www.lua.org/license.html
[golic]: http://golang.org/LICENSE
[bsd3]: http://opensource.org/licenses/BSD-3-Clause
