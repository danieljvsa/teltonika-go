// Package teltonika provides a stable, public API for decoding and encoding
// Teltonika GPS device protocol frames (login and AVL data packets).
//
// The package exposes codec agnostic types such as Packet, AVLRecord and
// Command, so callers can work with decoded data without importing any
// internal implementation packages.
//
// Typical usage:
//
//	// Decode a complete frame (login or data, TCP or UDP).
//	packet, err := teltonika.Decode(rawFrame)
//
//	// Encode a packet back to a complete, wire-ready frame.
//	frame, err := teltonika.Encode(packet)
//
// Supported codecs: 08, 8E, 12, 13, 14, 15 and 16.
package teltonika
