FROM alpine:3.3

COPY ./_output/httpClient.linux /bin/httpClient
RUN chmod +x /bin/httpClient

ENTRYPOINT ["/bin/httpClient"]
