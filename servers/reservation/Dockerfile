FROM alpine
RUN apk add --no-cache ca-certificates
COPY reservation ./reservation 
EXPOSE 80
ENTRYPOINT ["/reservation"]