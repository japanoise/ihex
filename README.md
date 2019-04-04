# ihex - Intel HEX pretty printer

Pretty prints simple [Intel HEX][1] format documents.

## Rationale

Recently I've become interested in the [Altair 8800][2], the legendary "first
home computer". One of the best resources on the 'net for this is the [website
of Mike Douglas' clone system][3] and the [corresponding YouTube channel.][4]

Much of the front panel software is available in Intel HEX format on Mike's
website, as well as original assembly and final binaries. However, when entering
the data into his computer, he uses some form of octal representation which is
not shown or available to download. So, because I'm bored, I wrote a pretty
printer for the Intel HEX format, so I can render the hex files to a similar
format.

Note that it only supports I8HEX format; I.E. it doesn't interpret linear or
segment addresses.

## Usage

The file `KILLBITS.HEX` is provided as an example.

For an octal listing:

``` shellsession
$ ihex -w -o -H KILLBITS.HEX | cut -f "2 4"
```
As above, but with useful colored backgrounds:

``` shellsession
$ ihex -w -c -o -H KILLBITS.HEX | cut -f "2 4"
```

To verify a file's checksums:

``` shellsession
$ ihex -w KILLBITS.HEX >/dev/null
```

## Installation

``` shellsession
$ go get github.com/japanoise/ihex
```

## Copying

Licensed MIT.

[1]: https://en.wikipedia.org/wiki/Intel_HEX
[2]: https://en.wikipedia.org/wiki/Altair_8800
[3]: http://altairclone.com/
[4]: https://www.youtube.com/user/deramp5113
