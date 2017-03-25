## diarizer.blabbertabber.com nginx configuration

* diarizer.blabbertabber.com has both IPv4 & IPv6 addresses
* currently maps to home.nono.io (73.15.134.22 and 2601:646:100:e8e8::101)
* tcp4/80,443,9443 is forwarded appropriately
* tcp6/22,80,443,9443 is allowed
* UploadServer listens on 9443
* nginx listens on 443

URLS:

* <https://diarizer.blabbertabber.com:9443/api/v1/upload>
  * creates `/var/blabbertabber/UploadServer/some-guid`
  * creates `/var/blabbertabber/UploadServer/some-guid/meeting.wav`
  * kicks off diarization, i.e.
    ```bash
    docker run blahblah
    ```
  * saves output to /var/blabbertabber/diarizer/some-guid

* <https://diarizer.blabbertabber.com/some-guid/>

Directory Structure:

* `/var/blabbertabber/` datadir
    * `UploadServer/some-guid` UploadServer saves `.wav` files here
    * `diarizer/index.html` index for <https://diarizer.blabbertabber.com>
    * `diarizer/some-guid` diarizer saves `stdout` here
    * `acme-challenge/` Let's encrypt work files (SSL certification)

### preparation for `acme-tiny`

Create the user that will update the certificates (nginx group to read key):
```
sudo adduser \
    --system \
    -c "acme-tiny" \
    -d /var/blabbertabber \
    -M \
    -s /sbin/nologin \
    -g nginx \
    acme_tiny
```

Modify `/etc/.gitignore`
```diff
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
CN=diarizer.blabbertabber.com
cd ~/workspace
git clone git@github.com:diafygi/acme-tiny.git
cd acme-tiny
```

Create the key and the CSR
```
sudo mkdir -p /etc/pki/nginx/private
sudo chown -R nginx:nginx /etc/pki/nginx
sudo chmod 750 /etc/pki/nginx/private
CN=diarizer.blabbertabber.com
openssl req \
  -new \
  -keyout $CN.key \
  -newkey rsa:4096 \
  -nodes \
  -sha256 \
  -subj "/C=US/ST=California/L=San Francisco/O=BlabberTabber/OU=/CN=${CN}/emailAddress=brian.cunnie@gmail.com/subjectAltName=DNS=diarizer.blabbertabber.com,DNS=home.nono.io" \
  -out $CN.csr
 # store $CN.key in LastPass
sudo mv $CN.key /etc/pki/nginx/private/
sudo chown nginx:nginx /etc/pki/nginx/private/$CN.key
sudo chmod 440 /etc/pki/nginx/private/$CN.key
sudo touch /etc/pki/nginx/$CN.crt
sudo chown acme_tiny:nginx /etc/pki/nginx/$CN.crt
sudo chmod 664 /etc/pki/nginx/$CN.crt
```

Procure the certificate
```
sudo -u acme_tiny \
    python acme_tiny.py \
        --account-key /etc/pki/nginx/private/$CN.key \
        --csr $CN.csr \
        --acme-dir /var/blabbertabber/acme-challenge/ \
        | sudo -u acme_tiny tee /etc/pki/nginx/$CN.crt
```

Modify nginx to use https
```
sudo vim /etc/nginx/nginx.conf
```

Set up cron
```
# TBD
```

Cleanup
```
rm ~/workspace/acme-tiny/{*.key,*.csr,*.crt}
```

### preparation for `nginx`

One-time set-up:
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

```
mkdir /var/blabbertabber/nginx
```
