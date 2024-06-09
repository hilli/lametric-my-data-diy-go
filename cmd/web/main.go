package main

import (
	"fmt"
	"log"
	"net/http"

	lametric "github.com/hilli/lametric-my-data-diy-go"
)

func main() {
	frame1 := lametric.MyDataFrame{
		Text: "Hello World!",
		Icon: "64",
		GoalData: &lametric.GoalData{
			Start:   0,
			Current: 34,
			End:     100,
			Unit:    "%",
		},
	}

	savings1 := lametric.MyDataFrame{
		Text: "Savings",
		Icon: "61",
	}
	savings2 := lametric.MyDataFrame{
		Icon: "34",
		GoalData: &lametric.GoalData{
			Start:   0,
			Current: 1500,
			End:     2000,
			Unit:    "$",
		},
	}
	savings3 := lametric.MyDataFrame{
		ChartData: []int{10, 20, 5, 40, 20, 15, 20, 34},
	}

	frames := &lametric.MyDataFrames{
		Frames: []lametric.MyDataFrame{frame1, savings1},
	}

	// or

	frames.AddFrame(savings2).AddFrame(savings3)
	frames.RemoveFrame(0) // Remove the first frame

	// Add the LaMetric frames to the default mux
	http.HandleFunc("/", frames.HttpFunc)

	fmt.Printf("Frames: %s\n", frames)

	log.Println("Server listening on port http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
