# nano-service

nano-service is an example of building a platform that allows rapid
deployment of web application code.

The server-side script that runs is JavaScript, and 
executes in the context of the request, returning the evaluated result.

## Deploying code

As an example, let's say we wanted to create a friendly service that
responds with a cheerful greeting, we can write the script like this:

```javascript
var result = "Hello World!";

result;
```

To deploy this code, we need to HTTP POST it to the 
nano-service `/deploy/` endpoint.

The format of the URL is `/deploy/{name}`

In this example we specify `helloworld` as the name 
(the last part of the path). Here's an example using 
[Fiddler](http://www.telerik.com/fiddler)

```
POST http://localhost:8181/deploy/helloworld HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 41

var result = "Hello World!";

result;
```

The response contains the hash of the function created:

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 20:57:06 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

h7vkaYrOMALxdos6hp2TNSX0xPZIMON3ojkE9Q1sUPw
```

## Running code

We can run the code using the `/run/` endpoint.

The format of the URL is `/run/{name}/{hash}[?{key1}={value1}&{key2}={value2}]`

```
GET http://localhost:8181/run/helloworld/h7vkaYrOMALxdos6hp2TNSX0xPZIMON3ojkE9Q1sUPw HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 0
```

Which returns the result:

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 20:58:00 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

Hello World!
```

## Parameters

Pass request parameters in the query string. These are then exposed to
the script as members of the global `$params` object.

In this example we take two parameters passed to the script, `val1`
and `val2`, add them together and return the result.

The JavaScript looks like this:

```javascript
Number($params.val1) + Number($params.val2);
```

So we deploy this code:

```
POST http://localhost:8181/deploy/add HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 44

Number($params.val1) + Number($params.val2);
```

And get back the hash of the validated code:

```
HTTP/1.1 200 OK
Date: Fri, 01 Jul 2016 22:00:13 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

NO2gb4GKaDoUtBccf5Ail3gEB0O9sZLdwWWaidlxBGU
```

We then execute it, passing `12` as the value of `val1`, and `34` as
the value of `val2`:

```
GET http://localhost:8181/run/add/NO2gb4GKaDoUtBccf5Ail3gEB0O9sZLdwWWaidlxBGU?val1=12&val2=34 HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 0
```

The response is the value `46`:

```
HTTP/1.1 200 OK
Date: Fri, 01 Jul 2016 22:01:19 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

46
```

## Simple HTTP

The nano-service scripting runtime has some simple synchronous http 
functions for `GET` and `POST` operations:

```javascript
var result = $get(url);
```

```javascript
var result = $post(url, "application/json", JSON.stringify(obj))
```

Here's an example which gets the contents of the specified url and
returns it.

`/deploy/` the code:

```
POST http://localhost:8181/deploy/get HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 16

$get($params.url)
```

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 21:01:06 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

rV00-kyEgoZfZa8D9MOxeXqqFiItt0_l61WD_n9Rbm4
```

`/run/` the code, passing an RSS feed url as the parameter:

```
GET http://localhost:8181/run/get/rV00-kyEgoZfZa8D9MOxeXqqFiItt0_l61WD_n9Rbm4?url=http%3A%2F%2Ffeeds.bbci.co.uk%2Fnews%2Frss.xml%3Fedition%3Duk HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 0
```

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 21:03:06 GMT
Content-Type: text/xml; charset=utf-8
Content-Length: 41148

<?xml version="1.0" encoding="UTF-8"?>
<?xml-stylesheet title="XSL_formatting" type="text/xsl" href="/shared/bsp/xsl/rss/nolsol.xsl"?>
<rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0" xmlns:media="http://search.yahoo.com/mrss/">
    <channel>
        <title><![CDATA[BBC News - Home]]></title>
        <description><![CDATA[BBC News - Home]]></description>
        <link>http://www.bbc.co.uk/news/</link>
        <image>
            <url>http://news.bbcimg.co.uk/nol/shared/img/bbc_news_120x60.gif</url>
            <title>BBC News - Home</title>
            <link>http://www.bbc.co.uk/news/</link>
        </image>

...truncated for brevity...

```

## Advanced

Here is a script which demonstrates inspecting the global request objects:

* `$uri` - The full URI of the request as a string.
* `$headers` - All HTTP Headers passed in the request.
* `$cookies` - All HTTP Cookies passed in the request.
* `$params` - All Query String parameters passed int the URI. 
Also if the request was a POST or PUT, and it's `Content-Type` was 
`application/x-www-form-urlencoded`, then the Form variables will be 
in here too.
* `$body` - This is the POST or PUT body of the request, as a string. 

It also demonstrates use of the `JSON` and `console` global objects. 

The script looks like this:

```javascript
var req  = {
  uri: $uri,
  headers: $headers,
  cookies: $cookies,
  params: $params,
  body: $body  
};

var reqString = JSON.stringify(req);

console.log(reqString);

reqString;
```

Lets deploy the script:

```
POST http://localhost:8181/deploy/inspect_request HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 192

var req  = {
  uri: $uri,
  headers: $headers,
  cookies: $cookies,
  params: $params,
  body: $body  
};

var reqString = JSON.stringify(req);

console.log(reqString);

reqString;
```

All is well, so we get back the hash version:

```
HTTP/1.1 200 OK
Date: Fri, 01 Jul 2016 21:49:31 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

HqTisjBvsBLPtpZ0E9t38Lr_LUSkBmMQOG5SDM26leU
```

Now we have succesfully deployed, we can execute the script. To make it
more interesting, lets expand the amount of information we pass in the 
request, and also use a HTTP POST rather than a GET:

```
POST http://localhost:8181/run/inspect_request/HqTisjBvsBLPtpZ0E9t38Lr_LUSkBmMQOG5SDM26leU?foo=bar&fan=ban HTTP/1.1
Host: localhost:8181
Connection: keep-alive
Cache-Control: max-age=0
Upgrade-Insecure-Requests: 1
User-Agent: Mozilla/5.0 (SAMSUNG_FRIDGEFREEZER_10.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8
DNT: 1
Referer: https://www.google.co.uk/
Accept-Encoding: gzip, deflate, sdch
Accept-Language: en-US,en;q=0.8
Cookie: __utmc=123456789; prov=412480918059358; __qca=2340975023957; _ga=238093583532523
If-Modified-Since: Fri, 01 Jul 2016 21:33:40 GMT
Content-Length: 38

This is some nice POST body data!!!111
```

And so, here is the JSON formatted response we get back:

```
HTTP/1.1 200 OK
Date: Fri, 01 Jul 2016 22:18:15 GMT
Content-Length: 915
Content-Type: text/plain; charset=utf-8

{"body":"This is some nice POST body data!!!111","cookies":{"__qca":"2340975023957","__utmc":"123456789","_ga":"238093583532523","prov":"412480918059358"},"headers":{"Accept":["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"],"Accept-Encoding":["gzip, deflate, sdch"],"Accept-Language":["en-US,en;q=0.8"],"Cache-Control":["max-age=0"],"Connection":["keep-alive"],"Content-Length":["38"],"Cookie":["__utmc=123456789; prov=412480918059358; __qca=2340975023957; _ga=238093583532523"],"Dnt":["1"],"If-Modified-Since":["Fri, 01 Jul 2016 21:33:40 GMT"],"Referer":["https://www.google.co.uk/"],"Upgrade-Insecure-Requests":["1"],"User-Agent":["Mozilla/5.0 (SAMSUNG_FRIDGEFREEZER_10.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"]},"params":{"fan":["ban"],"foo":["bar"]},"uri":"/run/inspect_request/HqTisjBvsBLPtpZ0E9t38Lr_LUSkBmMQOG5SDM26leU?foo=bar\u0026fan=ban"}
```

In a more readable format:

```json
{
	"body" : "This is some nice POST body data!!!111",
	"cookies" : {
		"__qca" : "2340975023957",
		"__utmc" : "123456789",
		"_ga" : "238093583532523",
		"prov" : "412480918059358"
	},
	"headers" : {
		"Accept" : ["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"],
		"Accept-Encoding" : ["gzip, deflate, sdch"],
		"Accept-Language" : ["en-US,en;q=0.8"],
		"Cache-Control" : ["max-age=0"],
		"Connection" : ["keep-alive"],
		"Content-Length" : ["38"],
		"Cookie" : ["__utmc=123456789; prov=412480918059358; __qca=2340975023957; _ga=238093583532523"],
		"Dnt" : ["1"],
		"If-Modified-Since" : ["Fri, 01 Jul 2016 21:33:40 GMT"],
		"Referer" : ["https://www.google.co.uk/"],
		"Upgrade-Insecure-Requests" : ["1"],
		"User-Agent" : ["Mozilla/5.0 (SAMSUNG_FRIDGEFREEZER_10.2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36"]
	},
	"params" : {
		"fan" : ["ban"],
		"foo" : ["bar"]
	},
	"uri" : "/run/inspect_request/HqTisjBvsBLPtpZ0E9t38Lr_LUSkBmMQOG5SDM26leU?foo=bar\u0026fan=ban"
}
```

## Coming Soon!

* Integrated Script Editor and Version Browser!

## Credits

* [otto Javascript Interpreter](https://github.com/robertkrimen/otto)
used to provide JavaScript execution framework.
* [esc](https://github.com/mjibson/esc) is a simple file embedder for Go.
* [Monaco Editor](https://github.com/Microsoft/monaco-editor) is a browser based code editor.