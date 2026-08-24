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
	if err := validateJSONObject(proof, field, maximum); err != nil {
		return err
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(proof, &object); err != nil || len(object) != 2 {
		return newProtocolError(ErrorCodeInvalidJSON, field, -1, ErrInvalidJSON)
	}
	for _, name := range []string{"creator", "signature_value"} {
		encoded, ok := object[name]
		if !ok {
			return newProtocolError(ErrorCodeInvalidJSON, field+"."+name, -1, ErrInvalidJSON)
		}
		var value string
		if err := json.Unmarshal(encoded, &value); err != nil || value == "" {
			return newProtocolError(ErrorCodeInvalidJSON, field+"."+name, -1, ErrInvalidJSON)
		}
	}
	return nil
}
