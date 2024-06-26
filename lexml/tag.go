package lexml

import (
    "io"
    "encoding/xml"
)

// Tags type represent fields of a XML file, those are generated on the go
type Tag struct {
	Id   int
	Name string
}


// findTags function find and return the exhaustive list of unique tags found in the xml file
func findTags(decoder *xml.Decoder) []Tag {

	tagList := make([]Tag, 0)
	index := 0

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch tk := tok.(type) {

		case xml.StartElement:

			if !tagExist(tk.Name.Local, tagList) {

				tagList = append(tagList, Tag{
					Id:   index,
					Name: tk.Name.Local,
				})
			}
		}
	}
	return tagList
}

// getTag function return a single tag from a tagList with a given name as string
func getTag(tagList []Tag, tok string) *Tag {

	for index, tag := range tagList {
		if tok == tag.Name {
			return &tagList[index]
		}
	}

	return nil
}

// tagExist function checks if a tag exist in the whole list of tags
func tagExist(word string, tagList []Tag) bool {

	for _, tag := range tagList {

		if tag.Name == word {
			return true
		}
	}

	return false
}
