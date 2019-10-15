FROM alpine:3.9

ENV USER=appuser
ENV HOME=/home/$USER

RUN adduser -S $USER

# Copy files to container
COPY target/server/app $HOME/app

# Configure permissions
RUN chown -R $USER $HOME

USER $USER
WORKDIR /home/$USER
CMD ["./app"]
