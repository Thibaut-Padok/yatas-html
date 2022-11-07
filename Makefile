default: build

test:
	go test ./...

build:
	go build -o bin/yatas-html

update:
	go get -u 
	go mod tidy

install: build
	mkdir -p ~/.yatas.d/plugins/github.com/StanGirard/yatas-html/local/
	mv ./bin/yatas-html ~/.yatas.d/plugins/github.com/StanGirard/yatas-html/local/

release: test
	standard-version
	git push --follow-tags origin main 