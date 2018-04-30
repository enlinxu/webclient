FROM alpine:3.3

COPY ./_output/webclient.linux /bin/webclient
RUN chmod +x /bin/webclient

ENTRYPOINT ["/bin/webclient"]
