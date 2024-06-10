// Description: This package is used to create data for the LaMetric My Data DIY app.
// The data is requested from the LaMetric device and displayed on the screen.

package lametricmydatadiygo

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// https://help.lametric.com/support/solutions/articles/6000225467-my-data-diy

// MyDataFrame is a struct that holds the data for a single frame.
// If the text owerflows the frame, the duration will be ignored.
type MyDataFrame struct {
	// Text is the text to display
	Text string `json:"text,omitempty"`

	// Icon is a string because it can be a number or a URL
	// List of icons: https://developer.lametric.com/icons
	Icon string `json:"icon,omitempty"`

	// Duration is the time in seconds the frame should be displayed
	Duration int `json:"duration,omitempty"`

	// GoalData is the data for the goal widget
	GoalData *GoalData `json:"goalData,omitempty"`

	// ChartData is the data for the chart widget - A list of integers
	ChartData []int `json:"chartData,omitempty"`

	// expireTime is the time when the frame should be removed
	expireTime int64
}

// GoalData is a struct that holds the data for the goal widget
type GoalData struct {
	// Start is the start value
	Start int `json:"start,omitempty"`
	// Current is the current value
	Current int `json:"current,omitempty"`
	// End is the end value for the goal
	End int `json:"end,omitempty"`
	// Unit is the unit for the goal. Will be displayed after the current value.
	Unit string `json:"unit,omitempty"`
}

// MyDataFrames is a struct that holds a slice of MyDataFrame
type MyDataFrames struct {
	sync.Mutex
	Frames []MyDataFrame `json:"frames"`
}

// ToJson will convert the MyDataFrames to a JSON byte slice
func (m *MyDataFrames) ToJson() ([]byte, error) {
	return json.Marshal(m)
}

// AddFrame will add a frame to the slice
// It will return the MyDataFrames struct for chaining
// Example:
//
//	frames := lametric.MyDataFrames{}
//	frames.AddFrame(lametric.MyDataFrame{Text: "Hello World!"})
func (m *MyDataFrames) AddFrame(frame MyDataFrame) *MyDataFrames {
	m.Lock()
	m.Frames = append(m.Frames, frame)
	m.Unlock()
	return m
}

// RemoveFrame will remove a frame from the slice
// It will return the MyDataFrames struct for chaining
// If the index is out of bounds, the slice will not be changed
// Example:
//
//	frames.RemoveFrame(0) // Remove the first frame
func (m *MyDataFrames) RemoveFrame(index int) *MyDataFrames {
	if index < 0 || index >= len(m.Frames) {
		return m
	}
	m.Lock()
	m.Frames = append(m.Frames[:index], m.Frames[index+1:]...)
	m.Unlock()
	return m
}

// String will convert the MyDataFrames to a JSON string
// It will return the JSON representation of the frames
// Example:
//
//	fmt.Printf("Frames: %s\n", frames)
//
// Output:
//
//	Frames: {"frames":[{"text":"Hello World!"}]}
func (m *MyDataFrames) String() string {
	jsonData, _ := m.ToJson()
	return string(jsonData)
}

// HttpFunc is a function that can be used as a http.HandlerFunc
// It will return the frames as JSON fot the LaMetric My Data DIY app
// Example:
//
//	http.HandleFunc("/", frames.HttpFunc)
func (m *MyDataFrames) HttpFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := m.ToJson()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func NewMyDataTextFrame(text, icon string, opts ...MyDataFrameOptions) MyDataFrame {
	return MyDataFrame{
		Text: text,
		Icon: icon,
	}
}

type MyDataFrameOptions func(*MyDataFrame) error

// MyDataFrameOptions is a struct that holds the options for the MyDataFrame
type MyDataFrameCustomizer interface {
	Customize(*MyDataFrame) error
}

func (opt MyDataFrameOptions) Customize(m *MyDataFrame) error {
	return opt(m)
}

func WithExpiryInSeconds(seconds int) MyDataFrameCustomizer {
	return func(m *MyDataFrame) error {
		m.expireTime = time.Now().Add(time.Duration(seconds) * time.Second).Unix()
		return nil
	}
}
