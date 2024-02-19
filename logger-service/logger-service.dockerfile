FROM alpine:latest

RUN mkdir /app
COPY loggerApp /app
CMD [ "/app/loggerApp" ]

LABEL authors="Mosich-dev"
