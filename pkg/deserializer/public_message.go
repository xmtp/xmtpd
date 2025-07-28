package deserializer

import "bytes"

type PublicMessageIn struct {
	Content FramedContentIn
	// and more
}

func (p *PublicMessageIn) TLSDeserialize(r *bytes.Reader) error {
	var content FramedContentIn
	if err := content.TLSDeserialize(r); err != nil {
		return err
	}

	p.Content = content
	return nil
}
