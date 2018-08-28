FROM alpine:latest
EXPOSE 3611 3612
RUN apk add --update --no-cache ca-certificates
RUN apk add --update --no-cache tzdata
ENV TZ America/New_York
WORKDIR /
ADD bin/* /dauth/bin/
ADD script/* /dauth/script/
ADD dauth_config.yaml /
CMD ["/bin/dauth", "serve"]
