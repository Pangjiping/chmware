test:
	rm -f ./*.out
	go test ./hash -v -run=Test -timeout=10m -cover -coverprofile consistent_test.out


.PHONY: test