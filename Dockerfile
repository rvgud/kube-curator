FROM debian
COPY ./app /app
ENTRYPOINT /app -f /config/config.yaml