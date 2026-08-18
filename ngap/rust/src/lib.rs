#![doc = include_str!("../README.md")]

use std::error::Error;
use std::fmt;

use asn1_codecs::aper::AperCodec as GeneratedAperCodec;
use asn1_codecs::{PerCodecData, PerCodecError};

/// SCTP payload protocol identifier used by NGAP, matching the Go package.
pub const PPID: u32 = 0x3c00_0000;

/// ASN.1 types generated from the repository's NGAP compatibility schema.
#[allow(
    clippy::all,
    non_camel_case_types,
    non_snake_case,
    non_upper_case_globals
)]
pub mod types {
    include!(concat!(env!("OUT_DIR"), "/ngap_generated.rs"));
}

pub use types::NGAP_PDU as NgapPdu;

/// Result type returned by this crate.
pub type Result<T> = std::result::Result<T, CodecError>;

/// Errors produced while encoding or decoding APER.
#[derive(Debug)]
pub enum CodecError {
    Encode(PerCodecError),
    Decode(PerCodecError),
}

impl fmt::Display for CodecError {
    fn fmt(&self, formatter: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Self::Encode(error) => write!(formatter, "NGAP APER encode failed: {error}"),
            Self::Decode(error) => write!(formatter, "NGAP APER decode failed: {error}"),
        }
    }
}

impl Error for CodecError {
    fn source(&self) -> Option<&(dyn Error + 'static)> {
        match self {
            Self::Encode(error) | Self::Decode(error) => Some(error),
        }
    }
}

/// Convenience APER operations for supported NGAP roots and transfer types.
pub trait AperMessage: Sized {
    fn encode_aper(&self) -> Result<Vec<u8>>;
    fn decode_aper(bytes: &[u8]) -> Result<Self>;
}

impl<T> AperMessage for T
where
    T: GeneratedAperCodec<Output = T>,
{
    fn encode_aper(&self) -> Result<Vec<u8>> {
        encode_generated(self)
    }

    fn decode_aper(bytes: &[u8]) -> Result<Self> {
        decode_generated(bytes)
    }
}

/// Encode a top-level NGAP PDU to APER bytes.
pub fn encode(pdu: &NgapPdu) -> Result<Vec<u8>> {
    pdu.encode_aper()
}

/// Decode APER bytes into a top-level NGAP PDU.
pub fn decode(bytes: &[u8]) -> Result<NgapPdu> {
    NgapPdu::decode_aper(bytes)
}

fn encode_generated<T>(value: &T) -> Result<Vec<u8>>
where
    T: GeneratedAperCodec,
{
    let mut codec = PerCodecData::new_aper();
    value.aper_encode(&mut codec).map_err(CodecError::Encode)?;
    Ok(codec.into_bytes())
}

fn decode_generated<T>(bytes: &[u8]) -> Result<T>
where
    T: GeneratedAperCodec<Output = T>,
{
    let mut codec = PerCodecData::from_slice_aper(bytes);
    T::aper_decode(&mut codec).map_err(CodecError::Decode)
}
