## diarizer.blabbertabber.com nginx configuration

* diarizer.blabbertabber.com has both IPv4 & IPv6 addresses
* currently maps to home.nono.io (73.15.134.22 and 2601:646:100:e8e8::101)
* tcp4/80,443,9443 is forwarded appropriately
* tcp6/22,80,443,9443 is allowed
* UploadServer listens on 9443
* nginx listens on 443

### preparation for `nginx`

#### Create nginx.conf



One-time set-up:

```
sudo mkdir /opt/blabbertabber
sudo adduser \
    --system \
    -c "BlabberTabber Diarizer" \
    -d /opt/blabbertabber \
    -M \
    -s /sbin/nologin \
    diarizer
```

Get certs:

```

```

URLS:

* <https://diarizer.blabbertabber.com:9443/api/v1/upload>
  * creates `/opt/blabbertabber/UploadServer/some-guid`
  * creates `/opt/blabbertabber/UploadServer/some-guid/meeting.wav`
  * kicks off diarization, i.e.
    ```bash
    docker run
    ```

* <https://diarizer.blabbertabber.com/some-guid/>

```
mkdir /opt/blabbertabber/nginx
