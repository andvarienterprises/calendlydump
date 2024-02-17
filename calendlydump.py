"""Dump Calendly appointments to stdout"""
import datetime
import requests
from absl import app
from absl import flags

FLAGS = flags.FLAGS
flags.DEFINE_string("calendly_keyfile", "calendly.key",
                    "Patch to Calendly token file")

class CalendlyError(Exception):
    """Generic Error when talking to Calendly."""

class Calendly(object):
    def __init__(self, key):
        self._key = key

    def get(self, method, params=None):
        """Generic request to Calendly"""
        headers = {"authorization": "Bearer " + self._key}
        r = requests.get("https://api.calendly.com/" + method,
                        params=params or {},
                        headers=headers,
                        timeout=10)
        return r.json()


    def me(self):
        """Get the info of the user the API key belongs to."""
        me = self.get("users/me")
        if 'resource' not in me:
            raise CalendlyError("no resource: " + str(me))
        return me["resource"]

    def _get_appts(self, params):
        me = self.me()
        params.update({
            "count": 100,
            "status": "active",
            "user": me["uri"],
        })
        events = self.get("scheduled_events", params)
        if "collection" not in events:
            raise CalendlyError("no collection: " + str(events))
        appts = {}
        for e in events["collection"]:
            uuid = e["uri"].split('/')[-1]
            etype = e["event_type"].split('/')[-1]
            ppl = self.get("scheduled_events/" + uuid + "/invitees")["collection"]
            appts[uuid] = (ppl[0]["email"], etype, e["start_time"])
        return appts

    def get_all_appointments(self):
        """ Return All Appointments."""
        return self._get_appts({})

    def get_future_appointments(self):
        """ Return a dict of email to datetime of next appointments."""
        params = {
            "min_start_time": datetime.datetime.now(
                tz=datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S%z")
        }
        return self._get_appts(params)

def main(_):
    """Do the Thing."""
    calendly_api_key = open(
        FLAGS.calendly_keyfile, "r", encoding="utf-8").readlines()[0].strip()

    c = Calendly(calendly_api_key)

    apps = c.get_all_appointments()

    for a in apps:
        print(f"{a},{apps[a][0]},{apps[a][1]},{apps[a][2]}")


if __name__ == '__main__':
    app.run(main)
