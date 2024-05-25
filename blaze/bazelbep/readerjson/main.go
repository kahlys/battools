package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/text"
	"github.com/tidwall/gjson"

	bes "github.com/kahlys/battools/proto/build_event_stream"
)

func main() {
	_ = bes.BuildEventId_TargetCompletedId{
		Label: "label",
	}

	return

	fileName := flag.String("file", "build_events.json", "build events file to read")
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	events := &buildEvents{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		event := gjson.Parse(scanner.Text())
		id := eventID(event)
		switch id {
		case "started":
			events.AddBuildEvent(buildEvent{id: id, event: parseBuildStarted(event)})
		case "targetCompleted":
			events.AddBuildEvent(buildEvent{id: id, event: parseBuildTargetCompleted(event)})
		case "buildFinished":
			events.AddBuildEvent(buildEvent{id: id, event: parseBuildFinished(event)})
		case "testSummary":
			events.AddBuildEvent(buildEvent{id: id, event: parseTestSummary(event)})
		case "testResult":
			events.AddBuildEvent(buildEvent{id: id, event: parseTestResult(event)})
		}
	}

	// for _, event := range events {
	// 	for key, _ := range event.Map() {
	// 		// Do something with the key and value...
	// 		fmt.Println("Key:", key)
	// 	}
	// }

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(events)
}

type buildEvent struct {
	id    string
	event string
}

type buildEvents []buildEvent

func (b *buildEvents) AddBuildEvent(event buildEvent) {
	*b = append(*b, event)
}

func (b *buildEvents) String() string {
	res := ""
	for _, event := range *b {
		res += fmt.Sprintf("%s\n", strings.ToUpper(strings.TrimSpace(event.id)))
		res += text.Color.Sprint(text.FgHiBlack, fmt.Sprintf("%s\n\n", strings.TrimSpace(event.event)))
	}
	return strings.TrimSpace(res)
}

func eventID(event gjson.Result) string {
	for key := range event.Get("id").Map() {
		return key
	}
	return ""
}

func parseBuildStarted(event gjson.Result) string {
	data := struct {
		Started struct {
			StartTime time.Time `json:"startTime"`
			Command   string    `json:"command"`
		} `json:"started"`
	}{}

	err := json.Unmarshal([]byte(event.Raw), &data)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%s start %s", data.Started.Command, data.Started.StartTime.Format(time.DateTime))
}

func parseBuildTargetCompleted(event gjson.Result) string {
	data := struct {
		ID struct {
			TargetCompleted struct {
				Label string `json:"label"`
			} `json:"targetCompleted"`
		} `json:"id"`
		Completed struct {
			Success bool `json:"success"`
		} `json:"completed"`
	}{}

	err := json.Unmarshal([]byte(event.Raw), &data)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%v %v", data.ID.TargetCompleted.Label, data.Completed.Success)
}

func parseBuildFinished(event gjson.Result) string {
	data := struct {
		Finished struct {
			ExitCode struct {
				Name string `json:"name,omitempty"`
				Code int32  `json:"code,omitempty"`
			} `json:"exitCode,omitempty"`
			FinishTime time.Time `json:"finishTime"`
		} `json:"finished"`
	}{}

	err := json.Unmarshal([]byte(event.Raw), &data)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(
		"%s code:%d at %s",
		data.Finished.ExitCode.Name,
		data.Finished.ExitCode.Code,
		data.Finished.FinishTime.Format(time.DateTime),
	)
}

func parseTestSummary(event gjson.Result) string {
	data := struct {
		ID struct {
			TestSummary struct {
				Label string `json:"label"`
			} `json:"testSummary"`
		} `json:"id"`
		TestSummary struct {
			TotalRunCount int `json:"totalRunCount"`
			Passed        []struct {
				URI string `json:"uri"`
			} `json:"passed"`
			OverallStatus    string `json:"overallStatus"`
			TotalRunDuration string `json:"totalRunDuration"`
		} `json:"testSummary"`
	}{}

	err := json.Unmarshal([]byte(event.Raw), &data)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(
		"%s %s %d/%d %s",
		data.ID.TestSummary.Label,
		data.TestSummary.OverallStatus,
		len(data.TestSummary.Passed),
		data.TestSummary.TotalRunCount,
		data.TestSummary.TotalRunDuration,
	)
}

func parseTestResult(event gjson.Result) string {
	data := struct {
		ID struct {
			TestResult struct {
				Label string `json:"label"`
			} `json:"testResult"`
		} `json:"id"`
		TestResult struct {
			Status        string `json:"status"`
			ExecutionInfo struct {
				Strategy string `json:"strategy"`
			} `json:"executionInfo"`
			TestAttemptDuration string `json:"testAttemptDuration"`
		} `json:"testResult"`
	}{}

	err := json.Unmarshal([]byte(event.Raw), &data)
	if err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf(
		"%s %s %s %s",
		data.ID.TestResult.Label,
		data.TestResult.Status,
		data.TestResult.TestAttemptDuration,
		data.TestResult.ExecutionInfo.Strategy,
	)
}
