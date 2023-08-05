.global _main	            // Provide program starting address
.align 4

_main:	
	// Setup
    // stp = store pair
    // stores x29, LR in sp, sp-8, and the ! sets the stack pointer to be sp-16
    // x29 is the frame pointer
	stp	x29, LR, [sp, #-16]!     ; Save LR, FR
    // I don't know and don't think I care
	adrp  	    X0, ptfStr@PAGE // printf format str
	add	X0, X0, ptfStr@PAGEOFF


	mov	x2, #4711
	mov	x3, #3845
	mov     x10, #65
    // stores x10, x2, and x3 in the referenced locations.
    // first argument updates sp to be sp-32
    // I _think_ that it doesn't pick -24 because the program is
    // four byte alligned
    // honestly I don't know why it stores these
	str	x10, [SP, #-32]!
	str	x2, [SP, #8]
	str	x3, [SP, #16]

    // bl is not branch label but branch link. This expects to return
	bl	    _printf	// call printf

	add	    SP, SP, #32	// Clean up stack

	MOV	X0, #0		// return code
	ldp	x29, LR, [sp], #16     ; Restore FR, LR
	RET
.data
ptfStr: .asciz	"Hello World %c %ld %ld\n"
.align 4
.text
