program: runtime.o assembly.o
	ld -o program assembly.o runtime.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64

assembly.o: assembly.s
	as -arch arm64 -o assembly.o assembly.s

assembly.s: ./app/*
	go run ./app/* $(PROG)

runtime.o: runtime.c runtime.h
	clang -c -g -std=c99 runtime.c

clean:
	rm -f program assembly.o runtime.o assembly.s
