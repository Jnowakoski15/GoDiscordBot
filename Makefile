build:
	docker build . -t jnowakoski/go-team-maker-bot:scratch
pushs:
	docker push jnowakoski/go-team-maker-bot:scratch
push:
	docker push jnowakoski/go-team-maker-bot
run:
	docker run  -e DISCORD_TOKEN=NzAyMTU5ODU2ODUzMDU3NTg4.Xp8UPQ.UQ8yJts_N3Qwxigw_i7hKpZwcjE jnowakoski/go-team-maker-bot