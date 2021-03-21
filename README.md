# Vaeryn Virtual Computer

A hobby project to model a general-purpose computer architecture.

Requirements: `make`, `docker`, `docker-compose`.

## Try it

A packaged docker image contains some tools to try out VVC.

```
make vvc-run
```

This image contains 2 executables:
* `compile` which will compile a `.vvb` program into an executable binary
* `execute` will spin up a VVC, load and execute a specified binary file.

The `/programs` directory contains some sample `.vvb` files. The following
compiles and executes `loops.vvb`, which simply loops up to 10.

```
$ compile /programs/loops.vvb > /tmp/myprogram
Compiling file /programs/loops.vvb
Done
$ execute /tmp/myprogram 
2021/03/21 22:55:42 CPU: Register 3: 0x0 (0)
2021/03/21 22:55:42 CPU: Register 3: 0x1 (1)
2021/03/21 22:55:42 CPU: Register 3: 0x2 (2)
2021/03/21 22:55:42 CPU: Register 3: 0x3 (3)
2021/03/21 22:55:42 CPU: Register 3: 0x4 (4)
2021/03/21 22:55:42 CPU: Register 3: 0x5 (5)
2021/03/21 22:55:42 CPU: Register 3: 0x6 (6)
2021/03/21 22:55:42 CPU: Register 3: 0x7 (7)
2021/03/21 22:55:42 CPU: Register 3: 0x8 (8)
2021/03/21 22:55:42 CPU: Register 3: 0x9 (9)
```

## Architecture

Some rough notes on the architecture:
* Components:
  * `machine` provides an single interface for creating, booting and operating the computer.
  * `cpu` which aligns with the CPU in a modern architecture. Executes instructions serially
    loaded from memory.
  * `memory` is a simple I/O interface to randomly access storage.
* Fixed-length words. A word is 4 bytes, `int32`.
* Words are the addressable unit in memory.
* The CPU contains a small number of registers to facilitate computation. Each register stores
  a single word.
* Only unsigned integers are supported (for now).

## Project Layout

Modelled on https://github.com/golang-standards/project-layout.

## .vvb files

Vaeryn Virtual Bytecode files (`*.vvb`) are VVC program source code. A compiler
is included in which will compile `.vvb` files to binary files in a format which
can be executed on VVC.

See [./programs](./programs) for some examples of `.vvb` files.

The `sample` make target demonstrates compilation of a `.vvb` file followed by
its execution.

A bytecode file is a series of instruction, each written on a new line. All extra
whitespace is ignored.

```
[label:] INSTRUCTION [argA] [argB] # Comments follow
```

* `[label:]` is an optional label that can be used later in jump instruction arguments.
  Labels can only contain letters (no digits or whitespace).
* `INSTRUCTION` is the 3-letter instruction code. See 
  [instruction_set.go](./internal/core/instruction_set.go) for a list of available instructions.
* `[argA]` and `[argB]` are arguments. These can take the form:
    * `rX` (e.g. `r0`) where X is the CPU register being addressed.
    * `labelRef` (e.g. `loopStart`) which the compiler will replace down to the address of
      the instruction with this label.
    * A number, in decimal.
* `# Comment...` is a comment. Anything after a `#` on a line will be ignored by the compiler.

Some checks performed by the compiler are:
* Instructions have the correct number of arguments
* Instruction arguments are of the expected type, e.g. an instruction expecting a register argument
  will reject unless its argument is of the form `rX`.
* All label references are valid.
* The resulting compiled file is >0 bytes in length.