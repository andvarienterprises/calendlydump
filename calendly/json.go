package calendly

import "time"

type ExternalEvent struct {
	external_id, kind string
}

type Event struct {
	name, event_type, status, uri                string
	start_time, end_time, created_at, updated_at *time.Time
	calendar_event                               *ExternalEvent
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

	if c, ok := j.(map[string]interface{})["collection"]; ok {
		for _, je := range c.([]interface{}) {
			e, err := populateEvent(je)

			if err != nil {
				return nil, err
			}
			ret = append(ret, e)
		}
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

	if ce, ok := event_info["calendar_event"].(map[string]interface{}); ok {
		//ce := event_info["calendar_event"].(map[string]interface{})
		e.calendar_event = &ExternalEvent{
			external_id: ce["external_id"].(string),
			kind:        ce["kind"].(string),
		}
	}

	created_at, err := time.Parse(time.RFC3339, event_info["created_at"].(string))
	if err != nil {
		return nil, err
	}
	e.created_at = &created_at

	start_time, err := time.Parse(time.RFC3339, event_info["start_time"].(string))
	if err != nil {
		return nil, err
	}
	e.start_time = &start_time

	end_time, err := time.Parse(time.RFC3339, event_info["end_time"].(string))
	if err != nil {
		return nil, err
	}
	e.end_time = &end_time

	if updated_at, err := time.Parse(time.RFC3339, event_info["updated_at"].(string)); err == nil {
		e.updated_at = &updated_at
	}

	return &e, nil
}
