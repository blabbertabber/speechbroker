### Diarizer Proxy

This is a Cloud Foundry application whose purpose is to proxy requests to the diarizer.

***[You probably don't need this & are better left off ignoring this directory.]***

I have only one IP address, and only one port 80, and only one port 443. I'm using them
for my Cloud Foundry installation, so I forward traffic to the diarizer
via this wonderfully small Go program.

Here's how to install the program:

```sh
cf api api.cf.nono.io
cf target -o org -s space
cf create-private-domain org diarizer.com
cf push
cf create-security-group diarizer.nono.io security-group-diarizer.json
cf bind-security-group diarizer.nono.io org --space space
cf restart diarizer-proxy
```