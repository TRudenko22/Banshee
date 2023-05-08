FROM fedora:latest
WORKDIR /app

RUN mkdir /data && chmod 777 /data
COPY ./banshee/ /app/
RUN dnf update -y && dnf install go -y
RUN go build .

EXPOSE 587
CMD ["/app/banshee"]
