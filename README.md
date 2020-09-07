# Ninja proxy
<p align="center">
  <img alt="ninja cat image" src="https://i.ibb.co/S54DPbc/ninja.png" />
</p>

## Generate temporary URLs with a stateless service
This project grew out of the need to share authenticated proxies with third-party services, so that they can all operate with the same IP. Several proxy-providers provide mechanisms for this, but they usually require sending your username and password to the third-party and this can be very risky. Ideally you want some temporary URL that opaquely routes the traffic from the third-party to your proxy provider, and that is exactly what this project provides.

### Quick start
Clone this repo and build the project. Assuming that golang is installed, you build the server using

`make binary`

that will put it into `./bin/ninja-proxy`.

You can install the helper utilities by running

```bash
source activate.sh
```

and load a key into your terminal with

```bash
KEY=$(ninja-key)
```

Now you can run the service (in the background) using

```bash
bin/ninja-proxy -key $KEY &
```
which by default will run on 0.0.0.0:7777.

To generate a link for the service call

```bash
link=$(ninja-link 120 $KEY http://username:password@httpbin.org:80 --headers test-header=foo another=bar)
```

```bash
$ curl $link/headers
{
  "headers": {
    "Accept": "*/*",
    "Another": "bar",
    "Authorization": "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
    "Host": "httpbin.org",
    "Test-Header": "foo",
    "User-Agent": "curl/7.72.0",
    "X-Amzn-Trace-Id": "Root=1-5f56788f-d9513fb01690f9080adeb528"
  }
}
```
Notice how the username and password are automatically added to the `Authorization` header, you can do the same thing with a proxy URL and it will just work.

The generated link contains a username and header field with all of this information encrypted so that it can't be read. It looks quite messy, but it works.

#### Proxies
There's not very much else to say, use exactly the same procedure as above to wrap  your proxy URL and use it as normal

```bash
proxy_link=$(ninja-link 120 $KEY "http://$PROXY_USER:$PROXY_PASS@$PROXY_HOST:$PROXY_PORT")
curl -x $proxy_link -k http://httpbin.org/ip
```

### Docker image
A dockerfile has been included with this repository and can be built using `make docker`. Then you can run this using

```
NINJA_PORT=7777
docker run --rm -d -p $NINJA_PORT:7777 ninja-proxy
```

### Limitations
This is a young project created for a specific problem at hand. You can use this to proxy both HTTP and HTTPS connections, but right now it only supports HTTP proxies or URLs using HTTP/1.x.
