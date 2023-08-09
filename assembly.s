	.global _start
	.align 4
_start:
	mov x9, 2
	mov x9, x9
	add x9, x9, 3
	add x9, x9, 1
	mov x0, x9
	ret
