# Developer manual

This manual is for developers who use and change teltonika-go, a Go
library that decodes and encodes binary data from Teltonika GPS
devices.

The library reads and writes AVL data frames and IMEI login frames
over TCP and UDP. It supports Codecs 8, 8E, 12, 13, 14, 15, and 16.

Use the public package in new code. The public package is
`github.com/danieljvsa/teltonika-go/public`. It exposes the full API
without internal implementation types.

## How this manual is organized

| Document | What it covers |
| --- | --- |
| getting-started.md | Add the library, decode a frame, encode a frame |
| datamodel.md | Wire frames, codecs, and the public data types |
| using-the-api.md | Decode and encode with options, command codecs, legacy API |
| testing.md | Build, test, vet, format, and fuzz the library |

## Workflow

Read the sections in this order:

1. Read [getting-started.md](getting-started.md) and run the two examples.
2. Read [datamodel.md](datamodel.md) to learn the data types.
3. Read [using-the-api.md](using-the-api.md) before you design your integration.
4. Read [testing.md](testing.md) before you change the code.

## How this manual is written

This manual follows Simplified Technical English. It uses short
sentences and a fixed vocabulary.

One name means one thing. The library decodes a frame and encodes a
packet. A packet is the decoded form of a frame. A record is one AVL
data point. An element is one I/O value.

Code identifiers, file names, and commands keep their exact spelling.
