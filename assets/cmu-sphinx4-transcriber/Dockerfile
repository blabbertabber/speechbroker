# This is a Dockerfile of CMU Sphinx 4 + customized BlabberTabber java front-end

# The front-end invokes CMU Sphinx 4's transcription features
# It depends on the generic CMU Sphinx 4 Docker image

# To run:
#   docker run -it blabbertabber/cmu-sphinx4-transcriber bash

FROM blabbertabber/cmu-sphinx4

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN curl -OL "https://raw.githubusercontent.com/blabbertabber/DiarizerServer/master/assets/cmu-sphinx4-transcriber/Transcriber.java"

RUN javac Transcriber.java -cp /sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar
