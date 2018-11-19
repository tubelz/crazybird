# build stage
FROM rennomarcus/macaw:latest
# some information about the docker image
LABEL Name="Crazybird" Version="0.2" Maintainer="Marcus Renno <me@rennomarcus.com>"

ARG appName
ARG repository

WORKDIR /go/src/${repository}/${appName}