# ⛨ Buckler ⛨

[![Build Status](https://travis-ci.org/gittip/img.shields.io.png)](https://travis-ci.org/gittip/img.shields.io)
[![Buckler Shield](http://b.repl.ca/v1/use-buckler-blue.png)](http://buckler.repl.ca)
[![Get Hype](http://b.repl.ca/v1/GET-HYPE!-orange.png)](http://buckler.repl.ca)
[![MIT License](http://b.repl.ca/v1/License-MIT-red.png)](LICENSE)
[![CLI interface](http://b.repl.ca/v1/command-line-blue.png)](#command-line)

Buckler is [Shields](https://github.com/olivierlacan/shields) as a Service (ShaaS, or alternatively, Badges as a Service)
for use in GitHub READMEs, or anywhere else. Use buckler with your favorite continuous integration tool, performance
monitoring service API, or ridiculous in-joke to surface information.

Buckler is available hosted at [b.repl.ca](http://buckler.repl.ca). You may use the [API](#API) to generate shields at runtime,
pregenerate them and host them on your own service, or run your own copy of Buckler to protect important company secrets.

# API

Buckler tries to make creating shields easy. Each shield request is a url that has three parts:
- `subject`
- `status`
- `colour`

Parts are separated by a hyphen. The request is suffixed by `.png` and prefixed with the Buckler host and API version, likely
`b.repl.ca/v1/`. Requests will take the form: `http://b.repl.ca/v1/$SUBJECT-$STATUS-$COLOR.png`

## Examples

- http://b.repl.ca/v1/build-passing-brightgreen.png ⇨ ![](http://b.repl.ca/v1/build-passing-brightgreen.png)
- http://b.repl.ca/v1/downloads-3.4K-blue.png ⇨ ![](http://b.repl.ca/v1/downloads-3.4K-blue.png)
- http://b.repl.ca/v1/coverage-unknown-lightgrey.png ⇨ ![](http://b.repl.ca/v1/coverage-unknown-lightgrey.png)
- http://b.repl.ca/v1/review-NACKED-red.png ⇨ ![](http://b.repl.ca/v1/review-NACKED-red.png)
- http://b.repl.ca/v1/enterprise-ready-ff69b4.png ⇨ ![](http://b.repl.ca/v1/enterprise-ready-ff69b4.png)

## Valid Colours

- `brightgreen` ⇨ ![](http://b.repl.ca/v1/colour-brightgreen-brightgreen.png)
- `green` ⇨ ![](http://b.repl.ca/v1/colour-green-green.png)
- `yellowgreen` ⇨ ![](http://b.repl.ca/v1/colour-yellowgreen-yellowgreen.png)
- `yellow` ⇨ ![](http://b.repl.ca/v1/colour-yellow-yellow.png)
- `orange` ⇨ ![](http://b.repl.ca/v1/colour-orange-orange.png)
- `red` ⇨ ![](http://b.repl.ca/v1/colour-red-red.png)
- `grey` ⇨ ![](http://b.repl.ca/v1/colour-grey-grey.png)
- `lightgrey` ⇨ ![](http://b.repl.ca/v1/colour-lightgrey-lightgrey.png)
- `blue` ⇨ ![](http://b.repl.ca/v1/colour-blue-blue.png)

Six digit RGB hexidecimal colour values work as well:

- `804000` - ![](http://b.repl.ca/v1/colour-brown-804000.png)

### Grey?

Don't worry; `gray` and `lightgray` work too.

## Escaping Underscores and Hyphens

Hyphens (`-`) are used to delimit individual fields in your shield request. To include a literal hyphen, use two hyphens (`--`):

http://b.repl.ca/v1/really--cool-status-yellow.png ⇨ ![](http://b.repl.ca/v1/really--cool-status-yellow.png)

Similarly, underscores (`_`) are used to indicated spaces. To include a literal underscore, use two underscores (`__`):

http://b.repl.ca/v1/__private-method_name-lightgrey.png ⇨ ![](http://b.repl.ca/v1/__private-method_name-lightgrey.png)

## URL Safe

Buckler API requests are just HTTP GETs, so remember to URL encode!

http://b.repl.ca/v1/uptime-99.99%25-yellowgreen.png ⇨ ![](http://b.repl.ca/v1/uptime-99.99%25-yellowgreen.png)

# Try It Out

Play around with the simple form on [b.repl.ca](http://b.repl.ca)

# Installing

```bash
go get github.com/gittip/img.shields.io
```

Alternatively, `git clone` and `go build` to run from source.

# Command Line

Buckler also provides a command line interface:

```bash
# writes to build-passing-brightgreen.png
buckler -v build -s passing -c brightgreen

# writes to my-custom-filename.png
buckler -v build -s passing -c green my-custom-filename.png

# writes to standard out
buckler -v license -s MIT -c blue -

# writes 2 shields
buckler build-passing-brightgreen.png license-MIT-blue.png
```

# Thanks

- Olivier Lacan for the [shields](https://github.com/olivierlacan/shields) repo
- Steve Matteson for [Open Sans](http://opensans.com/)
