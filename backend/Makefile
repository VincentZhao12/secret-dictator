build: 
	go build -o bin main.go
run: build
	./bin/main
clean:
	rm -rf bin
types:
	rm -rf generated-types/modules/*
	tygo generate
	touch generated-types/types.ts
	cat generated-types/modules/* > generated-types/types.ts
	cp generated-types/types.ts ../web/src/types/index.ts