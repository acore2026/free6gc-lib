package acn

import "encoding/json"

func (c *Codec) encodeProof(output *encoder, proof json.RawMessage, field string) error {
	if err := validateJSONObject(proof, field, c.limits.MaxJSONContainer); err != nil {
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
	if err := validateJSONObject(proof, field, c.limits.MaxJSONContainer); err != nil {
		return nil, err
	}
	return cloneBytes(proof), nil
}
