.PHONY: all
all: cli collector api

.PHONY: cli
cli:
	cd cli && go build

.PHONY: collector
collector:
	cd collector && go build

.PHONY: api
api:
	cd api && go build

.PHONY: clean
clean:
	rm -vf cli/cli collector/collector api/api
