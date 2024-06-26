package lexml

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

// Set type is the entrypoint struct to interact with this program with Fields as a list of Tag types, Content that represent the whole XML file, Raw wich represent the raw byte data of the file
type Set struct {
	Fields  []Tag
	Content Data
	Raw     []byte
}

// Parse method convert the whole file into a DataSet datastructure
func (set Set) Parse() ([]Tag, Data, error) {

	reader := bytes.NewReader(set.Raw)
	decoder := xml.NewDecoder(reader)

	tagList := findTags(decoder)

	reader = bytes.NewReader(set.Raw)
	decoder = xml.NewDecoder(reader)

	content, err := genData(decoder, tagList, 0)
	if err != nil {
		return []Tag{}, Data{}, fmt.Errorf("Parse() -> %s", err)
	}

	return tagList, content, nil
}

// NewSet function generates a set that's retuned as pointer
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
