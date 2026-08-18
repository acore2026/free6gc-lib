use std::fs;
use std::path::PathBuf;

use ngap_codec::types::{
    AMF_UE_NGAP_ID, Cause, CauseRadioNetwork, Criticality, InitialContextSetupFailure,
    InitialContextSetupFailureProtocolIEs, InitialContextSetupFailureProtocolIEs_Entry,
    InitialContextSetupFailureProtocolIEs_EntryValue,
    InitialContextSetupRequestProtocolIEs_EntryValue,
    InitialContextSetupResponseProtocolIEs_EntryValue, InitialUEMessageProtocolIEs_EntryValue,
    InitiatingMessage, InitiatingMessageValue, NGAP_PDU, NGSetupRequestProtocolIEs_EntryValue,
    NGSetupResponseProtocolIEs_EntryValue, PDUSessionID, PDUSessionResourceSetupItemSUReq,
    PDUSessionResourceSetupItemSUReqPDUSessionResourceSetupRequestTransfer,
    PDUSessionResourceSetupListSUReq, PDUSessionResourceSetupRequest,
    PDUSessionResourceSetupRequestProtocolIEs, PDUSessionResourceSetupRequestProtocolIEs_Entry,
    PDUSessionResourceSetupRequestProtocolIEs_EntryValue, PDUSessionResourceSetupRequestTransfer,
    PDUSessionResourceSetupResponseProtocolIEs_EntryValue, PDUSessionResourceSetupResponseTransfer,
    PDUSessionResourceSetupUnsuccessfulTransfer, ProcedureCode, ProtocolIE_ID, RAN_UE_NGAP_ID,
    S_NSSAI, SD, SST, SuccessfulOutcomeValue, UnsuccessfulOutcome, UnsuccessfulOutcomeValue,
};
use ngap_codec::{AperMessage, CodecError, decode, encode};

const NG_SETUP_RESPONSE: &str = include_str!("../../testdata/interop/ng_setup_response.hex");
const NG_SETUP_REQUEST: &str = include_str!("../../testdata/interop/ng_setup_request.hex");
const INITIAL_UE_MESSAGE: &str = include_str!("../../testdata/interop/initial_ue_message.hex");
const PDU_SESSION_RESOURCE_SETUP_REQUEST: &str =
    include_str!("../../testdata/interop/pdu_session_resource_setup_request.hex");
const UNKNOWN_PROCEDURE: &str = "005500170000020052400601805445535400524006018054455354";
const UNKNOWN_IE: &str = "001500170000025555400601805445535400524006018054455354";
const OUT_OF_SCOPE_AMF_STATUS_INDICATION: &str =
    "000140150000010078000e002002f839cafe000180414d4631";

fn bytes(encoded: &str) -> Vec<u8> {
    hex::decode(encoded.trim()).expect("test vector must be valid hexadecimal")
}

fn go_fuzz_bytes(encoded: &str) -> Vec<u8> {
    let literal = encoded
        .strip_prefix("go test fuzz v1\n[]byte(\"")
        .and_then(|value| value.strip_suffix("\")\n"))
        .expect("corpus entry must contain one Go byte string");
    let source = literal.as_bytes();
    let mut decoded = Vec::with_capacity(source.len());
    let mut index = 0;

    while index < source.len() {
        if source[index] != b'\\' {
            decoded.push(source[index]);
            index += 1;
            continue;
        }

        match source.get(index + 1) {
            Some(b'x') => {
                let digits = std::str::from_utf8(
                    source
                        .get(index + 2..index + 4)
                        .expect("hex escape must contain two digits"),
                )
                .expect("hex escape must be ASCII");
                decoded.push(u8::from_str_radix(digits, 16).expect("hex escape must be valid"));
                index += 4;
            }
            Some(b'f') => {
                decoded.push(0x0c);
                index += 2;
            }
            Some(b'n') => {
                decoded.push(b'\n');
                index += 2;
            }
            escape => panic!("unsupported Go corpus escape: {escape:?}"),
        }
    }

    decoded
}

