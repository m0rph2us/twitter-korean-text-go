# Korean Language Processor

Ported project of https://github.com/twitter/twitter-korean-text with go language.

Use it, PR, bug reports are all welcome.

# Requirements

You will need PCRE library/header for regular expression.(Due to golang's too many unsupported feature for regular expression.)

# How to use

Once you have got this module. You need to do some go-get for another dependent module.

```
go get github.com/acidd15/go-scala-util
go get github.com/glenn-brown/golang-pkg-pcre
```

And then, you should put some environment variables like this.

```
# To make current directory to GOPATH
export GOPATH=$PWD
export GOPATH=$GOPATH:$PWD/src/gitlab.com/acidd15/twitter-korean-text-go

# For pcre library header
export CGO_CFLAGS=-I/usr/local/include

# Dictionary resource
export KRGO_DIC_RSRC=$PWD/src/gitlab.com/acidd15/twitter-korean-text-go/src/tktg/resources/
```

Once you have completed preceding instructions then see the test files for your information.

