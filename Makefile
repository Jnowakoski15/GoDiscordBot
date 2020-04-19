build:
	docker build . -t go-dock
run:
	docker run -p 3000:3000 go-dock