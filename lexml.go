package lexml

import (
    "fmt"
    //"encoding/xml"
    "bufio"
)

/* STRUCTURES */

// Tags represent fields of a XML file, those are generated on the go
type Tag struct {
    Id int
    Name string
}

// Data is a recursive representation of a parsed XML file content
type Data struct {
    Type *Tag
    Index int
    Value string
    Inners []Data
}

// a DataSet is the full content of a parsed XML file
type DataSet struct {
    Fields []Data
}

type Coordinates struct {
    Start int
    End int
}

/* LOCAL FUNCTIONS */

func tagExist(word string, tagList []Tag) bool {
    for _, t := range tagList {
        if t.Name == word {
            return true
        }
    }
    return false
}

func isClosingTag(buff []byte, index int) bool {
    if buff[index+1] == '/' {
        return true
    } else {
        return false
    }
}

func findTags(buff []byte) ([]Tag, error ) {

    var check bool
    index := 0
    coordinates := make([]Coordinates, 0)
    tags := make([]Tag, 0)

    for ; index < len(buff); index = index + 1 {
        switch buff[index] {
        case '<':
            if index + 1 > len(buff) {
                return []Tag{}, fmt.Errorf("findTags() -> Out Of Bounds, got '%d' with buffer length '%d'", index, len(buff))
            }
            check = isClosingTag(buff[:], index)
            if check == true {
                index = index + 1
                continue
            }
            coordinates = append(coordinates, Coordinates{Start: index+1, End: 0})
        case '>':
            if check == true {
                check = false
                continue
            }
            coordinates[len(coordinates) - 1].End = index - 1
        }
    }
    
    for i := range coordinates {
        if len(tags) == 0 {
            tags = append(tags, Tag{
                Id: 0,
                Name: string(buff[coordinates[i].Start:coordinates[i].End]),
            })
        } else if tagExist(string(buff[coordinates[i].Start:coordinates[i].End]), tags) == false {
            tags = append(tags, Tag{
                Id: len(tags) - 1,
                Name: string(buff[coordinates[i].Start:coordinates[i].End]),
            })
        } else {
            continue
        }
    }
    return tags, nil
}

/* PUBLIC FUNCTIONS */

// this function is a lexer that convert the whole file into a DataSet datastructure
func Lexml(scanner *bufio.Scanner) DataSet {
    
    //set := DataSet{ Fields: make([]Data, 0) }
    contentBuffer := make([]byte, 0)
    
    tags, err := findTags(contentBuffer[:])
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(tags)

    return DataSet{}
}
