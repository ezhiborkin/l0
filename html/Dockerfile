# Используем базовый образ Nginx
FROM nginx:alpine

# Копируем файлы вашего веб-приложения в контейнер
COPY index.html /usr/share/nginx/html/index.html
COPY /nginx-config/nginx.conf /etc/nginx/nginx.conf

RUN apk add --update nano