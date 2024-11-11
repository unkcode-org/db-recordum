# Dockerfile
FROM golang:1.22-alpine

# Instala herramientas necesarias
RUN apk update && apk add --no-cache mysql-client

# Crea directorio de la aplicación
WORKDIR /app

# Copia los archivos de la aplicación
COPY . .

# Compila la aplicación
RUN go build -o db-recordum main.go

# Comando para ejecutar el servicio
CMD ["/app/db-recordum"]
