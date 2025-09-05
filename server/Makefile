# Variables
COMPOSE         ?= docker compose      
DOCKER_DIR      ?= docker
COMPOSE_FILE    ?= docker-compose.yml     

# Commandes prioritaire aux fichiers
.PHONY: build compose compose-down compose-logs compose-ps

# Version de dev 
dev:
	@cd src && go run *.go

# Version compilée
build:
	@cd src && go mod tidy && go build -o ../bin/web *.go

# Suppression de la version compilée
rm:
	@rm ./bin/web 

# Version image docker
image:
	@docker build -t web:latest . --no-cache

# Démarrage de la version docker
run:
	@docker run -d --name dev --rm -p 9090:9090 web:latest

# Démarrage de la version docker
stop:
	@docker stop dev

# Démarrage stack
up:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) up -d

# Arrêt de la stack
down:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) down

# Journaux de la stack
logs:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) logs -f

# Liste des processus de la stack
ps:
	@$(COMPOSE) --project-directory $(DOCKER_DIR) -f $(DOCKER_DIR)/$(COMPOSE_FILE) ps

compose-down: down

compose-up: up