all: cli/cli collector/collector

cli/cli: cli/*.go
	cd cli && go build

collector/collector: collector/*.go
	cd collector && go build

clean:
	rm -vf cli/cli collector/collector

