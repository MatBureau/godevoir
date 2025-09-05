# Image pour la compilation
FROM golang:1.25 AS builder

# Copie de l'application et du Makefile
COPY src /go/src
COPY Makefile /go

# Compilation
RUN make build 

# Image finale
FROM busybox

# Installation des fichiers
COPY src/www /www
COPY --from=builder /go/bin/web /web

# Ouverture du port
EXPOSE 9090

# Démarrage automatique de agent
ENTRYPOINT ["/web"]