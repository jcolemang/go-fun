//
// Assembler program to print "Hello World!"
// to stdout.
//
// X0-X2 - parameters to Unix system calls
// X16 - Mach System Call function number
//

.global _start			// Provide program starting address to linker
.align 2			// Make sure everything is aligned properly

// Setup the parameters to print hello world
// and then call the Kernel to do it.
_start:
	stp	x29, LR, [sp, #-16]!

    bl _read_int
    bl _add_five
    bl _print_int

	ldp	x29, LR, [sp], #16
    mov X0, 0
    ret
