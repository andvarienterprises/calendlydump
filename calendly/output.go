package calendly

import "fmt"

var CSVHeader = "name,event_type,status,uri,created_at,updated_at,start_time,end_time"

func (e *Event) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", e.name, e.event_type, e.status, e.uri, e.created_at, e.updated_at, e.start_time, e.end_time)

}

func EventsCSV(e []*Event) string {
	ret := CSVHeader + "\n"
	for _, ev := range e {
		ret += ev.String() + "\n"
	}
	return ret
}
