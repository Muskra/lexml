package lexml

import (
	"fmt"
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
    
    //tagList
    _, tagLoc, err := getTags(set.Raw)
    if err != nil {
        return Data{}, fmt.Errorf("%s", err)
    }

    //fmt.Printf("tagList: %v\n", tagList)
    for _, loc := range tagLoc {
        open := set.Raw[loc.Open.Start:loc.Open.End]
        close := set.Raw[loc.Close.Start:loc.Close.End]
        fmt.Println("open:", loc.Open, open, string(open), "\nclose:", loc.Close, close, string(close))
    }

    tagLoc = getWeights(tagLoc)
    tagLoc = make([]locations, 0)
    for _, l := range tagLoc {

        fmt.Println(l)
    }
    //getData(set.Raw, tagList, tagLoc)

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
    loc := make([]locations, 0)
    
    splitted := split(buff)
    tagList = getListTags(buff, splitted)
    open := getOpenTags(buff, splitted)
    close := getCloseTags(buff, splitted)

    loc, err := locAppend(buff, open, close)
    if err != nil {
        return nil, nil, fmt.Errorf("getTags() -> %s", err)
    }

    return tagList, loc[:], nil
}

// getData returns a Data type that recursively contains the tags and their contained informations
func getData(buff []byte, tags []Tag, loc []locations) {
    
    for _, lc := range loc {
        fmt.Println(string(buff[lc.Open.End:lc.Close.Start]))
    }
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

// getListTags return a list of every tag found in the file
func getListTags(buff []byte, split []coordinates) []Tag {

    tagList := make([]Tag, 0)

    for index := range split {

        start := split[index].Start
        end := split[index].End
        name := buff[start:end]
       
        if !(name[0] == '/') {
            if !tagExist(string(name), tagList) {
                tagList = append(tagList, Tag{Id: len(tagList), Name: string(name)})
            }
        }
    }
    return tagList
}

// getOpenTags return a list of every opening tag found in the file
func getOpenTags(buff []byte, split []coordinates) []coordinates {

    open := make([]coordinates, 0)

    for index := range split {

        start := split[index].Start
        end := split[index].End
        name := buff[start:end]
       
        if !(name[0] == '/') {
            open = append(open, coordinates{
                Start: start,
                End: end,
            })
        }
    }
    return open
}

// getClosTag retrieves the coordinates of every closing tags
func getCloseTags(buff []byte, split []coordinates) []coordinates {

    close := make([]coordinates, 0)

    for index := range split {
        start := split[index].Start
        end := split[index].End
        name := buff[start:end]

        if name[0] == '/' {
            close = append(close, coordinates{
                Start: start,
                End: end,
            })
        }
    }
    return close
}

// locAppend merge two []coordinates lists into a list of locations type
func locAppend(buff []byte, open []coordinates, close []coordinates) ([]locations, error) {
    
    if len(open) != len(close) {
        return nil, fmt.Errorf("locAppend() -> Malformed XML file, got '%d' opening tags with '%d' closing tags.", len(open), len(close))
    }

    loc := make([]locations, 0)
    
    for index, op := range open {
        idx := index
        for ; idx < len(close); idx = idx + 1 {
            cl := close[idx]
            if slices.Equal(buff[op.Start:op.End], buff[cl.Start+1:cl.End]) {
                loc = append(loc, locations{
                    Open: op,
                    Close: close[idx],
                })
                break
            }
        }
    }
    return loc[:], nil
}

// getWeight return the weight of each tags from their locations. This weight value is used to check wether a specific tag is into another one.
func getWeights(loc []locations) []locations {
    
    for index, tag := range loc {
        if index == 0 {
            tag.Weight = 0
            continue
        }
        loc[index].Weight = countWeight(index, loc)
        fmt.Println(loc[index].Weight)
    }
    return loc
}

// compWeights finds the weight of a specific locations
func countWeight(index int, loc []locations) int {

    chkOpn := loc[index].Open
    chkCls := loc[index].Close

    weight := 0

    /*
    
    check if the actual position is between any other positions
    we counts the number of tags that are higher value than the actual position

    */

    for _, check := range loc {
        cmp := (chkOpn.Start > check.Open.End) && (chkCls.End < check.Close.Start)
        if cmp {
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
