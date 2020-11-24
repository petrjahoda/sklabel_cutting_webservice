FROM alpine:latest as build
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Europe/Prague /etc/localtime

FROM scratch as final
COPY /css /css
COPY /html html
COPY /js js
COPY /linux /
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/localtime /etc/localtime
CMD ["/sklabel_cutting_webservice"]