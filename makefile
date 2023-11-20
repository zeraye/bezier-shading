build:
	go build -o bin/bezier-shading 

run: build
	./bin/bezier-shading
