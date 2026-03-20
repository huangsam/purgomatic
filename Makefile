.PHONY: all audit init clean

all: audit

init:
	go run main.go init

audit:
	go run main.go audit --file scans.json

clean:
	rm -f migration.json purgomatic.db
