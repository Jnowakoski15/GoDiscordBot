build:
	docker build . -t jnowakoski/go-team-maker-bot:scratch
pushs:
	docker push jnowakoski/go-team-maker-bot:scratch
push:
	docker push jnowakoski/go-team-maker-bot
run:
	docker run  -e DISCORD_TOKEN=<token> jnowakoski/go-team-maker-bot