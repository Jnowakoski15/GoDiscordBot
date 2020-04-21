build:
	docker build . -t team-maker-bot
run:
	docker run  -e DISCORD_TOKEN=<token> team-maker-bot