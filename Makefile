CGO=../cgo-gccgo/cgo

TARGET=alpm

CGOFILES=alpm.go\
	 handle.go\
	 error.go\
	 db.go\
	 package.go\
	 types.go

GOFILES=defs.go\
	enums.go

EXTRA_GOFILES=$(patsubst %.go,_obj/%.cgo1.go,$(CGOFILES)) _obj/_cgo_gotypes.go

$(TARGET).o: $(CGOFILES) $(GOFILES)
	$(CGO) -- $(CGOFILES)
	gccgo -o $(TARGET)1.o -c $(GOFILES) $(EXTRA_GOFILES)
	gcc -o alpm2.o -c _obj/_cgo_defun.c
	ar cru libalpm.a alpm1.o alpm2.o

test:
	gccgo -o alpm_test alpm_test.go $(TARGET).a -lalpm
	./alpm_test

clean:
	rm -rf *.o *.a *.so _obj _test _testmain.go *.exe _cgo* test.out build.out *.gox

