COMPOSE_FILES=docker-compose.yml
COMPOSE_PROFILES=
COMPOSE_COMMAND=docker-compose

ifeq (, $(shell which $(COMPOSE_COMMAND)))
	COMPOSE_COMMAND=docker compose
	ifeq (, $(shell which $(COMPOSE_COMMAND)))
		$(error "No docker compose in path, consider installing docker on your machine.")
	endif
endif

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq ($(APP_ENV),develop)
	COMPOSE_FILES=docker-compose.yml -f docker-compose-dev.yml
endif

# If the first argument is "log"...
ifeq (log,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif


help:
	@echo "env"
	@echo "==> Create .env file"
	@echo "    Options:"
	@echo "    sort=true   - sort variables alphabetically"
	@echo "    mirror=true - Make .env.example mirror variables in reader.go"
	@echo "    (Options can be combined)"
	@echo ""
	@echo "init"
	@echo "==> init project"
	@echo "setup"
	@echo "==> setup project"
	@echo ""
	@echo "up"
	@echo "==> Create and start containers"
	@echo ""
	@echo "build-up"
	@echo "==> Create and build all containers"
	@echo ""
	@echo "status"
	@echo "==> Show currently running containers"
	@echo ""
	@echo "destroy"
	@echo "==> Down all the containers, keeping their data"
	@echo ""
	@echo "mysql-shell"
	@echo "==> Create an interactive shell for mysql"
	@echo "swagger-generate"
	@echo "==> Create swagger (OpenAPI) document"
	@echo ""
env:
	@echo "Checking if .env exists..."
	@[ -e ./.env ] || cp -v ./.env.example ./.env
	@echo "Extracting environment variables from config/reader.go..."
	@grep -oP 'load\w+\("\K[A-Z_][A-Z0-9_]*(?="\))' internal/config/reader.go \
      | sort -u \
      > /tmp/env_vars.txt
	@echo "Found $$(wc -l < /tmp/env_vars.txt) environment variables in config/reader.go"
	@if [ "$(mirror)" = "true" ]; then \
		echo "mirroring .env.example to match variables in reader.go (preserving existing values)..."; \
		# First, save a list of all variables in the original .env.example \
		grep "^[A-Z]" .env.example | cut -d= -f1 > /tmp/env.example.original_names; \
		# Also save all variable-value pairs for reference \
		grep "^[A-Z]" .env.example > /tmp/env.example.original_pairs; \
		> /tmp/env.example.new; \
		for var in `cat /tmp/env_vars.txt`; do \
			if grep -q "^$$var=" .env.example; then \
				echo "Keeping $$var with existing value"; \
				grep "^$$var=" .env.example >> /tmp/env.example.new; \
			else \
				echo "Adding $$var to .env.example"; \
				echo "$$var=" >> /tmp/env.example.new; \
			fi; \
			# Remove this var from the list of original vars as we've processed it \
			sed -i "/^$$var$$/d" /tmp/env.example.original_names; \
		done; \
		# Now /tmp/env.example.original_names contains only vars that are not in reader.go \
		if [ -s /tmp/env.example.original_names ]; then \
			echo "Removed variables from .env.example:"; \
			for var in `cat /tmp/env.example.original_names`; do \
				val=$$(grep "^$$var=" /tmp/env.example.original_pairs | cut -d= -f2-); \
				echo -e "\033[0;31mRemoved $$var=$$val\033[0m"; \
			done; \
		fi; \
		mv /tmp/env.example.new .env.example; \
		rm -f /tmp/env.example.original_names /tmp/env.example.original_pairs; \
		if [ "$(sort)" = "true" ]; then \
			echo "sorting mirrored .env.example alphabetically..."; \
			grep "^[A-Z]" .env.example | sort > /tmp/env.example.sorted; \
			grep "^#" .env.example > /tmp/env.example.comments 2>/dev/null || true; \
			cat /tmp/env.example.comments /tmp/env.example.sorted > .env.example; \
			rm /tmp/env.example.sorted /tmp/env.example.comments; \
			# Remove empty lines from beginning of file \
			awk 'NF {p=1} p' .env.example > /tmp/env.example.clean; \
			mv /tmp/env.example.clean .env.example; \
		fi; \
	else \
		echo "Checking for missing variables in .env.example..."; \
		for var in `cat /tmp/env_vars.txt`; do \
			if ! grep -q "^$$var=" .env.example; then \
				echo "Adding $$var to .env.example"; \
				echo "$$var=" >> .env.example; \
			fi; \
		done; \
		if [ "$(sort)" = "true" ]; then \
			echo "sorting .env.example alphabetically..."; \
			grep "^[A-Z]" .env.example | sort > /tmp/env.example.sorted; \
			grep "^#" .env.example > /tmp/env.example.comments 2>/dev/null || true; \
			cat /tmp/env.example.comments /tmp/env.example.sorted > .env.example; \
			rm /tmp/env.example.sorted /tmp/env.example.comments; \
			# Let's use a simpler approach - write to a temp file without leading blank lines \
			awk 'NF {p=1} p' .env.example > /tmp/env.example.clean; \
			mv /tmp/env.example.clean .env.example; \
		else \
			echo "Adding variables to .env.example (keeping original order with new vars at end)..."; \
			mv .env.example .env.example.bak; \
			cat .env.example.bak > .env.example; \
			rm .env.example.bak; \
		fi; \
	fi
	@echo "Checking for missing variables in .env..."
	@for var in `cat /tmp/env_vars.txt`; do \
		if ! grep -q "^$$var=" .env; then \
			if grep -q "^$$var=" .env.example; then \
				echo "Adding $$var to .env (with value from .env.example)"; \
				grep "^$$var=" .env.example >> .env; \
			else \
				echo "Adding $$var to .env"; \
				echo "$$var=" >> .env; \
			fi; \
		fi; \
	done
	@if [ "$(sort)" = "true" ]; then \
		echo "sorting .env alphabetically (preserving existing values)..."; \
		grep "^[A-Z]" .env | sort > /tmp/env.sorted; \
		grep "^#" .env > /tmp/env.comments 2>/dev/null || true; \
		cat /tmp/env.comments /tmp/env.sorted > .env; \
		rm /tmp/env.sorted /tmp/env.comments; \
		# Let's use a simpler approach - write to a temp file without leading blank lines \
		awk 'NF {p=1} p' .env > /tmp/env.clean; \
		mv /tmp/env.clean .env; \
	else \
		echo "Adding variables to .env (keeping original order with new vars at end)..."; \
		mv .env .env.bak; \
		cat .env.bak > .env; \
		rm .env.bak; \
	fi
	@rm /tmp/env_vars.txt
	@echo "Environment files updated successfully!"
	@if [ "$(sort)" != "true" ] && [ "$(mirror)" != "true" ]; then \
		echo "Available options:"; \
		echo "- make env sort=true   (sort variables alphabetically)"; \
		echo "- make env mirror=true (make .env.example mirror reader.go variables)"; \
		echo "- make env sort=true mirror=true (combine both options)"; \
	fi

init:
	@echo installing swag
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	@echo installing golangcilint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.6


up:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) up -d

