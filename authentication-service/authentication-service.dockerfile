FROM alpine:latest

RUN mkdir /app

COPY authApp /app 

# COPY ./cmd/api/.env /app

CMD [ "/app/authApp" ]

