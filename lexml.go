package lexml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// Set is the entrypoint struct to interact with this program with Fields as a list of Tag types, Content that represent the whole XML file, Raw wich represent the raw byte data of the file
type Set struct {
	Fields  []Tag
	Content Data
	Raw     []byte
}

// Data is a recursive representation of a parsed XML file content
type Data struct {
	Type   *Tag
	Index  int
	Value  string
	Inners []Data
}

// Tags represent fields of a XML file, those are generated on the go
type Tag struct {
	Id   int
	Name string
}

// Parse convert the whole file into a DataSet datastructure
func (set Set) Parse() ([]Tag, Data, error) {

	reader := bytes.NewReader(set.Raw)
	decoder := xml.NewDecoder(reader)

	tagList := findTags(decoder)

	reader = bytes.NewReader(set.Raw)
	decoder = xml.NewDecoder(reader)

	data, err := recurse(decoder, tagList)
	if err != nil {
		return []Tag{}, Data{}, fmt.Errorf("Parse() -> %s", err)
	}

	return tagList, data, nil
}

// NewSet generates a set that's retuned as pointer
func NewSet(buff []byte) *Set {

	return &Set{

		Fields: make([]Tag, 0),

		Content: Data{
			Type:   nil,
			Index:  0,
			Value:  "",
			Inners: make([]Data, 0),
		},

		Raw: buff,
	}
}

// recurse recreates a recursive Data datastructure representation of the file itself. It's taking every tags, subtags and data to make them one recursive datastructure of Type Data
func recurse(decoder *xml.Decoder, tagList []Tag) (Data, error) {

	data := newData(0)
	index := 0

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}

		switch tk := tok.(type) {

		case xml.StartElement:
			name := tk.Name.Local

			data.Inners = append(data.Inners, newData(index))
			data.Inners[index].Type = getTag(tagList, name)

			data.Inners[index], err = recurse(decoder, tagList)
			if err != nil {
				return Data{}, fmt.Errorf("recurse() -> %s", err)
			}

			index = index + 1

		case xml.EndElement:
			return data, nil

		case xml.CharData:
			data.Value = string(tk)

		default:
			return Data{}, fmt.Errorf("recurse() -> Unknown or unused Type encountered, got %T", tok)
		}
	}

	return data, nil
}

// newData return an empty Data Type
func newData(index int) Data {

	return Data{
		Type:   nil,
		Index:  index,
		Value:  "",
		Inners: make([]Data, 0),
	}
}

// findTags find and return the exhaustive list of unique tags found in the xml file
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

// getTag return a single tag from a tagList with a given name as string
func getTag(tagList []Tag, tok string) *Tag {

	for index, tag := range tagList {

		if tok == tag.Name {
			return &tagList[index]
		}
	}

	return nil
}

// tagExist checks if a tag exist in the whole list of tags
func tagExist(word string, tagList []Tag) bool {

	for _, tag := range tagList {

		if tag.Name == word {
			return true
		}
	}

	return false
}
