# This is a Dockerfile of golang + ginkgo + gomega + node + jest test tools

FROM golang

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN go get github.com/onsi/ginkgo/ginkgo; \
    go get github.com/onsi/gomega; \
    apt update; \
    curl -sL https://deb.nodesource.com/setup_8.x | bash -; \
    apt install -y nodejs; \
    npm install --save-dev --global jest