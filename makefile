build:
	go build -o bin/bezier-shading cmd/bezier-shading/main.go 

run: build
	./bin/bezier-shading
