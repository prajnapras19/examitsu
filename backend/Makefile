# Source: https://www.linkedin.com/pulse/go-database-migrations-made-easy-example-using-mysql-tiago-melo/
# Modified to suffice our needs

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Version - optionally used on the goto command
V?=

# Number of migrations - optionally used on up and down commands
N?=

# Check environment variables for all commands except mysql-migrate-setup
ifneq ($(MAKECMDGOALS),mysql-migrate-setup)
    ifeq ($(MYSQL_USER),)
      $(error MYSQL_USER is not set)
    endif
    ifeq ($(MYSQL_PASSWORD),)
      $(error MYSQL_PASSWORD is not set)
    endif
    ifeq ($(MYSQL_HOST),)
      $(error MYSQL_HOST is not set)
    endif
    ifeq ($(MYSQL_DATABASE),)
      $(error MYSQL_DATABASE is not set)
    endif
    ifeq ($(MYSQL_MIGRATIONS_FOLDER),)
      $(error MYSQL_MIGRATIONS_FOLDER is not set)
    endif
endif

MYSQL_PORT ?= 3306

MYSQL_DSN = $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))/$(MYSQL_DATABASE)

mysql-migrate-setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing migrate command..."; go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; fi

mysql-migrate-up: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) up $(N)

mysql-migrate-down: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) down $(N)

mysql-migrate-to-version: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) goto $(V)

mysql-drop-db: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) drop

mysql-force-version: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) force $(V)

mysql-migration-version: mysql-migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path $(MYSQL_MIGRATIONS_FOLDER) version
