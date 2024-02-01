FROM alpine:latest

RUN mkdir /app
COPY authApp /app
CMD [ "/app/authApp" ]

LABEL authors="Mosich-dev"
