package lexml

import (
	"fmt"
    "bytes"
    "slices"
)

/* STRUCTURES */
// Set is the entrypoint struct to interact with this program with Fields as a list of Tag types, Content that represent the whole XML file, Raw wich represent the raw byte data of the file
type Set struct {
    Fields []Tag
    Content Data
    Raw []byte
}

// Data is a recursive representation of a parsed XML file content
type Data struct {
    Type *Tag
    Index int
    Value string
    Inners []Data
}

// Tags represent fields of a XML file, those are generated on the go
type Tag struct {
    Id int
    Name string
}

// locations is a superset of coordinates, it stores coordinates of the opening and the closing tags
type locations struct {
    Open coordinates
    Close coordinates
    Weight int
}

// coordinates represents coordinates of a single tag
type coordinates struct {
    Start int
    End int
}

// Parse convert the whole file into a DataSet datastructure
func (set Set) Parse() (Data, error) {
    
    tagList, tagLoc, err := getTags(set.Raw)
    if err != nil {
        return Data{}, fmt.Errorf("%s", err)
    }

    fmt.Printf("tagList: %v\n", tagList)

    // now that we generates the positions of very tags with their closing and opening locations, we can parse the content recursively as Data
    tagLoc = getWeights(tagLoc)

    /*for _, tag := range tagLoc {
        fmt.Println(tag.Weight)
    }*/

    return Data{}, nil
}

// NewSet generates a set that's retuned as pointer
func NewSet(buff []byte) *Set {
    return &Set{
        Fields: make([]Tag, 0),
        Content: Data{
            Type: nil,
            Index: 0,
            Value: "",
            Inners: make([]Data, 0),
        },
        Raw: buff[:],
    }
}

// getTags returns a list of Tag types wich defines the exhaustive list of tags. It also return a list of locations types wich gives indication on where each tags are in the buffer
func getTags(buff []byte) ([]Tag, []locations, error) {

    tagList := make([]Tag, 0)
    splitted := split(buff)
    loc := make([]locations, 0)

    for index := range splitted {
        start := splitted[index].Start
        end := splitted[index].End
        name := buff[start:end]
       
        if !(name[0] == '/') {
            
            if !tagExist(string(name), tagList) {
                tagList = append(tagList, Tag{Id: len(tagList), Name: string(name)})
            }

            open := coordinates{
                Start: start,
                End: end,
            }
            close, err := getClosTag(buff, splitted, index)
            if err != nil {
                return []Tag{}, []locations{}, fmt.Errorf("getTags() -> %s", err)
            }
            loc = append(loc, locations{
                Open: open,
                Close: close,
                Weight: 0,
            })
        }
    }

    return tagList, loc, nil
}

// split generate a coordinates list from every tag and tag ending bounds
func split(buff []byte) []coordinates {
    coo := make([]coordinates, 0)

    for index, char := range buff {
        switch char {
            case '<':
                coo = append(coo, coordinates{Start: index + 1, End: 0})
            case '>':
                coo[len(coo)-1].End = index
        }
    }
    return coo
}

// getClosTag retrieves the coordinates of every closing tags
func getClosTag(buff []byte, split []coordinates, index int) (coordinates, error) {

    start := split[index].Start
    end := split[index].End
    name := []byte{47}
    name = slices.Concat(name, buff[start:end])
    
    for indx := range split {
        next := buff[split[indx].Start:split[indx].End]
        if bytes.Equal(next, name) {
            return split[indx], nil
        }
    }
    
    return coordinates{}, fmt.Errorf("getClosTag() -> Closing tag not found for '%s'.", string(name))
}

// getWeight return the weight of each tags from their locations. This weight value is used to check wether a specific tag is into another one.
func getWeights(loc []locations) []locations {
    
    for index, tag := range loc {
        if index == 0 {
            tag.Weight = 0
            continue
        }
        tag.Weight = compWeights(index, loc)
        fmt.Println(tag.Weight)
    }
    return loc
}

// compWeights finds the weight of a specific locations
func compWeights(index int, loc []locations) int {
    
    weight := 0

    for _, tag := range loc {
        
        open := tag.Open.Start < loc[index].Open.Start
        close := tag.Close.End > loc[index].Open.End

        //fmt.Println(open, tag.Open.Start, loc[index].Open.Start, "|", close, tag.Close.End, loc[index].Open.End, "|", (open && close) )

        if open && close {
            weight = weight + 1
        }
    }
    return weight
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
