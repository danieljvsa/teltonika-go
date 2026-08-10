# Golden fixture provenance

Every file in this directory is a hexadecimal Teltonika AVL frame used as
decode/encode-conformance input. This document records where each frame came
from, whether it is complete, and any transformations applied.

All AVL-data examples originate from the official Teltonika product
documentation (wiki.teltonika-gps.com, "Codec" and "Command" pages). The
hex was copied verbatim; no bytes were modified. Files marked *truncated*
are genuine official examples that the Teltonika page itself shows
abridged, keeping their original declared-length header.

## AVL data frames

| File | Codec | Transport | Complete | Declared vs. delivered length |
| --- | --- | --- | --- | --- |
| codec8-tcp.hex | 8 | TCP | Yes | Data Length = frame − 12 (CRC excluded) |
| codec8e-tcp.hex | 8 Extended | TCP | Yes | Data Length = frame − 12 |
| codec16-tcp.hex | 16 | TCP | Yes | Data Length = frame − 12 |
| codec8-udp.hex | 8 | UDP | Yes | Length = frame − 2 |
| codec8e-udp.hex | 8 Extended | UDP | Yes | Length = frame − 2 |
| codec16-udp.hex | 16 | UDP | **Truncated** | Length 0x015B (=347) exceeds frame − 2 |

- **codec8-tcp.hex** — Teltonika wiki Codec 8 receiving example, 1 record.
  IMEI 351002040313538 encrypted in a TCP frame (PH/00001 FH/xxxx).
- **codec8e-tcp.hex** — Teltonika wiki Codec 8 Extended example, 1 record
  with a 2-byte I/O group.
- **codec16-tcp.hex** — Teltonika wiki Codec 16 example, 2 records.
- **codec8-udp.hex** — Teltonika wiki UDP receiving example, 1 record.
  IMEI 352093086403655. Declared length 0x003D equals frame − 2.
- **codec8e-udp.hex** — Teltonika wiki UDP receiving example, 1 record.
  Declared length 0x005F equals frame − 2.
- **codec16-udp.hex** — Teltonika wiki "receiving Codec 16 over UDP"
  example, showing the first record of a multi-record capture. The page
  presents the packet truncated (an ellipsis replaces the remainder), yet
  keeps the original header whose declared length 0x015B describes the full
  349-byte packet. Decoding therefore requires `WithLenientUDPLength()` and
  the frame cannot be re-encoded byte-for-byte (the encoder writes a
  truthful recomputed length). Device IMEI 352094085231592.

## Command frames

| File | Codec | Transport | Complete | Notes |
| --- | --- | --- | --- | --- |
| codec12-command.hex | 12 (Command) | TCP | Yes | Response type 6; device info message |
| codec13-response.hex | 13 (Command response) | TCP | Yes | 8-byte timestamp; "getinfo" |
| codec14-response.hex | 14 (Command response) | TCP | Yes | 8-byte IMEI prefix; firmware info message |
| codec15-response.hex | 15 (Command response) | TCP | Yes | Raw response type byte 0x0B |

- **codec12-command.hex** — Teltonika Command documentation example, 1
  response, message text beginning `IMEI:2019/7/22 7:22 RTC:...`.
- **codec13-response.hex** — 1 response, "getinfo", millisecond timestamp.
  Data Length 0x17.
- **codec14-response.hex** — 1 response, message carries firmware/device
  info such as `Ver:03.18.14_04 GPS:AXN_5.10_3333 Hw:FMB120`.
- **codec15-response.hex** — 1 response, second-precision timestamp, IMEI
  `0123456789123456`, message `Hello!\n`. The 0x0B response-type byte is a
  firmware convention preserved by the decoder as decimal `"11"` and again
  by the encoder, so the frame round-trips exactly.

## Verification status

- All complete frames decode and re-encode byte-for-byte
  (`TestRoundTripOfficialFixtures`).
- `codec16-udp.hex` decodes only (`TestDecodeOfficialCodec16UDPFixture`),
  and only under the lenient UDP option.
- The device model behind each capture is not recorded by the wiki source
  pages; where the documentation names a series it is noted above.