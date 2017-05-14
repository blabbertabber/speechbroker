# This is a Dockerfile of IBM 4 + customized BlabberTabber java front-end

# The front-end invokes CMU Sphinx 4's transcription features
# It depends on the generic CMU Sphinx 4 Docker image

# To run:
#   docker run -it blabbertabber/cmu-sphinx4-transcriber bash

FROM fedora

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN dnf groupinstall -y "Development Tools"

RUN git clone https://github.com/watson-developer-cloud/speech-to-text-websockets-python

# RUN   cd speech-to-text-websockets-python; \
#    curl -L https://gist.githubusercontent.com/cunnie/0825743040dbf9d889b50a28054354b3/raw/08ed8b680e0e82d4626884feb484fd31d6db641d/sttclient.patch -o sttClient.patch; \
#    git apply sttClient.patch
