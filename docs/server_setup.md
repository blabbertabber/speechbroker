## diarizer.com Networking

* diarizer.com has both IPv4 & IPv6 addresses
* currently maps to home.nono.io (73.15.134.22 and 2601:646:100:e8e8::101)
* tcp4/22,80,443,9443 is forwarded appropriately
* tcp6/22,80,443,9443 is allowed
* speechbroker listens on 9443
* nginx listens on 443

### URLS:

* <https://diarizer.com:9443/api/v1/upload>
  * creates `/var/blabbertabber/soundFiles/some-guid`
  * creates `/var/blabbertabber/soundFiles/some-guid/meeting.wav`
  * kicks off diarization, i.e.
    ```bash
    docker run blahblah
    ```
  * saves output to /var/blabbertabber/diarizationResults/some-guid

* <https://diarizer.com/some-guid/>

### Directory Structure:

* `/var/blabbertabber/` datadir
    * `soundFiles/some-guid` UploadServer saves `.wav` files here
    * `diarizer/index.html` index for <https://diarizer.com>
    * `diarizer/some-guid` diarizer saves `stdout` here
    * `acme-challenge/` Let's encrypt work files (SSL certification)

### GitHub

Files in https://github.com/cunnie/fedora.nono.io-etc take precedence over files
listed here in `/etc/`

### Prerequisites

```bash
sudo dnf install vim git nginx python golang htop
```

