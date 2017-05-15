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

## Developers: How to Update Diarizer Server

### 1. Updating the Golang-based Diarizer Server


First, test it on the test server:

```bash
ssh -i /c/Users/saint/.ssh/id_github saintbrendan@test.diarizer.com
cd $GOPATH/src/github.com/blabbertabber/DiarizerServer
git pull -r
go build
sudo systemctl stop diarizer.service
sudo -u diarizer ./DiarizerServer
 # run BlabberTabber, upload file, check output -- .txt files there?
 # if not, debug and repeat
sudo cp DiarizerServer /usr/local/bin/
sudo systemctl start diarizer.service
```

Then install and test on the production server (only change from above is first line, where we ssh into the server)

```bash
ssh -i /c/Users/saint/.ssh/id_github saintbrendan@diarizer.com
cd $GOPATH/src/github.com/blabbertabber/DiarizerServer
git pull -r
go build
sudo systemctl stop diarizer.service
sudo -u diarizer ./DiarizerServer
 # run BlabberTabber, upload file, check output -- .txt files there?
 # if not, debug and repeat
sudo cp DiarizerServer /usr/local/bin/
sudo systemctl start diarizer.service
```

### 2. Updating the HTML/CSS/JS

First, test. Webstorm should be configured to point to the test server by default, so you it should automatically
publish.

Here's a good url to test: <https://test.diarizer.com/?m=test>

When everything works, you can publish to production:
in Webstorm: **Tools &rarr; Deployment &rarr; Upload to ... &rarr; diarizer.com**

If you want to upload a wav file _without_ using the Android client, here's an example:

```
curl -F "meeting.wav=@/Users/cunnie/Google Drive/BlabberTabber/ICSI-diarizer-sample-meeting.wav" https://test.diarizer.com:9443/api/v1/upload
```