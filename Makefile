lints:
	$(GOPATH)/bin/golangci-lint run ./... --enable-all
scopelint:
	$(GOPATH)/bin/golangci-lint -v run ./... --disable-all -E scopelint --fix
revive:
	$(GOPATH)/bin/revive -config ./revive.toml -formatter friendly ./...
revive-no-tests:
	$(GOPATH)/bin/revive -config ./revive.toml \
	-exclude collections/concurrent_map_benchmark_test.go \
	-exclude collections/concurrent_map_test.go \
	-exclude collections/concurrent_set_test.go \
	-exclude collections/concurrent_set_benchmark_test.go \
	-exclude collections/concurrent_linked_list_test.go \
	-exclude collections/list_item_test.go \
	-exclude caches/lru_test.go \
	-exclude caches/lru_benchmark_test.go \
	-exclude caches/entity_list_test.go \
	-exclude caches/lru_entity_test.go \
    -formatter friendly ./...
