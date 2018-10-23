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

Now you can use epoxy to fetch every resource of the web page and embed them into the source file.

```
epoxy -source google-test.html -origin https://www.google.com/ -no-html
```
