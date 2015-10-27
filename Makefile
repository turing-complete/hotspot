root := $(shell pwd)
source := $(root)/source
target := $(root)/target
object := main.syso

export OUTPUT_DIR := $(target)

all: $(object)

install: $(object)
	go install

$(object): $(target)/libcircuit.a
	mkdir -p $(target)/$@
	cd $(target)/$@ && ar x ../libcircuit.a
	ld -r -o $@ $(target)/$@/*.o

$(target)/libcircuit.a: $(target)
	mkdir -p $(target)
	$(MAKE) -C $(source)

clean:
	$(MAKE) -C $(source) clean
	rm -rf $(target) $(object)

.PHONY: all install clean
