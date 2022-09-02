FROM alpine
EXPOSE 1997
RUN mkdir /app /core
COPY HalogenGhostCore /app
CMD ["/app/HalogenGhostCore"]