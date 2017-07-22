# This is a Dockerfile of IBM 4 + custom patched file
#  - only 1 alternative (i.e. no alternatives, just one version)
#  - enable diarization
#  - don't worry about continuous (streaming)/interim results

# To run:
#   docker run python ./sttClient.py -credentials 9f6c2cb4-d9d3-49db-96e4-58406a2fxxxx:8rgjxxxxxxxx -model en-US_NarrowbandModel -in <(echo /Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav) -out /tmp/junk


FROM fedora

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN dnf groupinstall -y "Development Tools" && \
    dnf install -y redhat-rpm-config && \
    dnf install -y python && \
    dnf install -y python-devel

RUN git clone https://github.com/watson-developer-cloud/speech-to-text-websockets-python && \
    pip install -r /speech-to-text-websockets-python/requirements.txt

RUN   cd speech-to-text-websockets-python; \
  curl -OL https://raw.githubusercontent.com/blabbertabber/speechbroker/master/assets/ibm-watson-stt/sttClient.patch; \
  git apply sttClient.patch
