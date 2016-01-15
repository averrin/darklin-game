#FROM centurylink/ca-certs
FROM ubuntu
EXPOSE 80
WORKDIR /app
# copy binary into image
COPY echo /app/
#COPY .env /app/
#ENTRYPOINT ["./echo"]
