

HelloWorld: runtime.o HelloWorld.o
	ld -o HelloWorld HelloWorld.o runtime.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64 

HelloWorld.o: HelloWorld.s
	as -arch arm64 -o HelloWorld.o HelloWorld.s

runtime.o: runtime.c runtime.h
	clang -c -g -std=c99 runtime.c

clean:
	rm -f HelloWorld HelloWorld.o runtime.o
