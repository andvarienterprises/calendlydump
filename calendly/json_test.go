package calendly

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
	"time"
)

func rawJsonInterfaceFromFile(fn string) (interface{}, error) {
	raw_json, err := os.ReadFile("testdata/" + fn)
	if err != nil {
		return nil, err
	}
	var ret interface{}
	err = json.Unmarshal(raw_json, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func getTime(t string) *time.Time {
	ret, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil
	}
	return &ret
}
func Test_populateEvent(t *testing.T) {
	type args struct {
		test_fn string
	}
	tests := []struct {
		name    string
		args    args
		want    *Event
		wantErr bool
	}{
		{
			name: "1 Event, ok",
			args: args{test_fn: "1event"},
			want: &Event{
				name:           "My Lovely Meeting",
				event_type:     "https://api.calendly.com/event_types/bbbb-bbbb-bbbb",
				status:         "active",
				uri:            "https://api.calendly.com/scheduled_events/aaaa-aaaa-aaaa",
				start_time:     getTime("2023-10-24T15:00:00.000000Z"),
				end_time:       getTime("2023-10-24T15:15:00.000000Z"),
				created_at:     getTime("2023-10-21T19:56:09.501460Z"),
				updated_at:     getTime("2023-10-21T19:56:11.453361Z"),
				calendar_event: nil,
			},
			wantErr: false,
		},
		{
			name: "1 Event, with gcal",
			args: args{test_fn: "1event-gcal"},
			want: &Event{
				name:       "My Lovely Meeting",
				event_type: "https://api.calendly.com/event_types/bbbb-bbbb-bbbb",
				status:     "active",
				uri:        "https://api.calendly.com/scheduled_events/aaaa-aaaa-aaaa",
				start_time: getTime("2023-10-24T15:00:00.000000Z"),
				end_time:   getTime("2023-10-24T15:15:00.000000Z"),
				created_at: getTime("2023-10-21T19:56:09.501460Z"),
				updated_at: getTime("2023-10-21T19:56:11.453361Z"),
				calendar_event: &ExternalEvent{
					external_id: "aaaaa",
					kind:        "google",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := rawJsonInterfaceFromFile(tt.args.test_fn)
			if err != nil {
				t.Errorf("Error reading test data: %v", err)
				return
			}
			got, err := populateEvent(i)
			if (err != nil) != tt.wantErr {
				t.Errorf("populateEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("populateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_populateEventsFromJSON(t *testing.T) {
	type args struct {
		test_fn string
	}
	tests := []struct {
		name           string
		args           args
		wantEventCount int
		wantErr        bool
	}{
		{
			name:           "2 Events",
			args:           args{test_fn: "2events"},
			wantEventCount: 2,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := rawJsonInterfaceFromFile(tt.args.test_fn)
			if err != nil {
				t.Errorf("Error reading test data: %v", err)
				return
			}
			got, err := populateEventsFromJSON(i)
			if (err != nil) != tt.wantErr {
				t.Errorf("populateEventsFromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantEventCount != len(got) {
				t.Errorf("populateEventsFromJSON() = %v, want %v", len(got), tt.wantEventCount)
			}
		})
	}
}