fn go_fuzz_corpus_entry(name: &str) -> Vec<u8> {
    let path = PathBuf::from(env!("CARGO_MANIFEST_DIR"))
        .join("../testdata/fuzz/FuzzNGAP")
        .join(name);
    let encoded = fs::read_to_string(path).expect("Go fuzz corpus entry must be readable");
    go_fuzz_bytes(&encoded)
}

#[test]
fn ng_setup_response_matches_go_vector() {
    let wire = bytes(NG_SETUP_RESPONSE);
    let pdu = decode(&wire).expect("NG Setup Response must decode");

    let NGAP_PDU::SuccessfulOutcome(outcome) = &pdu else {
        panic!("expected SuccessfulOutcome, got {pdu:?}");
    };
    let SuccessfulOutcomeValue::Id_NGSetup(response) = &outcome.value else {
        panic!("expected NG Setup Response, got {:?}", outcome.value);
    };

    assert_eq!(outcome.procedure_code.0, 21);
    assert_eq!(response.protocol_i_es.0.len(), 4);
    assert!(response.protocol_i_es.0.iter().any(|ie| matches!(
        &ie.value,
        NGSetupResponseProtocolIEs_EntryValue::Id_AMFName(name) if name.0 == "amf1"
    )));
    assert_eq!(encode(&pdu).expect("re-encode must succeed"), wire);
}

#[test]
fn ng_setup_request_matches_go_vector() {
    let wire = bytes(NG_SETUP_REQUEST);
    let pdu = decode(&wire).expect("NG Setup Request must decode");

    let NGAP_PDU::InitiatingMessage(message) = &pdu else {
        panic!("expected InitiatingMessage, got {pdu:?}");
    };
    let InitiatingMessageValue::Id_NGSetup(request) = &message.value else {
        panic!("expected NG Setup Request, got {:?}", message.value);
    };

    assert_eq!(message.procedure_code.0, 21);
    assert_eq!(request.protocol_i_es.0.len(), 4);
    assert!(request.protocol_i_es.0.iter().any(|ie| matches!(
        &ie.value,
        NGSetupRequestProtocolIEs_EntryValue::Id_RANNodeName(name)
            if name.0 == "UERANSIM-gnb-208-93-1"
    )));
    assert_eq!(encode(&pdu).expect("re-encode must succeed"), wire);
}

#[test]
fn initial_ue_message_matches_go_vector() {
    let wire = bytes(INITIAL_UE_MESSAGE);
    let pdu = decode(&wire).expect("Initial UE Message must decode");

    let NGAP_PDU::InitiatingMessage(message) = &pdu else {
        panic!("expected InitiatingMessage, got {pdu:?}");
    };
    let InitiatingMessageValue::Id_InitialUEMessage(initial) = &message.value else {
        panic!("expected Initial UE Message, got {:?}", message.value);
    };

    assert_eq!(message.procedure_code.0, 15);
    assert_eq!(initial.protocol_i_es.0.len(), 6);
    assert!(initial.protocol_i_es.0.iter().any(|ie| matches!(
        &ie.value,
        InitialUEMessageProtocolIEs_EntryValue::Id_SelectedPLMNIdentity(plmn)
            if plmn.0 == [0x64, 0xf6, 0x66]
    )));
    assert_eq!(encode(&pdu).expect("re-encode must succeed"), wire);
}

#[test]
fn unknown_and_out_of_scope_procedures_are_errors() {
    assert!(matches!(
        decode(&bytes(UNKNOWN_PROCEDURE)),
        Err(CodecError::Decode(_))
    ));
    assert!(matches!(
        decode(&bytes(UNKNOWN_IE)),
        Err(CodecError::Decode(_))
    ));
    assert!(matches!(
        decode(&bytes(OUT_OF_SCOPE_AMF_STATUS_INDICATION)),
        Err(CodecError::Decode(_))
    ));
}

