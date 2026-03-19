.PHONY: all scan report plan clean tidy

init:
	go run main.go init

scan:
	go run main.go scan --file scans.json

report:
	go run main.go report

plan:
	go run main.go plan

clean:
	rm -f migration.json purgomatic.db
