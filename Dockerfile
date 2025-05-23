FROM golang:1.23

WORKDIR /app

# Instala o Modd
RUN go install github.com/cortesi/modd/cmd/modd@latest

# Copia tudo (vamos sobrescrever com volume no compose)
COPY . .

# Exponha a porta usada pela API
EXPOSE 8080

# Comando inicial: roda modd
CMD ["/go/bin/modd", "-f", "modd.conf"]
