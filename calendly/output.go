package calendly

import "fmt"

var EventCSVHeader = "event_name,event_type,status,uri,created_at,updated_at,start_time,end_time"
var EventInvitesCSVHeader = EventCSVHeader + ",email,name,status"

func (e *Event) String() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", e.name, e.event_type, e.status, e.uri, e.created_at, e.updated_at, e.start_time, e.end_time)

}

func EventsCSV(e []*Event) string {
	ret := EventCSVHeader + "\n"
	for _, ev := range e {
		ret += ev.String() + "\n"
	}
	return ret
}

func EventInvitesCSV(e []*Event) string {
	// For each event, get the invites and return a line per invitee
	ret := EventInvitesCSVHeader + "\n"
	for _, ev := range e {
		for _, inv := range ev.invitees {
			ret += (ev.String() + "," + inv.email + "," + inv.name + "," + inv.status + "\n")
		}
	}
	return ret
}