#[test]
fn pdu_session_resource_setup_request_encodes_through_public_api() {
    let session = PDUSessionResourceSetupItemSUReq {
        pdu_session_id: PDUSessionID(10),
        pdu_session_nas_pdu: None,
        s_nssai: S_NSSAI {
            sst: SST(vec![1]),
            sd: Some(SD(vec![0x01, 0x02, 0x03])),
            ie_extensions: None,
        },
        pdu_session_resource_setup_request_transfer:
            PDUSessionResourceSetupItemSUReqPDUSessionResourceSetupRequestTransfer(vec![
                0x01, 0x02, 0x03, 0x04,
            ]),
        ie_extensions: None,
    };
    let pdu = NGAP_PDU::InitiatingMessage(InitiatingMessage {
        procedure_code: ProcedureCode(29),
        criticality: Criticality(Criticality::REJECT),
        value: InitiatingMessageValue::Id_PDUSessionResourceSetup(PDUSessionResourceSetupRequest {
            protocol_i_es: PDUSessionResourceSetupRequestProtocolIEs(vec![
                    PDUSessionResourceSetupRequestProtocolIEs_Entry {
                        id: ProtocolIE_ID(10),
                        criticality: Criticality(Criticality::REJECT),
                        value: PDUSessionResourceSetupRequestProtocolIEs_EntryValue::
                            Id_AMF_UE_NGAP_ID(AMF_UE_NGAP_ID(0x0001_0203_0405)),
                    },
                    PDUSessionResourceSetupRequestProtocolIEs_Entry {
                        id: ProtocolIE_ID(85),
                        criticality: Criticality(Criticality::REJECT),
                        value: PDUSessionResourceSetupRequestProtocolIEs_EntryValue::
                            Id_RAN_UE_NGAP_ID(RAN_UE_NGAP_ID(0x1122_3344)),
                    },
                    PDUSessionResourceSetupRequestProtocolIEs_Entry {
                        id: ProtocolIE_ID(74),
                        criticality: Criticality(Criticality::REJECT),
                        value: PDUSessionResourceSetupRequestProtocolIEs_EntryValue::
                            Id_PDUSessionResourceSetupListSUReq(
                                PDUSessionResourceSetupListSUReq(vec![session]),
                            ),
                    },
                ]),
        }),
    });

    let wire = bytes(PDU_SESSION_RESOURCE_SETUP_REQUEST);
    assert_eq!(
        encode(&pdu).expect("PDU Session Resource Setup Request must encode"),
        wire
    );
    assert_eq!(decode(&wire).expect("shared request must decode"), pdu);
}

#[test]
fn truncated_supported_messages_are_errors() {
    for wire in [
        bytes(NG_SETUP_RESPONSE),
        bytes(NG_SETUP_REQUEST),
        bytes(INITIAL_UE_MESSAGE),
    ] {
        for prefix_len in 0..wire.len() {
            assert!(
                decode(&wire[..prefix_len]).is_err(),
                "truncated prefix of length {prefix_len} unexpectedly decoded"
            );
        }
    }
}

#[test]
fn go_fuzz_corpus_never_panics() {
    let corpus = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../testdata/fuzz/FuzzNGAP");
    let entries = fs::read_dir(&corpus).expect("Go fuzz corpus must be present");

    for entry in entries {
        let path = entry.expect("corpus entry must be readable").path();
        let encoded = fs::read_to_string(&path).expect("corpus input must be readable");
        let input = go_fuzz_bytes(&encoded);
        let _ = decode(&input);
    }
}

#[test]
fn registration_and_session_establishment_messages_match_go_corpus() {
    for (name, expected_direction, expected_procedure) in [
        ("ci0003", "initiating", 4),
        ("ci0007", "initiating", 14),
        ("ci0008", "successful", 14),
        ("ci0004", "initiating", 46),
        ("ci0011", "initiating", 29),
        ("ci0012", "successful", 29),
        ("ci0036", "successful", 29),
        ("ci0158", "successful", 29),
        ("ci0854", "successful", 29),
    ] {
        let wire = go_fuzz_corpus_entry(name);
        let pdu = decode(&wire).unwrap_or_else(|error| {
            panic!("{name} must decode as procedure {expected_procedure}: {error}")
        });
        let (direction, procedure) = match &pdu {
            NGAP_PDU::InitiatingMessage(message) => ("initiating", message.procedure_code.0),
            NGAP_PDU::SuccessfulOutcome(message) => ("successful", message.procedure_code.0),
            NGAP_PDU::UnsuccessfulOutcome(message) => ("unsuccessful", message.procedure_code.0),
        };
        assert_eq!(
            (direction, procedure),
            (expected_direction, expected_procedure)
        );
        assert_eq!(
            encode(&pdu).unwrap_or_else(|error| panic!("{name} must re-encode: {error}")),
            wire,
            "{name} must match the Go APER bytes"
        );
    }
}

