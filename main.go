package main

import (
	"fmt"
	"os"
	"pmp/lexml"
    "pmp/repl"
)

func abort(err error) {
	fmt.Printf("< ! > Error: %s", err)
}

func main() {

	/*
	   {0 procmon}
	   {0 processlist}
	   {0 process}
	   {0 ProcessIndex}
	   {0 ProcessId}
	   {0 ParentProcessId}
	   {0 ParentProcessIndex}
	   {0 AuthenticationId}
	   {0 CreateTime}
	   {0 FinishTime}
	   {0 IsVirtualized}
	   {0 Is64bit}
	   {0 Integrity}
	   {0 Owner}
	   {0 ProcessName}
	   {0 ImagePath}
	   {0 CommandLine}
	   {0 CompanyName}
	   {0 Version}
	   {0 Description}
	   {0 modulelist}
	   {0 module}
	   {0 Timestamp}
	   {0 BaseAddress}
	   {0 Size}
	   {0 Path}
	   {0 Company}
	   {0 eventlist}
	   {0 event}
	   {0 Time_of_Day}
	   {0 Process_Name}
	   {0 PID}
	   {0 Operation}
	   {0 Result}
	   {0 Event_Class}
	   {0 Category}
	   {0 Detail}
	   {0 Command_Line}
	   {0 Date___Time}
	   {0 Completion_Time}
	   {0 TID}
	   {0 Parent_PID}
	   {0 Session}
	   {0 User}
	   {0 Authentication_ID}
	   {0 Virtualized}
	   {0 Relative_Time}
	   {0 Duration}
	   {0 Sequence}
	   {0 Image_Path}
	   {0 Architecture}
	*/

	var err error

	//file, err := os.ReadFile("/Users/tartintosh/Documents/devroom/go/pmp/Logfile.XML")
	file, err := os.ReadFile("/Users/tartintosh/Documents/devroom/go/pmp/simple.xml")
	if err != nil {
		abort(fmt.Errorf("main.go line 13 -> %s", err))
	}

	set := lexml.NewSet(file)
	set.Fields, set.Content, err = set.Parse()
	if err != nil {
		abort(fmt.Errorf("main.go line 23 -> %s", err))
	}

    //set.Content.DisplayIndex("")
	//lexml.FormatPrint(set.Content.PreFormatAll())
	repl.TestLaunch(set)
}
