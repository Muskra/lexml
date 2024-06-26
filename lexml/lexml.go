package lexml

import (
	"encoding/xml"
	"fmt"
	"io"
)

// genData function recreates a recursive Data datastructure representation of the file itself. It's taking every tags, subtags and data to make them one recursive datastructure of Type Data
func genData(decoder *xml.Decoder, tagList []Tag, offset int) (Data, error) {

	data := NewData(offset)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch tk := tok.(type) {

		case xml.StartElement:
			name := tk.Name.Local

			data.Inners = append(data.Inners, NewData(0))

			data.Inners[len(data.Inners) - 1], err = genData(decoder, tagList, len(data.Inners) - 1)
			if err != nil {
				return Data{}, fmt.Errorf("recurse() -> %s", err)
			}

            data.Inners[len(data.Inners) - 1].Type = getTag(tagList, name)

		case xml.EndElement:
			return data, nil

		case xml.CharData:
            if string(tok.(xml.CharData)) != "\n" {
			    data.Value = string(tk)
            }

		case xml.ProcInst:
			continue

		default:
			return Data{}, fmt.Errorf("recurse() -> Unknown or unused Type encountered, got %T", tok)
		}
	}

	return data, nil
}

