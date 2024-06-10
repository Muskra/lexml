package lexml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
  "slices"
	"strconv"
	"strings"

)

// Set type is the entrypoint struct to interact with this program with Fields as a list of Tag types, Content that represent the whole XML file, Raw wich represent the raw byte data of the file
type Set struct {
	Fields  []Tag
	Content Data
	Raw     []byte
}

// Data type is a recursive representation of a parsed XML file content
type Data struct {
	Type   *Tag
	Index  int
	Value  string
	Inners []Data
}

// DataAlt type is an altered Data Type without the Inners element, it's used when returning a list of specific Data elements from a LookupId, LookupName or LookupIndex functions
type DataAlt struct {
    Type *Tag
    Index int
    Value string
}

// Tags type represent fields of a XML file, those are generated on the go
type Tag struct {
	Id   int
	Name string
}

// Parse method convert the whole file into a DataSet datastructure
func (set Set) Parse() ([]Tag, Data, error) {


	reader := bytes.NewReader(set.Raw)
	decoder := xml.NewDecoder(reader)

	tagList := findTags(decoder)

	reader = bytes.NewReader(set.Raw)
	decoder = xml.NewDecoder(reader)

	content, err := genData(decoder, tagList,0)
	//133610189360577891
	//findProcess(content, "",0)
	findCreateKey(content,"",0)
	//displayIndex(content,"")

	if err != nil {
		return []Tag{}, Data{}, fmt.Errorf("Parse() -> %s", err)
	}

	return tagList, content, nil
}

// Alter method generate a DataAlt Type
func (data Data) Alter() DataAlt {
    return DataAlt{
        Type: data.Type,
        Index: data.Index,
        Value: data.Value,
    }
}

// LookupId method search recursively throught the Content of a given Set Type and returns a list of pointers to every Data elements that are equal to the given id
func (data Data) LookupId(id int) []DataAlt {

    dataList := make([]DataAlt, 0)
    givenId := 0

    if intEq(id, data.Type.Id) {
        dataList = append(dataList, data.Alter())
    }

    for index, dt := range data.Inners {

        givenId = dt.Type.Id

        if intEq(id, givenId) {
            dataList = append(dataList, dt.Alter())
        }

        dataList = slices.Concat(dataList, dt.LookupId(index))
    }

    return dataList
}

// LookupName method search recursively throught the Content of a given Set Type and returns a list of pointers to every Data elements that are equal to the given id
func (data Data) LookupName(name string) []DataAlt {

    dataList := make([]DataAlt, 0)
    givenName := ""

    if strEq(name, data.Type.Name) {
        dataList = append(dataList, data.Alter())
    }

    for _, dt := range data.Inners {

        givenName = dt.Type.Name

        if strEq(name, givenName) {
            dataList = append(dataList, dt.Alter())
        }

        dataList = slices.Concat(dataList, dt.LookupName(name))
    }

    return dataList
}

// LookupIndex method search recursively througt the Content of a given Set Type and returns the element present at a specific index of a given depth. depth is based on x and y is the index in the Inners Type
func (data Data) LookupIndex(depth int, x int, y int) DataAlt {

    retData := DataAlt{
        Type: nil,
        Index: 0,
        Value: "",
    }

    for index, dt := range data.Inners {

        if !(depth == x) {
            retData = dt.LookupIndex(depth+1, x, y)

        } else if (depth == x) && (index == y) {
            retData = dt.Alter()
            break
        }
    }

    return retData
}

func (data Data) PreFormatAll() []DataAlt {

    dataList := make([]DataAlt, 0)
    //dataList = append(dataList, data.Alter())
    for _, dt := range data.Inners {
        dataList = append(dataList, dt.Alter())

        dataList = slices.Concat(dataList, dt.PreFormatAll())
    }
    return dataList
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

func formatPrint(dataAltSlice []DataAlt) {
    for _, dt := range dataAltSlice {
        fmt.Printf("Id: %d\tName: %s\tIndex: %d\tValue: %s\n", dt.Type.Id, dt.Type.Name, dt.Index, dt.Value)
    }
}

func displayIndex(data Data, prefix string) {
	fmt.Printf("%sPosition: %d Contenu: %s TagName: %s\n", prefix, data.Index, data.Value, data.Type.Name)
	for _, inner := range data.Inners {
			displayIndex(inner, prefix+"  ")
	}
}

//in development.....
//si on passe 0 dans l'argument processIndex alors, on applique aucun filtrage sur les processus
func findCreateKey(data Data, prefix string, processIndex int){
	if data.Type.Name == "event"{
		dataProcessIndex, err := strconv.Atoi(data.Inners[0].Value)
    if err != nil {
        fmt.Println("Erreur de conversion:", err)
		}
		if dataProcessIndex == processIndex || processIndex == 0 {
			if data.Inners[4].Type.Name == "Operation" && data.Inners[4].Value == "RegCreateKey" {
				if data.Inners[7].Type.Name == "Detail" && strings.Contains(data.Inners[7].Value , "REG_CREATED_NEW_KEY"){
					fmt.Println("-----------------------------------------------------")
					fmt.Printf("%sTime of Day: %s \n", prefix, data.Inners[1].Value)
					fmt.Printf("%sProcess Name of Day: %s \n", prefix, data.Inners[2].Value)
					fmt.Printf("%sPath: %s \n", prefix, data.Inners[5].Value)
				}
			}
		}
	}
	for _, inner := range data.Inners {
		findCreateKey(inner,prefix+"  ",processIndex)
	}
}

// genData function recreates a recursive Data datastructure representation of the file itself. It's taking every tags, subtags and data to make them one recursive datastructure of Type Data
func genData(decoder *xml.Decoder, tagList []Tag, index int) (Data, error) {

	data := newData(index)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		switch tk := tok.(type) {

		case xml.StartElement:

			name := tk.Name.Local
			data.Inners = append(data.Inners, newData(index))
			data.Inners[len(data.Inners)-1], err = genData(decoder, tagList,len(data.Inners)-1)
			if err != nil {
				return Data{}, fmt.Errorf("recurse() -> %s", err)
			}
			data.Inners[len(data.Inners)-1].Type = getTag(tagList, name)


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

// newData function return an empty Data Type
func newData(index int) Data {

	return Data{
    Type:   &Tag{Id: 0, Name: "XMLROOT"},
		Index:  index,
		Value:  "",
		Inners: make([]Data, 0),
	}
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

// getId function returns the Id of a given Data Type
func getId(data Data) int {
    return data.Type.Id
}

// intEq function checks equality of two given Int values
func intEq(orig int, given int) bool {
    return (orig == given)
}

// getName function returns the Name of a given Data Type
func getName(data Data) int {
    return data.Type.Id
}

// strEq function checks equality of two given String values
func strEq(orig string, given string) bool {
    return bytes.Equal([]byte(orig), []byte(given))
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