#[test]
fn session_establishment_transfer_containers_match_go_corpus() {
    let request_pdu = decode(&go_fuzz_corpus_entry("ci0011"))
        .expect("Go PDU Session Resource Setup Request must decode");
    let NGAP_PDU::InitiatingMessage(request_message) = request_pdu else {
        panic!("expected initiating message");
    };
    let InitiatingMessageValue::Id_PDUSessionResourceSetup(request) = request_message.value else {
        panic!("expected PDU Session Resource Setup Request");
    };
    let request_transfer_wire = request
        .protocol_i_es
        .0
        .iter()
        .find_map(|ie| {
            match &ie.value {
            PDUSessionResourceSetupRequestProtocolIEs_EntryValue::
                Id_PDUSessionResourceSetupListSUReq(list) => Some(
                    list.0[0]
                        .pdu_session_resource_setup_request_transfer
                        .0
                        .as_slice(),
                ),
            _ => None,
        }
        })
        .expect("request must contain a setup transfer");
    let request_transfer =
        PDUSessionResourceSetupRequestTransfer::decode_aper(request_transfer_wire)
            .expect("request transfer must decode");
    assert_eq!(
        request_transfer
            .encode_aper()
            .expect("request transfer must re-encode"),
        request_transfer_wire
    );

    let response_pdu = decode(&go_fuzz_corpus_entry("ci0012"))
        .expect("Go PDU Session Resource Setup Response must decode");
    let NGAP_PDU::SuccessfulOutcome(response_message) = response_pdu else {
        panic!("expected successful outcome");
    };
    let SuccessfulOutcomeValue::Id_PDUSessionResourceSetup(response) = response_message.value
    else {
        panic!("expected PDU Session Resource Setup Response");
    };
    let response_transfer_wire = response
        .protocol_i_es
        .0
        .iter()
        .find_map(|ie| {
            match &ie.value {
            PDUSessionResourceSetupResponseProtocolIEs_EntryValue::
                Id_PDUSessionResourceSetupListSURes(list) => Some(
                    list.0[0]
                        .pdu_session_resource_setup_response_transfer
                        .0
                        .as_slice(),
                ),
            _ => None,
        }
        })
        .expect("response must contain a setup transfer");
    let response_transfer =
        PDUSessionResourceSetupResponseTransfer::decode_aper(response_transfer_wire)
            .expect("response transfer must decode");
    assert_eq!(
        response_transfer
            .encode_aper()
            .expect("response transfer must re-encode"),
        response_transfer_wire
    );
}

