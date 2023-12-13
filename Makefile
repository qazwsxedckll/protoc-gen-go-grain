.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative --plugin=protoc-gen-go-grain=protoc-gen-go-grain --go-grain_out=. --go-grain_opt=paths=source_relative testdata/hello.proto
