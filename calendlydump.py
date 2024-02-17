"""Check calendly appointments against a janky spreadsheet."""
import datetime
import requests
from absl import app
from absl import flags
from googleapiclient.discovery import build 
from google_auth_oauthlib.flow import InstalledAppFlow

FLAGS = flags.FLAGS
flags.DEFINE_string("calendly_keyfile", "calendly.key",
                    "Patch to Calendly token file")
flags.DEFINE_string("sheets_keyfile", "sheets.key.json",
                    "Patch to oauth 2.0 json keyfile")
flags.DEFINE_string("sheet_id", "",
                    "Google Sheet ID")

flags.DEFINE_string("email_range", "C2:C100",
                    "Range for emails")
flags.DEFINE_string("state_range", "I2:I100",
                    "range for state")

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

    def get_appointments(self):
        """ Return a dict of email to datetime of next appointments."""
        me = self.me()

        params = {
            "count": 100,
            "status": "active",
            "user": me["uri"],
            "min_start_time": datetime.datetime.now(
                tz=datetime.timezone.utc).strftime("%Y-%m-%dT%H:%M:%S%z")
        }
        events = self.get("scheduled_events", params)
        if "collection" not in events:
            raise CalendlyError("no collection: " + str(events))
        appts = {}
        for e in events["collection"]:
            start = datetime.datetime.fromisoformat(e["start_time"])
            uuid = e["uri"].split('/')[-1]
            ppl = self.get("scheduled_events/" + uuid + "/invitees")
            for p in ppl["collection"]:
                if p["email"] in appts:
                    print(f"More than one appointment for {p['email']}")
                    # Keep the soonest date
                    if start > appts[p["email"]]:
                        continue
                appts[p["email"]] = e["start_time"]

        return appts

class Sheet(object):
    def __init__(self, keyfile, sheet_id):
        self._sheet_id = sheet_id
        self._key = keyfile
        self._flow = InstalledAppFlow.from_client_secrets_file(
            keyfile, ["https://www.googleapis.com/auth/spreadsheets"])
        self._flow.run_local_server()
        self._creds = self._flow.credentials
    
    def get_range(self, range):
        sheets_svc = build('sheets', 'v4', credentials=self._creds)
        return sheets_svc.spreadsheets().values().get(
            spreadsheetId=self._sheet_id, range=range).execute()

def main(_):
    """Do the Thing."""
    print("Hello.")
    calendly_api_key = open(
        FLAGS.calendly_keyfile, "r", encoding="utf-8").readlines()[0].strip()

    c = Calendly(calendly_api_key)
    s = Sheet(FLAGS.sheets_keyfile, FLAGS.sheet_id)

    emails = s.get_range(FLAGS.email_range)['values']
    apps = c.get_appointments()

    for e in emails:
        if e[0] in apps:
            print(f"{e[0]}\t\t: {apps[e[0]]}")
        else:
            print(f"{e[0]}\t\tNo Appointment")


if __name__ == '__main__':
    app.run(main)
