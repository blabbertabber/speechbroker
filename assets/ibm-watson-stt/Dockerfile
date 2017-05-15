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

RUN   cd speech-to-text-websockets-python; \
  curl -OL https://raw.githubusercontent.com/blabbertabber/DiarizerServer/master/assets/ibm-watson-stt/sttClient.patch; \
  git apply sttClient.patch
