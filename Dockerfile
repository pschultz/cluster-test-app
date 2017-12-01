FROM busybox

COPY server /server

CMD ["/server"]
