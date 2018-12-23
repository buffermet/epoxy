Asynchronous and recursive data URL generator/embedder for single files or web pages.

# Install

Make sure that your **go environment** is configured correctly and that `$GOPATH/bin` is added to `$PATH`.

```
go get github.com/yungtravla/epoxy
cd $GOPATH/src/github.com/yungtravla/epoxy
go install
```

# Usage

First grab the source of a web page and save it locally.

```
curl https://twitter.com/ > twitter-index.html
```

Now you can use epoxy to fetch every resource in the web page and embed them into the source file.

```
epoxy -source twitter-index.html -origin https://twitter.com/ -recurse 3 -no-html
```

You can set the recursion limit with `-recurse` to choose how many nested resources should be embedded as data URLs for every resource.

![screenshot from 2018-12-23 18-58-11](https://user-images.githubusercontent.com/29265684/50382162-ed984400-06e4-11e9-813d-b0a4c8b64a16.png)

If you want to turn a single file into a data URL, set the recursion to 0 and epoxy will generate a data URL for the specified file.

```
epoxy -source twitter-index.html -recurse 0
```

# Options

```
  -print          print payload to stdout.

  -source PATH    path to source file.
  -origin URL     full URL to source file.

  -recurse INT    limit of recursions for resource embedding (default=1).
  -cores INT      limit of procs for async parsing (default=4).

  -no-unknown     don't embed unknown filetypes.
  -no-svg         don't embed svg files.
  -no-jpg         don't embed jpg files.
  -no-png         don't embed png files.
  -no-gif         don't embed gif files.
  -no-webp        don't embed webp files.
  -no-cr2         don't embed cr2 files.
  -no-tif         don't embed tif files.
  -no-bmp         don't embed bmp files.
  -no-jxr         don't embed jxr files.
  -no-psd         don't embed psd files.
  -no-ico         don't embed ico files.
  -no-mp4         don't embed mp4 files.
  -no-m4v         don't embed m4v files.
  -no-mkv         don't embed mkv files.
  -no-webm        don't embed webm files.
  -no-mov         don't embed mov files.
  -no-avi         don't embed avi files.
  -no-wmv         don't embed wmv files.
  -no-mpg         don't embed mpg files.
  -no-flv         don't embed flv files.
  -no-mid         don't embed mid files.
  -no-mp3         don't embed mp3 files.
  -no-m4a         don't embed m4a files.
  -no-ogg         don't embed ogg files.
  -no-flac        don't embed flac files.
  -no-wav         don't embed wav files.
  -no-amr         don't embed amr files.
  -no-epub        don't embed epub files.
  -no-zip         don't embed zip files.
  -no-tar         don't embed tar files.
  -no-rar         don't embed rar files.
  -no-gz          don't embed gz files.
  -no-bz2         don't embed bz2 files.
  -no-7z          don't embed 7z files.
  -no-xz          don't embed xz files.
  -no-pdf         don't embed pdf files.
  -no-exe         don't embed exe files.
  -no-swf         don't embed swf files.
  -no-rtf         don't embed rtf files.
  -no-eot         don't embed eot files.
  -no-ps          don't embed ps files.
  -no-sqlite      don't embed sqlite files.
  -no-nes         don't embed nes files.
  -no-crx         don't embed crx files.
  -no-cab         don't embed cab files.
  -no-deb         don't embed deb files.
  -no-ar          don't embed ar files.
  -no-z           don't embed z files.
  -no-lz          don't embed lz files.
  -no-rpm         don't embed rpm files.
  -no-elf         don't embed elf files.
  -no-doc         don't embed doc files.
  -no-docx        don't embed docx files.
  -no-xls         don't embed xls files.
  -no-xlsx        don't embed xlsx files.
  -no-ppt         don't embed ppt files.
  -no-pptx        don't embed pptx files.
  -no-woff        don't embed woff files.
  -no-woff2       don't embed woff2 files.
  -no-ttf         don't embed ttf files.
  -no-otf         don't embed otf files.
  -no-css         don't embed css files.
  -no-html        don't embed html files.
  -no-js          don't embed js files.
  -no-json        don't embed json files.
```
