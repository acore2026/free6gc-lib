use std::env;
use std::error::Error;
use std::path::PathBuf;

use asn1_compiler::{
    Asn1Compiler,
    generator::{Codec, Derive, Visibility},
};

fn main() -> Result<(), Box<dyn Error>> {
    const SCHEMA: &str = "asn1/ngap-15.8.0-min.asn1";

    println!("cargo:rerun-if-changed={SCHEMA}");
    println!("cargo:rerun-if-changed=build.rs");

    let output = PathBuf::from(env::var_os("OUT_DIR").ok_or("OUT_DIR is not set")?)
        .join("ngap_generated.rs");
    let output = output
        .to_str()
        .ok_or("generated-code path is not valid UTF-8")?;

    let mut compiler = Asn1Compiler::new(
        output,
        &Visibility::Public,
        vec![Codec::Aper],
        vec![Derive::Clone, Derive::Debug, Derive::Eq, Derive::PartialEq],
    );
    compiler.compile_files(&[SCHEMA])?;

    Ok(())
}