disable selinux (it's the biggest goddamn pain in the butt)

```bash
vim /etc/sysconfig/selinux
```

```diff
-SELINUX=enforcing
+SELINUX=permissive
```

```bash
sudo shutdown -r now
```

Install docker

```bash
sudo dnf config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo
sudo dnf makecache fast
sudo dnf remove docker docker-common; sudo shutdown -r now # if you have previously installed Fedora's version
sudo dnf install docker-ce
sudo systemctl enable docker.service
sudo groupadd --system docker
sudo usermod -aG docker cunnie
sudo usermod -aG docker diarizer
sudo shutdown -r now
```

Disable firewall. Note that we `mask` rather than `disable`. _I believe_ that's
because Docker goes rogue and brings up the firewall when started.

```bash
sudo systemctl stop firewalld; sudo systemctl mask firewalld
```

And, since we are an NTP server, we modify `/etc/chrony.conf`:

```diff
+# When making changes:
+#   sudo systemctl restart chronyd.service
+# Status
+#   echo -e "tracking\nsources\nsourcestats" | chronyc
+# Use google's time servers
+pool time1.google.com iburst
+pool time2.google.com iburst
+pool time3.google.com iburst
+pool time4.google.com iburst

+allow

+clientloglimit 16777216

+log measurements statistics tracking
+
+# latency tweaks (probably not necessary)
+sched_priority 1
+lock_all
```

And then restart:

```bash
sudo systemctl restart chronyd.service
```

The combination of Docker with NTP-server-with-many-clients will exhaust
the connection-tracking that Docker enables, which will disrupt clients with
the following kernel message (`/var/log/messages`):

```
nf_conntrack: table full, dropping packet
```

The fix is tricky; first we need to bump the kernel's number of connections,
but that is not enough: the difficulty is that the boot process attempts to
set the number of connections _before_ the `nf_conntrack` module is loaded, and
fails. When the module is loaded much later as a result of `dockerd` starting,
it defaults to the too-small default setting of 65,536. Our solution is to set both the
`sysctl` setting _and_ force the module to be loaded on boot:

```bash

echo nf_conntrack | sudo tee /etc/modules-load.d/nf_conntrack.conf
echo br_netfilter | sudo tee /etc/modules-load.d/br_netfilter.conf
export SYS=/etc/sysctl.conf; grep conntrack $SYS || ( echo 'net.netfilter.nf_conntrack_max = 524288' | sudo tee -a $SYS )
```

### preparation for `acme-tiny`

Create the user that will update the certificates:
```
sudo adduser \
    --system \
    -c "acme-tiny" \
    -d /var/blabbertabber \
    -M \
    -s /sbin/nologin \
    acme_tiny
```

Create the user that will diarize the meetings:
```
sudo mkdir /var/blabbertabber
sudo adduser \
    --system \
    -c "BlabberTabber Diarizer" \
    -d /var/blabbertabber \
    -M \
    -s /sbin/nologin \
    diarizer
```

Create the "Let's Encrypt" account key (different from the HTTPS key)
and store it on the Filesystem
```
sudo touch /etc/pki/letsencrypt.key
sudo chown acme_tiny:acme_tiny /etc/pki/letsencrypt.key
sudo chmod 600 /etc/pki/letsencrypt.key
openssl genrsa 4096 | sudo -u acme_tiny tee /etc/pki/letsencrypt.key
```

Modify `/etc/.gitignore`
```diff
+ pki/*.key
+ pki/nginx/private/
```

Create the challenge directory
```bash
sudo mkdir /var/blabbertabber/acme-challenge
sudo chown acme_tiny:nginx /var/blabbertabber/acme-challenge
```

Modify `/etc/nginx/nginx.conf`
```diff
+   location /.well-known/acme-challenge/ {
+       alias /var/blabbertabber/acme-challenge/;
+       try_files $uri =404;
+   }
```

Restart nginx
```
sudo systemctl restart nginx.service
```

Download & install `acme-tiny`
```
cd ~/workspace
git clone git@github.com:diafygi/acme-tiny.git
cd acme-tiny
```

Create the key and the CSR; remember to adjust the subjectAltName to whichever server you're configuring.
```
CN=diarizer.com
sudo mkdir -p /etc/pki/nginx/private
sudo chown -R nginx:nginx /etc/pki/nginx
sudo chmod 750 /etc/pki/nginx/private
openssl genrsa 4096 | sudo -u nginx tee /etc/pki/nginx/private/server.key
sudo chmod 440 /etc/pki/nginx/private/server.key
sudo chown -R nginx:diarizer /etc/pki/nginx/private/
openssl req \
  -new \
  -key <(sudo -u nginx cat /etc/pki/nginx/private/server.key) \
  -sha256 \
  -subj "/C=US/ST=California/L=San Francisco/O=BlabberTabber/OU=/CN=${CN}/emailAddress=brian.cunnie@gmail.com" \
  -reqexts SAN \
  -config <(cat /etc/pki/tls/openssl.cnf <(printf "[SAN]\nsubjectAltName=DNS:diarizer.com,DNS:home.nono.io,DNS:home.nono.com,DNS:diarizer.com")) \
  -out server.csr
 # prepare empty certificate file with proper permissions
sudo touch /etc/pki/nginx/server.crt
sudo chown acme_tiny:nginx /etc/pki/nginx/server.crt
sudo chmod 664 /etc/pki/nginx/server.crt
```

Procure the certificate
```
sudo -u acme_tiny \
    python acme_tiny.py \
        --account-key /etc/pki/letsencrypt.key \
        --csr server.csr \
        --acme-dir /var/blabbertabber/acme-challenge/ \
        | sudo -u acme_tiny tee /etc/pki/nginx/server.crt
```

Modify `/etc/nginx/nginx.conf` to use https
```diff
+    server {
+        listen       443 ssl http2 default_server;
+        listen       [::]:443 ssl http2 default_server;
+        server_name  _;
+        root         /var/blabbertabber/diarizationResults;
+
+        ssl_certificate "/etc/pki/nginx/server.crt";
+        ssl_certificate_key "/etc/pki/nginx/private/server.key";
+        ssl_session_cache shared:SSL:1m;
+        ssl_session_timeout  10m;
+        ssl_ciphers HIGH:!aNULL:!MD5;
+        ssl_prefer_server_ciphers on;
+
+        # Load configuration files for the default server block.
+        include /etc/nginx/default.d/*.conf;

+        location / {
+            autoindex on;
+        }
+
+        error_page 404 /404.html;
+            location = /40x.html {
+        }
+
+        error_page 500 502 503 504 /50x.html;
+            location = /50x.html {
+        }
+    }
```

```bash
sudo systemctl enable --now nginx.service
sudo systemctl restart nginx.service
```

Skeletal page to avoid 404:
```bash
sudo mkdir /var/blabbertabber/diarizationResults
sudo chown diarizer /var/blabbertabber/diarizationResults
echo '<html><title>BlabberTabber</title><body><h1>BlabberTabber</h1></body></html' |
    sudo -u diarizer tee /var/blabbertabber/diarizationResults/index.html
```

Move `acme_tiny.py` into an appropriate directory
```
sudo cp ~/workspace/acme-tiny/acme_tiny.py /usr/local/bin/
```

Copy cronjob into place to keep SSL certificates fresh:
```
sudo cp $GOPATH/src/github.com/blabbertabber/speechbroker/assets/cron.weekly.sh \
  /etc/cron.weekly/letsencrypt.sh
```

Make it executable and test
```
sudo /etc/cron.weekly/letsencrypt.sh
```

Cleanup
```
rm ~/workspace/acme-tiny/{*.key,*.csr,*.crt}
```

Copy keys (letsencrypt.key & diarizer.com) into LastPass

### Prep Upload Server

Create directories and run test
```
sudo mkdir -p /var/blabbertabber/soundFiles
sudo chown diarizer /var/blabbertabber/soundFiles
sudo -u diarizer mkdir -p /var/blabbertabber/soundFiles/3426dfcc-fe5f-4686-9279-d997ef9fb0da
cd /var/blabbertabber/soundFiles/3426dfcc-fe5f-4686-9279-d997ef9fb0da
sudo -u diarizer curl -OL https://nono.io/meeting.wav
sudo -u diarizer mkdir /var/blabbertabber/diarizationResults/3426dfcc-fe5f-4686-9279-d997ef9fb0da
sudo -u diarizer \
     docker run \
        --volume=/var/blabbertabber:/blabbertabber \
        --workdir=/speaker-diarization \
        blabbertabber/aalto-speech-diarizer \
            /speaker-diarization/spk-diarization2.py \
                /blabbertabber/soundFiles/3426dfcc-fe5f-4686-9279-d997ef9fb0da/meeting.wav \
                -o /blabbertabber/diarizationResults/3426dfcc-fe5f-4686-9279-d997ef9fb0da/results.txt
```

Download the IBM Bluemix Watson Speech to Text Service Credentials and copy into `/etc/speechbroker`

* Log into Bluemix
* Click on the upper-left hamburger
* **Services &rarr; Watson**
* click on **Speech-to-Text-*xx***
* click on **Service credentials**
* click **View credentials**
* copy credentials into clipboard

```
sudo mkdir /etc/speechbroker
sudo tee /etc/speechbroker/ibm_service_creds.json # paste from clipboard and hit enter+^D
```

Install service
```
cd
mkdir go
go get github.com/blabbertabber/speechbroker/
go build github.com/blabbertabber/speechbroker
go build -o /tmp/ibmjson github.com/blabbertabber/speechbroker/ibmjson
sudo cp speechbroker /tmp/ibmjson /usr/local/bin/
sudo setcap cap_setgid+ep /usr/local/bin/speechbroker
sudo cp assets/diarizer.service /usr/lib/systemd/system/
echo enable diarizer.service | sudo tee /usr/lib/systemd/system-preset/50-diarizer.preset
sudo systemctl daemon-reload
sudo systemctl enable --now --system diarizer.service
```

Copy the transcription/diarization speed factors file into place (if you're
on the test server, the source file name is `speedfactors-test.json`):

```
sudo cp assets/speedfactors.json /etc/speechbroker/speedfactors.json
```

Privacy Policy (7 days, prune anything older than 6 days = 24 * 60 * 6 = 8640 minutes). Append the following
line to `/etc/crontab`
```bash
23   0  *  *  * diarizer   find /var/blabbertabber/soundFiles/ -name '*-*-*-*' -type d -mmin +8640 -exec rm -rf {} \;
23   0  *  *  * diarizer   docker system prune --all --force
```

Fix ``
```bash
export SYS=/etc/sysctl.d/99-sysctl.conf; grep -q conntrack $SYS ||
  ( echo 'net.netfilter.nf_conntrack_max = 524288' | sudo tee -a $SYS )
```


Testing Locally (i.e. not on diarizer.com or test.diarizer.com)

* Create the key & certificate

```
openssl req -new -newkey rsa:2048 -days 3650 -nodes -x509 -keyout /tmp/server.key -out /tmp/server.crt
```

* Run code locally:

```
go run main.go -keyPath=/tmp/server.key -certPath=/tmp/server.crt
```

For Windows, check Brendan's _Daily Learned_ file.