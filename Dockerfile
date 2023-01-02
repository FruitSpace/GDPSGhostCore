FROM alpine
EXPOSE 1997
RUN apk add --no-cache tzdata
RUN mkdir /app /core
COPY HalogenGhostCore /app
CMD ["/app/HalogenGhostCore"]