## diarizer.com Networking

* diarizer.com has both IPv4 & IPv6 addresses
* currently maps to home.nono.io (73.15.134.22 and 2601:646:100:e8e8::101)
* tcp4/22,80,443,9443 is forwarded appropriately
* tcp6/22,80,443,9443 is allowed
* DiarizerServer listens on 9443
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
sudo dnf install vim git docker nginx python golang htop
```

disable selinux (it's the biggest goddamn pain in the butt)

```
vim /etc/sysconfig/selinux
```

```diff
-SELINUX=enforcing
+SELINUX=permissive
```

```
sudo shutdown -r now
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

Create the key and the CSR
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
+        ssl_certificate "/etc/pki/nginx/diarizer.com.crt";
+        ssl_certificate_key "/etc/pki/nginx/private/diarizer.com.key";
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

Create `/etc/cron.weekly/letsencrypt.sh`
```
#!/bin/bash
set -eux
sudo -u acme_tiny python /usr/local/bin/acme_tiny.py \
    --account-key /etc/pki/letsencrypt.key \
    --csr /etc/pki/nginx/server.csr \
    --acme-dir /var/blabbertabber/acme-challenge > /tmp/signed.crt || exit
wget -O- https://letsencrypt.org/certs/lets-encrypt-x3-cross-signed.pem > /tmp/intermediate.pem
cat /tmp/signed.crt /tmp/intermediate.pem |
    sudo -u acme_tiny tee /etc/pki/nginx/server.crt
sudo systemctl restart nginx.service
```

Make it executable and test
```
sudo chmod +x /etc/cron.weekly/letsencrypt.sh
sudo /etc/cron.weekly/letsencrypt.sh
```

Cleanup
```
rm ~/workspace/acme-tiny/{*.key,*.csr,*.crt}
```

Copy keys (letsencrypt.key & diarizer.com) into LastPass

### Prep Upload Server

Install Docker
```
sudo systemctl enable docker.service
sudo groupadd --system docker
sudo usermod -aG docker cunnie
sudo usermod -aG docker diarizer
sudo shutdown -r now
```

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

Install service
```
sudo cp DiarizerServer /usr/local/bin/
sudo chown diarizer:diarizer /usr/local/bin/DiarizerServer
 # the following is bad; should have other ways to set uid
sudo chown diarizer:diarizer /usr/local/bin/DiarizerServer
sudo chmod 6755 /usr/local/bin/DiarizerServer
sudo cp assets/diarizer.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now diarizer.service
```

Privacy Policy (7 days, prune anything older than 6 days = 24 * 60 * 6 = 8640 minutes). Append the following
line to `/etc/crontab`
```bash
23 0 *  *  *  * diarizer   find /var/blabbertabber/soundFiles/ -name '*-*-*-*' -type d -mmin +8640 -exec rm -rf {} \;
```

Updating service
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
