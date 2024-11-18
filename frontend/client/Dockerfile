# Etapa de construcción
FROM node:16-alpine AS builder

WORKDIR /app

# Copia el archivo package.json y package-lock.json si existe
COPY package*.json ./

# Instala las dependencias
RUN rm -rf node_modules && npm ci

# Copia el resto de los archivos
COPY . .

# Construye la aplicación
RUN npm run build

# Etapa de producción
FROM nginx:alpine

# Copia los archivos compilados de la etapa de construcción
COPY --from=builder /app/build /usr/share/nginx/html

# Exponer el puerto 80 para nginx
EXPOSE 80

# Comando de inicio para nginx
CMD ["nginx", "-g", "daemon off;"]
