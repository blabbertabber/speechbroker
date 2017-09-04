## Running Tests

Install ginkgo and gomega if you haven't already:
```bash
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
```

Reformat & run ginkgo:
```bash
gofmt -w .
ginkgo -r .
```

## How to Update Diarizer Server

### 1. Updating the Golang-based Diarizer Server


First, test it on the test server:

```bash
ssh -i /c/Users/saint/.ssh/id_github saintbrendan@test.diarizer.com
cd $GOPATH/src/github.com/blabbertabber/speechbroker
git pull -r
go build
sudo setcap cap_setgid+ep speechbroker
sudo systemctl stop diarizer.service
sudo -u diarizer ./speechbroker \
    -ibmServiceCredsPath=/etc/speechbroker/ibm_service_creds.json \
    -speedfactorsPath=/etc/speechbroker/speedfactors.json
 # run BlabberTabber, upload file, check output -- .txt files there?
 # if not, debug and repeat
sudo cp speechbroker /usr/local/bin/
sudo setcap cap_setgid+ep /usr/local/bin/speechbroker
sudo systemctl start diarizer.service
```

Then install and test on the production server (identical instructions as the test server's, with the exception of
the first line, where we ssh into the server):

```bash
ssh -i /c/Users/saint/.ssh/id_github saintbrendan@diarizer.com
cd $GOPATH/src/github.com/blabbertabber/speechbroker
git pull -r
go build
sudo setcap cap_setgid+ep speechbroker
sudo systemctl stop diarizer.service
sudo -u diarizer ./speechbroker \
    -ibmServiceCredsPath=/etc/speechbroker/ibm_service_creds.json \
    -speedfactorsPath=/etc/speechbroker/speedfactors.json
 # run BlabberTabber, upload file, check output -- .txt files there?
 # if not, debug and repeat
sudo cp speechbroker /usr/local/bin/
sudo setcap cap_setgid+ep /usr/local/bin/speechbroker
sudo systemctl start diarizer.service
```

### 2. Updating the HTML/CSS/JS

First, test. Webstorm should be configured to point to the test server by default, so your changes should be
automatically, instantaneously, published there.

Here are several good urls to test (8, 10, 12, 16 speaker colors, respectively):

* <https://test.diarizer.com/?m=test>
* <https://test.diarizer.com/?m=test-10>
* <https://test.diarizer.com/?m=test-12>
* <https://test.diarizer.com/?m=test-16>

Also, to view the _raw files_ within a directory, remove the `?m=` from the URL (or the `?meeting=`),
e.g.:

* <https://test.diarizer.com/test>

When everything works, you can publish to production:
in Webstorm: **Tools &rarr; Deployment &rarr; Upload to ... &rarr; diarizer.com**

If you want to upload a wav file _without_ using the Android client, here's an example:

```
curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://test.diarizer.com:9443/api/v1/upload
```
