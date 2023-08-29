
.PHONY: all
all: dsclean aosearch aocrawler

.PHONY: aosearch
aosearch:
	cd cmd/ao3search && go build -o aos main.go && mv aos ${GOPATH}/bin/aos

.PHONY: aocrawler
aocrawler:
	cd cmd/ao3crawler && go build -o aoc main.go && mv aoc ${GOPATH}/bin/aoc

.PHONY: dsclean
dsclean:
	cd cmd/dsclean && go build -o dsc main.go && mv dsc ${GOPATH}/bin/dsc