#[test]
fn initial_context_session_transfers_in_go_corpus_decode() {
    let corpus = PathBuf::from(env!("CARGO_MANIFEST_DIR")).join("../testdata/fuzz/FuzzNGAP");
    let mut request_transfers = 0;
    let mut response_transfers = 0;

    for entry in fs::read_dir(corpus).expect("Go fuzz corpus must be present") {
        let path = entry.expect("corpus entry must be readable").path();
        let encoded = fs::read_to_string(path).expect("corpus input must be readable");
        let Ok(pdu) = decode(&go_fuzz_bytes(&encoded)) else {
            continue;
        };

        match pdu {
            NGAP_PDU::InitiatingMessage(message) => {
                let InitiatingMessageValue::Id_InitialContextSetup(request) = message.value else {
                    continue;
                };
                for ie in request.protocol_i_es.0 {
                    if let InitialContextSetupRequestProtocolIEs_EntryValue::
                        Id_PDUSessionResourceSetupListCxtReq(list) = ie.value
                    {
                        for item in list.0 {
                            let wire = item.pdu_session_resource_setup_request_transfer.0;
                            let transfer =
                                PDUSessionResourceSetupRequestTransfer::decode_aper(&wire)
                                    .expect("initial-context request transfer must decode");
                            assert_eq!(
                                transfer
                                    .encode_aper()
                                    .expect("initial-context request transfer must re-encode"),
                                wire
                            );
                            request_transfers += 1;
                        }
                    }
                }
            }
            NGAP_PDU::SuccessfulOutcome(message) => {
                let SuccessfulOutcomeValue::Id_InitialContextSetup(response) = message.value else {
                    continue;
                };
                for ie in response.protocol_i_es.0 {
                    if let InitialContextSetupResponseProtocolIEs_EntryValue::
                        Id_PDUSessionResourceSetupListCxtRes(list) = ie.value
                    {
                        for item in list.0 {
                            let wire = item.pdu_session_resource_setup_response_transfer.0;
                            let transfer =
                                PDUSessionResourceSetupResponseTransfer::decode_aper(&wire)
                                    .expect("initial-context response transfer must decode");
                            assert_eq!(
                                transfer
                                    .encode_aper()
                                    .expect("initial-context response transfer must re-encode"),
                                wire
                            );
                            response_transfers += 1;
                        }
                    }
                }
            }
            NGAP_PDU::UnsuccessfulOutcome(_) => {}
        }
    }

    assert!(
        request_transfers > 0,
        "no initial-context request transfer tested"
    );
    assert!(
        response_transfers > 0,
        "no initial-context response transfer tested"
    );
}

#[test]
fn registration_and_session_failure_types_use_public_api() {
    let cause = Cause::RadioNetwork(CauseRadioNetwork(
        CauseRadioNetwork::RADIO_RESOURCES_NOT_AVAILABLE,
    ));
    let failure = NGAP_PDU::UnsuccessfulOutcome(UnsuccessfulOutcome {
        procedure_code: ProcedureCode(14),
        criticality: Criticality(Criticality::REJECT),
        value: UnsuccessfulOutcomeValue::Id_InitialContextSetup(InitialContextSetupFailure {
            protocol_i_es: InitialContextSetupFailureProtocolIEs(vec![
                InitialContextSetupFailureProtocolIEs_Entry {
                    id: ProtocolIE_ID(10),
                    criticality: Criticality(Criticality::IGNORE),
                    value: InitialContextSetupFailureProtocolIEs_EntryValue::Id_AMF_UE_NGAP_ID(
                        AMF_UE_NGAP_ID(1),
                    ),
                },
                InitialContextSetupFailureProtocolIEs_Entry {
                    id: ProtocolIE_ID(85),
                    criticality: Criticality(Criticality::IGNORE),
                    value: InitialContextSetupFailureProtocolIEs_EntryValue::Id_RAN_UE_NGAP_ID(
                        RAN_UE_NGAP_ID(1),
                    ),
                },
                InitialContextSetupFailureProtocolIEs_Entry {
                    id: ProtocolIE_ID(15),
                    criticality: Criticality(Criticality::IGNORE),
                    value: InitialContextSetupFailureProtocolIEs_EntryValue::Id_Cause(
                        cause.clone(),
                    ),
                },
            ]),
        }),
    });
    let failure_wire = encode(&failure).expect("Initial Context Setup Failure must encode");
    assert_eq!(
        decode(&failure_wire).expect("Initial Context Setup Failure must decode"),
        failure
    );

    let transfer = PDUSessionResourceSetupUnsuccessfulTransfer {
        cause,
        criticality_diagnostics: None,
        ie_extensions: None,
    };
    let transfer_wire = transfer
        .encode_aper()
        .expect("unsuccessful setup transfer must encode");
    assert_eq!(
        PDUSessionResourceSetupUnsuccessfulTransfer::decode_aper(&transfer_wire)
            .expect("unsuccessful setup transfer must decode"),
        transfer
    );
}
