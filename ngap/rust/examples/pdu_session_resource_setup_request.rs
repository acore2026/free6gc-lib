use ngap_codec::types::{
    AMF_UE_NGAP_ID, Criticality, InitiatingMessage, InitiatingMessageValue, NGAP_PDU, PDUSessionID,
    PDUSessionResourceSetupItemSUReq,
    PDUSessionResourceSetupItemSUReqPDUSessionResourceSetupRequestTransfer,
    PDUSessionResourceSetupListSUReq, PDUSessionResourceSetupRequest,
    PDUSessionResourceSetupRequestProtocolIEs, PDUSessionResourceSetupRequestProtocolIEs_Entry,
    PDUSessionResourceSetupRequestProtocolIEs_EntryValue, ProcedureCode, ProtocolIE_ID,
    RAN_UE_NGAP_ID, S_NSSAI, SD, SST,
};
use ngap_codec::{Result, encode};

fn main() -> Result<()> {
    let session = PDUSessionResourceSetupItemSUReq {
        pdu_session_id: PDUSessionID(10),
        pdu_session_nas_pdu: None,
        s_nssai: S_NSSAI {
            sst: SST(vec![1]),
            sd: Some(SD(vec![0x01, 0x02, 0x03])),
            ie_extensions: None,
        },
        // Replace this example value with an APER-encoded
        // PDUSessionResourceSetupRequestTransfer from the SMF.
        pdu_session_resource_setup_request_transfer:
            PDUSessionResourceSetupItemSUReqPDUSessionResourceSetupRequestTransfer(vec![
                0x01, 0x02, 0x03, 0x04,
            ]),
        ie_extensions: None,
    };

    let request = PDUSessionResourceSetupRequest {
        protocol_i_es: PDUSessionResourceSetupRequestProtocolIEs(vec![
            PDUSessionResourceSetupRequestProtocolIEs_Entry {
                id: ProtocolIE_ID(10),
                criticality: Criticality(Criticality::REJECT),
                value:
                    PDUSessionResourceSetupRequestProtocolIEs_EntryValue::Id_AMF_UE_NGAP_ID(
                        AMF_UE_NGAP_ID(0x0001_0203_0405),
                    ),
            },
            PDUSessionResourceSetupRequestProtocolIEs_Entry {
                id: ProtocolIE_ID(85),
                criticality: Criticality(Criticality::REJECT),
                value:
                    PDUSessionResourceSetupRequestProtocolIEs_EntryValue::Id_RAN_UE_NGAP_ID(
                        RAN_UE_NGAP_ID(0x1122_3344),
                    ),
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
    };

    let pdu = NGAP_PDU::InitiatingMessage(InitiatingMessage {
        procedure_code: ProcedureCode(29),
        criticality: Criticality(Criticality::REJECT),
        value: InitiatingMessageValue::Id_PDUSessionResourceSetup(request),
    });

    let wire = encode(&pdu)?;
    println!("{wire:02x?}");
    Ok(())
}
