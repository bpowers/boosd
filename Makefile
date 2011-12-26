include $(GOROOT)/src/Make.inc

all: install

DIRS=\
	token\
	parser\
	driver\

NOTEST=
NOBENCH=

TEST=\
	$(filter-out $(NOTEST),$(DIRS))

BENCH=\
	$(filter-out $(NOBENCH),$(TEST))

clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
gofmt.dirs: $(addsuffix .gofmt, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
testshort.dirs: $(addsuffix .testshort, $(TEST))
bench.dirs: $(addsuffix .bench, $(BENCH))

%.clean:
	+$(MAKE) -C $* clean

%.install:
	+@echo install $*
	+@$(MAKE) -C $* install.clean >$*/build.out 2>&1 || (echo INSTALL FAIL $*; cat $*/build.out; exit 1)

%.gofmt:
	+@echo gofmt $*
	+gofmt -w $*/*.go

%.nuke:
	+$(MAKE) -C $* nuke

%.test:
	+@echo test $*
	+@$(MAKE) -C $* test.clean >$*/test.out 2>&1 || (echo TEST FAIL $*; cat $*/test.out; exit 1)

%.testshort:
	+@echo test $*
	+@$(MAKE) -C $* testshort.clean >$*/test.out 2>&1 || (echo TEST FAIL $*; cat $*/test.out; exit 1)

%.bench:
	+$(MAKE) -C $* bench	

clean: clean.dirs

install: install.dirs

gofmt: gofmt.dirs

test:	test.dirs

testshort: testshort.dirs

nuke: nuke.dirs
	rm -rf "$(GOROOT)"/pkg/*

deps:
	./deps.bash

echo-dirs:
	@echo $(DIRS)
