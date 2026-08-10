# Data model

This section explains how frames are structured on the wire and how the
library represents them in memory.

## Frame types

A device sends two kinds of frames.

- A login frame carries the device IMEI. It opens the session.
- A data frame carries AVL records or command responses.

`Decode` returns a packet with the field `Kind` set to `KindLogin` or
`KindData`.

## TCP data frames

A TCP data frame has this layout:

```
header (8 bytes) + codec id (1 byte) + payload + CRC (4 bytes)
```

The header starts with four zero bytes. The next four bytes hold the
Data Field Length.

The Data Field Length covers the codec id, the records, and the
trailing record count. It excludes the eight-byte header and the
four-byte CRC.

TCP does not preserve message boundaries. One socket read can hold part
of a frame, one full frame, or several frames. `Decode` requires a
complete frame. See [using-the-api.md](using-the-api.md) for the
receiving steps.

## UDP data frames

A UDP data frame has this layout:

```
length (2) + packet id (2) + version (1) + AVL packet id (1)
+ IMEI length (2) + IMEI + codec id (1) + payload
```

The length field covers everything after it. A complete UDP frame has
no CRC.

The decoder rejects a frame unless the declared length equals the
delivered length. Use the option `WithLenientUDPLength()`. The decoder
then accepts a declared length larger than the delivered frame.

## Login frames

A login frame has this layout:

```
length (2) + IMEI (ASCII)
```

The length field holds the byte count of the IMEI. `Decode` treats a
login frame as a login packet. The packet carries the IMEI as a string.

## Codecs

The library supports these codecs.

| Codec | Byte value | Carries |
| --- | --- | --- |
| 8 | 0x08 | AVL records |
| 8E | 0x8E | AVL records with extended I/O |
| 16 | 0x10 | AVL records with generation type |
| 12 | 0x0C | Command responses |
| 13 | 0x0D | Command responses with timestamps |
| 14 | 0x0E | Command responses with IMEI |
| 15 | 0x0F | Command responses with timestamps and IMEI |

The constants `Codec8`, `Codec8Ext`, `Codec12`, `Codec13`, `Codec14`,
`Codec15`, and `Codec16` name the codecs in code.

## Packet

A `Packet` is the top-level value. `Decode` returns it. `Encode`
accepts it.

```
Kind      login or data
Protocol  TCP or UDP
Header    transport header
IMEI      set for a login packet
Codec     codec identifier
Records   AVL records
Commands  command responses
```

## AVLRecord

An `AVLRecord` is one AVL data point.

```
Timestamp        time of the measurement
Priority         record priority
GPS              position and speed
EventIO          event I/O identifier
IOElements       I/O values
GenerationType   Codec 16 only
```

## GPSData

`GPSData` holds the GPS portion of a record.

- Latitude is in degrees. A positive value means north.
- Longitude is in degrees. A positive value means east.
- Altitude is in meters.
- Angle is in degrees from north, clockwise.
- Satellites is the count of visible satellites.
- Speed is a raw value from the device.

## IOElement

An `IOElement` has an ID and a value. The value is the raw byte content
in lowercase hexadecimal. The library keeps the value as hex so the
bytes round-trip without change.

## Command and CommandResponse

Command codecs carry responses.

A `Command` groups responses of one type. `Command.Type` is the string
`"Command"` or `"Response"`. A Codec 15 frame may carry a raw response
type byte. The decoder stores that byte as its decimal string, for
example `"11"`. The encoder writes the string back to the byte.

A `CommandResponse` holds one response. It has an optional timestamp,
the response text, the response as hex, the command type, and an IMEI.

## Transport header

A data packet carries a `Header` value.

- For TCP the header holds the `DataLength`.
- For UDP the header holds the length, the packet ID, the AVL packet
  ID, the IMEI length, and the IMEI.

The encoder writes the header for you. You set the fields for UDP
only.
