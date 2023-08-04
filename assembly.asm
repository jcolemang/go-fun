main:
	movq 2 %rax
	addq 3 %rax
	movq %rax %rcx
	movq %rcx %rcx
	movq 5 %rax
	addq %rcx %rax
	movq %rax %rcx
	movq %rcx %rcx
	movq %rcx %rax
	addq 10 %rax
	movq %rax %rcx
	movq %rcx %rax
	addq %rcx %rax
	movq %rax %rcx
	movq 1 %rax
	addq %rcx %rax
	movq %rax %rcx
