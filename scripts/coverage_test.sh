go test -json ./... -coverprofile=coverprofile.tmp -coverpkg=./... | jq; \
	cat coverprofile.tmp | grep -v _mock.go | grep -v _easyjson.go | grep -v .pb.go | grep -v _grpc.go > coverprofile.out ; \
	rm coverprofile.tmp ; \
	go tool cover -func=coverprofile.out  ; \
	go tool cover -html=coverprofile.out
