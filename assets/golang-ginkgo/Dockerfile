# This is a Dockerfile of golang + ginkgo + gomega + node + jest test tools

FROM golang

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN go get github.com/onsi/ginkgo/ginkgo; \
    go get github.com/onsi/gomega; \
    apt-get update; \
    curl -sL https://deb.nodesource.com/setup_12.x | bash -; \
    apt-get install -y nodejs; \
    apt-get install -y gcc g++ make; \
    curl -sL https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add -; \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list; \
    apt-get update; \
    apt-get install -y yarn; \
    yarn global add jest --dev
