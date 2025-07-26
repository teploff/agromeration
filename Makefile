DC=docker-compose

up:
	$(DC) up -d pg-auth

down:
	$(DC) down

logs:
	$(DC) logs -f pg-auth

ps:
	$(DC) ps

clean:
	$(DC) down -v