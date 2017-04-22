# This is a Dockerfile of CMU Sphinx4
#
# To run:
#   docker run -it blabbertabber/cmu-sphinx4-transcriber bash

FROM blabbertabber/cmu-sphinx4

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

RUN curl -OL "https://raw.githubusercontent.com/blabbertabber/DiarizerServer/master/assets/cmu-sphinx4-transcriber/Transcriber.java"

RUN javac Transcriber.java -cp /sphinx4-5prealpha-src/sphinx4-core/build/libs/sphinx4-core-5prealpha-SNAPSHOT.jar
