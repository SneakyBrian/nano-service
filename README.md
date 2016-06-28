# nano-service

nano-service is an example of building a platform that allows rapid
deployment of web application code.

The server-side script that runs is written in JavaScript, and 
executes in the context of the request, returning the evaluated result.

## Deploying code

As an example, let's say we wanted to create a friendly service that
responds with a cheerful greeting, we can write the script like this:

```
var result = "Hello World!";

result;
```

In order to deploy this code, we need to HTTP POST it to the 
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

Request parameters are passed in the query string. They are exposed to
the script as members of the `$query` object.

In this example we take two parameters passed to the script, `val1`
and `val2`, add them together and return the result.

The JavaScript looks like this:

```
Number($query.val1) + Number($query.val2);
```

So we deploy this code:

```
POST http://localhost:8181/deploy/add HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 42

Number($query.val1) + Number($query.val2);
```

And get back the hash of the validated code:

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 20:59:23 GMT
Content-Length: 43
Content-Type: text/plain; charset=utf-8

uH6BOHvnm0VUCBpa7TltWSUwYgoR1jLRctrZVgiBfpo
```

We then execute it, passing `12` as the value of `val1`, and `34` as
the value of `val2`:

```
GET http://localhost:8181/run/add/uH6BOHvnm0VUCBpa7TltWSUwYgoR1jLRctrZVgiBfpo?val1=12&val2=34 HTTP/1.1
User-Agent: Fiddler
Host: localhost:8181
Content-Length: 0
```

The response is the value `46`:

```
HTTP/1.1 200 OK
Date: Tue, 28 Jun 2016 21:00:21 GMT
Content-Length: 2
Content-Type: text/plain; charset=utf-8

46
```

## Simple HTTP

The nano-service scripting runtime has some simple synchronous http 
functions for `GET` and `POST` operations:

```
var result = $get(url);
```

```
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

$get($query.url)
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

## Credits

* [otto Javascript Interpreter](https://github.com/robertkrimen/otto)
used to provide JavaScript execution framework. 