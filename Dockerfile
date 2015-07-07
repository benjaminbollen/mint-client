# Pull base image.
FROM eris/base

# Set the env variables to non-interactive
ENV DEBIAN_FRONTEND noninteractive
ENV DEBIAN_PRIORITY critical
ENV DEBCONF_NOWARNINGS yes
ENV TERM linux
RUN echo 'debconf debconf/frontend select Noninteractive' | debconf-set-selections

# grab deps (gmp)
RUN apt-get update && \
  apt-get install -y --no-install-recommends \
    libgmp3-dev && \
  rm -rf /var/lib/apt/lists/*

# set the repo and install tendermint
ENV repo /go/src/github.com/eris-ltd/mint-client
ADD . $repo
WORKDIR $repo
RUN go install ./...

# grab eris-keys
ENV keys_repo $GOPATH/src/github.com/eris-ltd/eris-keys
RUN git clone https://github.com/eris-ltd/eris-keys $keys_repo
WORKDIR $keys_repo
RUN git checkout develop
RUN go install
WORKDIR $repo

ADD ./test.sh /test.sh
RUN chown $USER:$USER /test.sh

# set user
USER $USER
ENV TMROOT /home/eris/.eris/
WORKDIR /home/eris

CMD ["/test.sh"]