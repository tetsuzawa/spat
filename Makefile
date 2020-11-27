PROGRAMS := \
bin/dxxconv \
bin/fadein-fadeout \
bin/make_fadein_fadeout_filter_fourier \
bin/make_pinknoise \
bin/overlap-add

.PHONY: all clean

all: ${PROGRAMS}

install: all
	install ${PROGRAMS} ${HOME}/local/bin

clean:
	rm -rf bin

lint:
	go fmt -x ./...

bin/dxxconv: *.go cmd/dxxconv/*.go
	go build -o $@ ./cmd/dxxconv/

bin/fadein-fadeout: *.go cmd/fadein-fadeout/*.go
	go build -o $@ ./cmd/fadein-fadeout/

bin/make_fadein_fadeout_filter_fourier: *.go cmd/make_fadein_fadeout_filter_fourier/*.go
	go build -o $@ ./cmd/make_fadein_fadeout_filter_fourier/

bin/make_pinknoise: *.go cmd/make_pinknoise/*.go
	go build -o $@ ./cmd/make_pinknoise/

bin/overlap-add: *.go cmd/overlap-add/*.go
	go build -o $@ ./cmd/overlap-add/
