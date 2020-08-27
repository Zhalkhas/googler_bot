FROM node:latest
RUN git clone https://github.com/postlight/mercury-parser-api/ 
WORKDIR /mercury-parser-api
RUN yarn install
COPY ./reverse_proxy/reverse_proxy reverse_proxy
COPY ./start.sh /start.sh
RUN chmod +x /start.sh
ENTRYPOINT [ "/start.sh" ]