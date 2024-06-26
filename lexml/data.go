package lexml

import (
    "fmt"
    "slices"
    "strings"
)

// Data type is a recursive representation of a parsed XML file content
type Data struct {
	Type   *Tag
	Index  int
	Value  string
	Inners []Data
}

// DataAlt type is an altered Data Type without the Inners element, it's used when returning a list of specific Data elements from a LookupId, LookupName or LookupIndex functions
type DataAlt struct {
	Type  *Tag
	Index int
	Value string
}

// Alter method generate a DataAlt Type
func (data Data) Alter() DataAlt {
	return DataAlt{
		Type:  data.Type,
		Index: data.Index,
		Value: data.Value,
	}
}

// LookupId method search recursively throught the Content of a given Set Type and returns a list of pointers to every Data elements that are equal to the given id
func (data Data) LookupId(id int) []DataAlt {

	dataList := make([]DataAlt, 0)
	givenId := 0

	if IntEq(id, data.Type.Id) {
		dataList = append(dataList, data.Alter())
	}

	for index, dt := range data.Inners {

		givenId = dt.Type.Id

		if IntEq(id, givenId) {
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

	if StrEq(name, data.Type.Name) {
		dataList = append(dataList, data.Alter())
	}

	for _, dt := range data.Inners {

		givenName = dt.Type.Name

		if StrEq(name, givenName) {
			dataList = append(dataList, dt.Alter())
		}

		dataList = slices.Concat(dataList, dt.LookupName(name))
	}

	return dataList
}

// LookupIndex method search recursively througt the Content of a given Set Type and returns the element present at a specific index of a given depth. depth is based on x and y is the index in the Inners Type
func (data Data) LookupIndex(depth int, x int, y int) DataAlt {

	retData := DataAlt{
		Type:  nil,
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

func (data Data) DisplayIndex(prefix string) {
	fmt.Printf("%sPosition: %d Contenu: %s TagName: %s\n", prefix, data.Index, data.Value, data.Type.Name)
	for _, inner := range data.Inners {
			inner.DisplayIndex(fmt.Sprintf("%s\t", prefix))
	}
}

func FormatPrint(dataAltSlice []DataAlt) {
	for _, dt := range dataAltSlice {
		fmt.Printf("Id: %d\tName: %s\tIndex: %d\tValue: %s\n", dt.Type.Id, dt.Type.Name, dt.Index, strings.Trim(dt.Value, "\n"))
	}
}

// newData function return an empty Data Type
func NewData(index int) Data {

	return Data{
		Type:   &Tag{Id: 0, Name: "XMLROOT"},
		Index:  index,
		Value:  "",
		Inners: make([]Data, 0),
	}
}

// in development.....
// si on passe 0 dans l'argument processIndex alors, on applique aucun filtrage sur les processus
// je pense que ça devrait être substitué par un nouveau type pour éviter d'avoir 150 000 option
// ADD DESCRIPTION OF THE FUNCTION FOR THE DOC

// getId function returns the Id of a given Data Type
func getId(data Data) int {
	return data.Type.Id
}

// getName function returns the Name of a given Data Type
func getName(data Data) int {
	return data.Type.Id
}
