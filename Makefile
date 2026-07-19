up:
	docker compose --profile local up --build -d

down:
	docker compose --profile local down

logs:
	docker compose --profile local logs -f
