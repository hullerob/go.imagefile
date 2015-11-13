Go.Image File
=============

[![Build Status](https://travis-ci.org/hullerob/go.imagefile.svg?branch=master)](https://travis-ci.org/hullerob/go.imagefile)

About
-----

This is Go implementation of [`farbfeld` image format](http://git.2f30.org/farbfeld/).

`imagefile` format was deprecated by `farbfeld`; documentation and refence
implementation is no longer available.

It uses Go's `image` interface, similar to `image/png`.


Install
-------

    go get github.com/hullerob/go.imagefile

Usage
-----

`Encode` and `Decode` for old `imagefile`.

`FFEncode` and `FFDecode` for new `farbfeld`.

See `examples`.
