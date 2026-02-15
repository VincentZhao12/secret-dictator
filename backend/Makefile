build-web:
	cd ../web && npm run build
	mkdir -p static
	cd ../backend && cp -r ../web/dist/* ./static
	

build: build-web
	go build -o bin main.go

run: build
	set -a && [ -f .env.dev ] && . ./.env.dev && set +a && ./bin/main

dev:
	set -a && [ -f .env.dev ] && . ./.env.dev && set +a && go run main.go

run-prod: build
	set -a && [ -f .env.prod ] && . ./.env.prod && set +a && ./bin/main

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