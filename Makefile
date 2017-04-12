SHELL = /bin/bash
TARGETS = mabread

all: $(TARGETS)

$(TARGETS): %: cmd/%/main.go
	go build -o $@ $<

clean:
	rm -f $(TARGETS)
