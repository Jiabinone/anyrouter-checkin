.PHONY: dev build fmt lint clean

dev:
	@npm run dev

build:
	@npm run build

fmt:
	@npx prettier --write "src/**/*.{vue,ts,js,css}"

lint:
	@npm run lint
	@npm run type-check

clean:
	@rm -rf dist/ node_modules/.cache/
