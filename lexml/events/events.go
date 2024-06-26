package event

import (

)

type Event struct {
    ProcessIndex int
    CreateTime int
    // if true output Data
    // if false output string
    SearchFor bool
}

// garder la boucle for de récursion, faire un tronc commun pour les recherches d'évènements et séparer les éléments spécifiques
func (event Event) FindEvent(data Data) {

    for ...
    switch event.Data.Type.Name {
        case "process":
             if event.SearchFor {
                event.Process()
            } else {
                // returns a Data type instead of the print
            }
        case "pid":
            event.Pid()
        defaut:
    }

    for _, inner.........
}



	}
}

func (event Event) EventPrintProcess() DataAlt {
	if data.Type.Name == "process" {
		dataProcessIndex, err := strconv.Atoi(data.Inners[1].Value)
		if err != nil {
			fmt.Println("Erreur de conversion:", err)
		}
		dataCreateTime, err := strconv.Atoi(data.Inners[5].Value)
		if err != nil {
			fmt.Println("Erreur de conversion:", err)
	    }   


            return fmt.Sprintf("")
        }
        
       		if dataProcessIndex == processIndex || processIndex == 0 {
			if dataCreateTime >= createTime || createTime == 0 {
				fmt.Println("-----------------------------------------------------")
				fmt.Printf(v"%sCreateTime: %s \n", prefix, data.Inners[5].Value)
				fmt.Printf("%sProcess ID: %s \n", prefix, data.Inners[1].Value)
				fmt.Printf("%sProcess Name: %s \n", prefix, data.Inners[11].Value)
				fmt.Printf("%sImage path: %s \n", prefix, data.Inners[12].Value)
				fmt.Printf("%sCommand Line: %s \n", prefix, data.Inners[13].Value)
				fmt.Printf("%sDescription: %s \n", prefix, data.Inners[16].Value)
			}
		}
}

        func (event Event) EventDataOfProcess() Data {
            //....
        }



