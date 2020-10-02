FROM alpine:latest as build
RUN apk add tzdata

FROM scratch as final
COPY /css /css
COPY /html html
COPY /js js
COPY /linux /
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
CMD ["/sklabel_cutting_webservice_linux"]