build-up:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) up --build -d

build-no-cache:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) build --no-cache

status:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) ps $(RUN_ARGS)

down:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) down --remove-orphans $(RUN_ARGS)

purge:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) down --remove-orphans --volumes $(RUN_ARGS)

mysql-shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec -u 0 mysql mysql -hmysql -u$(DATABASE_MYSQL_USER) -D$(DATABASE_MYSQL_NAME) -p$(DATABASE_MYSQL_PASSWORD)

redis-shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec -u 0 redis redis-cli

.PHONY: log
log:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) logs -f $(RUN_ARGS)

generate-swagger:
	swag init -q --parseDependency --parseInternal -g router.go -d internal/api/rest

lint:
	golangci-lint run

build-clean:
	rm -rf ./build

build-arm: build-clean
	go generate ./...
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -v -a -ldflags "-w -s" \
            -o build/ ./cmd/...

build-linux: build-clean
	go generate ./...
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -v -a -ldflags "-w -s" \
			-o build/ ./cmd/...

.PHONY: log
shell:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec $(RUN_ARGS) bash

shell-as-root:
	$(COMPOSE_COMMAND) -f $(COMPOSE_FILES) exec -u 0 $(RUN_ARGS) bash

precommit-hook:
	@[ -f ./precommit-hook.sh ] && cp -v ./precommit-hook.sh ./.git/hooks/pre-commit && echo "precommit hook installed" || echo "error: could not find precommit-hook.sh"