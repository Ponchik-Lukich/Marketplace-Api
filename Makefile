dump_db:
	docker-compose exec -T postgres pg_dump "dbname=avito user=user password=password" >./pkg/storage/migrations/tables.sql