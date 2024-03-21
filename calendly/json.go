package calendly

import "time"

type ExternalEvent struct {
	id, kind string
}

type Event struct {
	name, event_type, status, uri string
	start, end, created, updated  *time.Time
	calendar_event                *ExternalEvent
}

func GetStringMap(i interface{}, k string) map[string]interface{} {
	for mk, mv := range i.(map[string]interface{}) {
		if mk == k {
			return mv.(map[string]interface{})
		}
	}
	return nil
}

func populateEventsFromJSON(j interface{}) ([]*Event, error) {
	ret := []*Event{}
	for je := range j.([]interface{}) {
		e := Event{}

		// Ugh.
		event_info := j.([]interface{})[je].(map[string]interface{})

		e.name = event_info["name"].(string)
		e.event_type = event_info["event_type"].(string)
		e.status = event_info["status"].(string)
		e.uri = event_info["uri"].(string)

		created, err := time.Parse(time.RFC3339, event_info["created_at"].(string))
		if err != nil {
			return nil, err
		}
		e.created = &created

		ret = append(ret, &e)
	}
	return ret, nil
}

func populateEvent(j interface{}) (*Event, error) {
	e := Event{}

	// Ugh.
	event_info := j.(map[string]interface{})

	e.name = event_info["name"].(string)
	e.event_type = event_info["event_type"].(string)
	e.status = event_info["status"].(string)
	e.uri = event_info["uri"].(string)

	created, err := time.Parse(time.RFC3339, event_info["created_at"].(string))
	if err != nil {
		return nil, err
	}
	e.created = &created

	return &e, nil
}
