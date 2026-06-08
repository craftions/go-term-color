.PHONY: test coverage cover-html clean

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	findstr /V /C:"/example" coverage.out > coverage_filtered.out
	move /Y coverage_filtered.out coverage.out
	go tool cover -func=coverage.out

cover-html:
	go tool cover -html=coverage.out

clean:
	go clean -testcache
	if exist *.out del /q /f *.out
	if exist *.html del /q /f *.html
