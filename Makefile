lints:
	$(GOPATH)/bin/golangci-lint run ./... --enable-all
scopelint:
	$(GOPATH)/bin/golangci-lint -v run ./... --disable-all -E scopelint --fix
