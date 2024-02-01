FROM alpine:latest

RUN mkdir /app
COPY brokerApp /app
CMD [ "/app/brokerApp" ]

LABEL authors="Mosich-dev"
