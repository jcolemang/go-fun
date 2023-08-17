
PROG ?= ./test-files/6.lang

program: c/runtime.o assembly.o
	ld -o program assembly.o c/runtime.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64
	chmod +x program

assembly.o: assembly.s
	as -arch arm64 -o assembly.o assembly.s

assembly.s: ./app/* ./pkg/**/*
	go run ./app/* $(PROG) assembly.s

c/runtime.o: c/runtime.c c/runtime.h
	clang -c -g -std=c99 c/runtime.c -o c/runtime.o

clean:
	rm -f program assembly.o runtime.o assembly.s
	rm -f ./test-files/*.s
	rm -f ./test-files/*.o
	rm -f ./c/*.o
	rm -f ./test-files/*.prog
