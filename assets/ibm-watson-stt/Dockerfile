# This is a Dockerfile of IBM 4 + custom patched file
#  - only 1 alternative (i.e. no alternatives, just one version)
#  - enable diarization
#  - don't worry about continuous (streaming)/interim results

# To run:
#   docker run -it blabbertabber/ibm-watson-stt bash

FROM fedora

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN dnf groupinstall -y "Development Tools"

RUN git clone https://github.com/watson-developer-cloud/speech-to-text-websockets-python

RUN   cd speech-to-text-websockets-python; \
  curl -OL https://raw.githubusercontent.com/blabbertabber/DiarizerServer/master/assets/ibm-watson-stt/sttClient.patch; \
  git apply sttClient.patch
