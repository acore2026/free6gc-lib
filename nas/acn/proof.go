package acn

import "encoding/json"

func (c *Codec) encodeProof(output *encoder, proof json.RawMessage, field string) error {
	if err := validateProof(proof, field, c.limits.MaxJSONContainer); err != nil {
		return err
	}
	output.lve(proof)
	return nil
}

func (c *Codec) decodeProof(input *decoder, field string) (json.RawMessage, error) {
	proof, err := input.lveBytes(field, c.limits.MaxJSONContainer)
	if err != nil {
		return nil, err
	}
	if err := validateProof(proof, field, c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return cloneBytes(proof), nil
}

func validateProof(proof json.RawMessage, field string, maximum int) error {
	// The codec owns binary framing and JSON-container integrity. Proof-suite
	// fields and signatures are validated by the destination that owns the
	// trusted verification key; enforcing one suite here would reject valid ACN
	// JsonWebSignature2020 messages before IDM or ACF can authenticate them.
	return validateJSONObject(proof, field, maximum)
}
