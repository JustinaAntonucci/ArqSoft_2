FROM nginx:alpine

# Copia tu archivo de configuración personalizado a nginx
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 8084
CMD ["nginx", "-g", "daemon off;"]
