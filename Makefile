g9cc: main.go

test: g9cc
	./test.sh

clean:
	rm -f tmp*

gdb: tmp.s
	gcc -c -nostdlib tmp.s
	gdb tmp.o
