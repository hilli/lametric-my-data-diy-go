package lametricmydatadiygo_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	lametric "github.com/hilli/lametric-my-data-diy-go"
)

func Test_MyDataDIY(t *testing.T) {
	t.Run("Test MyDataFrame", func(t *testing.T) {
		frame := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
			GoalData: &lametric.GoalData{
				Start:   0,
				Current: 34,
				End:     100,
				Unit:    "%",
			},
		}

		if frame.Text != "Hello World!" {
			t.Errorf("Expected 'Hello World!' but got %s", frame.Text)
		}
		if frame.Icon != "64" {
			t.Errorf("Expected '64' but got %s", frame.Icon)
		}
		if frame.GoalData.Start != 0 {
			t.Errorf("Expected 0 but got %d", frame.GoalData.Start)
		}
		if frame.GoalData.Current != 34 {
			t.Errorf("Expected 34 but got %d", frame.GoalData.Current)
		}
		if frame.GoalData.End != 100 {
			t.Errorf("Expected 100 but got %d", frame.GoalData.End)
		}
		if frame.GoalData.Unit != "%" {
			t.Errorf("Expected '%%' but got %s", frame.GoalData.Unit)
		}
	})

	t.Run("Test MyDataFrame with empty GoalData", func(t *testing.T) {
		frame := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
		}

		if frame.Text != "Hello World!" {
			t.Errorf("Expected 'Hello World!' but got %s", frame.Text)
		}
		if frame.Icon != "64" {
			t.Errorf("Expected '64' but got %s", frame.Icon)
		}
		if frame.GoalData != nil {
			t.Errorf("Expected nil but got %v", frame.GoalData)
		}
	})

	t.Run("Test JSON returned is correct", func(t *testing.T) {
		frame := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
			GoalData: &lametric.GoalData{
				Start:   0,
				Current: 34,
				End:     100,
				Unit:    "%",
			},
			ChartData: []int{10, 20, 5},
		}

		frames := lametric.MyDataFrames{}
		frames.AddFrame(frame)

		json, err := frames.ToJson()
		if err != nil {
			t.Errorf("Expected nil but got %v", err)
		}
		if string(json) != `{"frames":[{"text":"Hello World!","icon":"64","goalData":{"current":34,"end":100,"unit":"%"},"chartData":[10,20,5]}]}` {
			t.Errorf("Expected JSON but got %s", json)
		}
		if frames.String() != `{"frames":[{"text":"Hello World!","icon":"64","goalData":{"current":34,"end":100,"unit":"%"},"chartData":[10,20,5]}]}` {
			t.Errorf("Expected JSON but got %s from the stringer", frames.String())
		}

	})

	t.Run("Test AddFrame and RemoveFrame", func(t *testing.T) {
		frame1 := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
		}

		frame2 := lametric.MyDataFrame{
			Text: "Not Hello World!",
			Icon: "61",
		}
		var frames lametric.MyDataFrames
		frames.AddFrame(frame1).AddFrame(frame2)
		if len(frames.Frames) != 2 {
			t.Errorf("Expected 2 but got %d", len(frames.Frames))
		}

		frames.RemoveFrame(0)
		if len(frames.Frames) != 1 {
			t.Errorf("Expected 1 but got %d", len(frames.Frames))
		}

		if frames.Frames[0].Text != "Not Hello World!" {
			t.Errorf("Expected 'Not Hello World!' but got %s", frames.Frames[0].Text)
		}

		frames.RemoveFrame(1)
		if len(frames.Frames) != 1 {
			t.Errorf("Expected 1 but got %d", len(frames.Frames))
		}
	})

	t.Run("Test httpFunc", func(t *testing.T) {
		frame := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
		}

		frames := lametric.MyDataFrames{}
		frames.AddFrame(frame)

		ts := httptest.NewServer(http.HandlerFunc(frames.HttpFunc))
		defer ts.Close()

		resp, err := ts.Client().Get(ts.URL)
		if err != nil {
			t.Errorf("Expected nil but got %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 but got %d", resp.StatusCode)
		}
		if resp.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected 'application/json' but got %s", resp.Header.Get("Content-Type"))
		}
	})
}

func Test_Push(t *testing.T) {
	t.Run("Test Push", func(t *testing.T) {
		frame := lametric.MyDataFrame{
			Text: "Hello World!",
			Icon: "64",
		}
		frames := lametric.MyDataFrames{}
		frames.AddFrame(frame)

		token := "77charactersofrandomness"
		basicauth := "Basic ZGV2Ojc3Y2hhcmFjdGVyc29mcmFuZG9tbmVzcw=="

		ts := httptest.NewServer(http.Handler(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Authorization") != basicauth {
					t.Errorf("Expected '%s' but got %s", basicauth, r.Header.Get("Authorization"))
				}
				w.WriteHeader(http.StatusOK)
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST but got %s", r.Method)
				}
				if r.Body == nil {
					t.Errorf("Expected a body but got nil")
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Errorf("Expected to be able to read the body but got an error %v", err)
				}
				if string(body) != `{"frames":[{"text":"Hello World!","icon":"64"}]}` {
					t.Errorf("Expected a body %s, but got 'adasd'", r.Body)
				}
			},
		)))
		defer ts.Close()

		err := frames.Push(ts.URL, token)

		// err := frames.Push(url, token)
		if err != nil {
			t.Errorf("Expected nil but got %v", err)
		}
	})
}
