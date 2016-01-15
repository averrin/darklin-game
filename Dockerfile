#FROM centurylink/ca-certs
FROM ubuntu
EXPOSE 80
WORKDIR /app
# copy binary into image
COPY core /app/
#COPY .env /app/
#ENTRYPOINT ["./core"]
