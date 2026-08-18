// Package acn implements the private ACN control-PDU format documented by the
// Free6GC ACN prototype specification.
//
// ACN PDUs are carried in the payload container of an existing UL NAS
// TRANSPORT or DL NAS TRANSPORT message. They are not 3GPP 5GMM message types,
// and payload-container type 0x0e is a prototype-local assignment.
//
// This package encodes and decodes plain NAS messages. Callers are responsible
// for applying or removing NAS security, tracking transactions, verifying
// signatures and timestamps, and binding payload identities to authenticated UE
// context.
package acn
