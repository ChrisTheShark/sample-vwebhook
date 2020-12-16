FROM alpine:3.12

ENV KUBECTL_VERSION=v1.19.4

ADD . /predeploy/
WORKDIR /predeploy

RUN apk update && \
    apk add curl openssl && \
    chmod u+x predeploy.sh && \ 
    chmod u+x generate-secret.sh && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl && \
    chmod u+x kubectl && mv kubectl /bin/kubectl

ENTRYPOINT [ "./predeploy.sh" ]