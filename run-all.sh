set -xe

for f in `find ./test-files -name '*.lang'`
do
	go run ./app/* $f $f.s
    chmod +x $f.s
	clang -c -g -std=c99 c/runtime.c
	as -arch arm64 -o $f.o $f.s
	ld -o $f.prog $f.o runtime.o -lSystem -syslibroot `xcrun -sdk macosx --show-sdk-path` -e _start -arch arm64
	chmod +x $f.prog
done

set +xe
