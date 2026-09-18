###############################################################################
# BEGIN build-stage
# Compile the binary
FROM --platform=$BUILDPLATFORM docker.io/library/golang:1.27.1 AS build-stage

ARG BUILDPLATFORM
ARG TARGETARCH

WORKDIR /app

COPY . ./

RUN GOOS=linux GOARCH="${TARGETARCH}" hack/build.sh

#
# END build-stage
###############################################################################

###############################################################################
# BEGIN final-stage
# Create final docker image
FROM scratch AS final-stage

COPY --from=build-stage /app/bin/netrouse /netrouse

USER 65534:65534

WORKDIR /data
VOLUME /data

ENTRYPOINT ["/netrouse"]

CMD ["server"]

#
# END final-stage
###############################################################################
