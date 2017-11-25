# This is a Dockerfile of golang + ginkgo + gomega test tools

FROM golang

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN go get github.com/onsi/ginkgo/ginkgo; \
    go get github.com/onsi/gomega; \
    apt update; \
    apt install -y node npm; \
    npm install --save-dev jest --global