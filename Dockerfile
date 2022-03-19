FROM golang
WORKDIR /pzem-exporter
COPY . .
RUN go build -o PZEM_exporter .

FROM alpine
COPY --from=0 /pzem-exporter/PZEM_exporter .
ENTRYPOINT ["./PZEM_exporter"]
