FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY prometheus-telegram-bot /usr/bin/bot
#RUN ls -lh /usr/bin/prometheus-telegram-bot
COPY ./template template 
COPY ./mapping mapping 
EXPOSE 9087
ENTRYPOINT ["bot"]

