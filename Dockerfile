FROM scratch
WORKDIR /srv
COPY ./googler_bot googler_bot
# RUN go build -o googler_bot main.go
ENTRYPOINT ["/srv/googler_bot" ]