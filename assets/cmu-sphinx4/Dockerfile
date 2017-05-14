# This is a Dockerfile of CMU Sphinx 4

# To run:
#   docker run -it blabbertabber/cmu-sphinx4 bash
#   java -Xmx2g -cp sphinx4-5prealpha-src/sphinx4-samples/build/libs/sphinx4-samples-5prealpha-SNAPSHOT.jar edu.cmu.sphinx.demo.speakerid.SpeakerIdentificationDemo
#   java -Xmx2g -cp sphinx4-5prealpha-src/sphinx4-samples/build/libs/sphinx4-samples-5prealpha-SNAPSHOT.jar edu.cmu.sphinx.demo.transcriber.TranscriberDemo 2> /dev/null

FROM fedora

LABEL authors="Brian Cunnie <brian.cunnie@gmail.com>, Brendan Cunnie <saintbrendan@gmail.com>"

CMD ["/bin/bash"]

RUN dnf groupinstall -y "Development Tools"; \
    dnf install -y zip file less vim gradle

# fix: org.apache.http.ssl.SSLInitializationException: /usr/lib/jvm/java-1.8.0-openjdk-1.8.0.121-10.b14.fc25.x86_64/jre/lib/security/cacerts (No such file or directory)
RUN for DIR in /usr/lib/jvm/*/jre/lib/security/; do \
    [ ! -f $DIR/cacerts ] && ln -s /etc/pki/ca-trust/extracted/java/cacerts $DIR; \
  done || true

RUN curl -L "https://downloads.sourceforge.net/project/cmusphinx/sphinx4/5prealpha/sphinx4-5prealpha-src.zip?r=https%3A%2F%2Fsourceforge.net%2Fprojects%2Fcmusphinx%2Ffiles%2Fsphinx4%2F5prealpha%2F&ts=1492168729&use_mirror=iweb"  \
    -o sphinx4-5prealpha-src.zip ; \
  unzip sphinx4-5prealpha-src.zip ; \
  cd sphinx4-5prealpha-src/ ; \
  gradle build ; \
  ( cd sphinx4-samples/ ; gradle build )

RUN rm /sphinx4-5prealpha-src.zip
