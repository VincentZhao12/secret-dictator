build-web:
	cd ../web && npm run build
	mkdir -p static
	cd ../backend && cp -r ../web/dist/* ./static
	
	
build: build-web
	go build -o bin main.go

run: build
	export ENV="development" && ./bin/main

dev: build-web
	export ENV="development" && go run main.go

clean:
	rm -rf bin
	rm -rf ../web/dist
	rm -rf static

types:
	rm -rf generated-types/modules/*
	tygo generate
	touch generated-types/types.ts
	cat generated-types/modules/* > generated-types/types.ts
	cp generated-types/types.ts ../web/src/types/index.ts