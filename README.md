Embeds resources with data URLs to create single file web pages.

# Install

Make sure that your **go environment** is configured correctly and that `$GOPATH/bin` is added to `$PATH`.

```
cd $GOPATH/src/github.com/yungtravla/epoxy
go install
```

# Usage

First grab the source of a web page and save it locally.

```
curl https://www.google.com/ > google-test.html
```

Now you can use epoxy to fetch every resource in the web page and embed them into the source file.

```
epoxy -source google-test.html -origin https://www.google.com/ -no-html
```

![screenshot from 2018-10-23 14-57-34](https://user-images.githubusercontent.com/29265684/47336874-c74e2c00-d6d4-11e8-8406-ef5ea430238d.png)
