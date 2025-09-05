# Commandes prioritaire aux fichiers
.PHONY: build 

# Version de dev 
dev:
	@cd src && go run *.go

# Version compilée
build:
	@cd src && go mod tidy && go build -o ../bin/agent *.go

# Suppression de la version compilée
rm:
	@rm ./bin/agent 
