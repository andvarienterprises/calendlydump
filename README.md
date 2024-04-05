Dump calendly appointments to a spreadsheet (or stdout)

Usage:

```
# Get the tool
go install github.com/andvarienterprises/calendlydump@0.0.1

# Print info about me on calendly (mainly to see if the API key is working)
calendlydump --calendly-auth-token=my.token.file me

# Print a CSV format list of all events, one line per event
calendlydump --calendly-auth-token=my.token.file dump events

# Print a CSV format list of all invites, one line per invite
# i.e. if there is an event with 4 invitees, it will print 4 lines
# for that event, one per invitee
calendlydump --calendly-auth-token=my.token.file dump invites
```

Useful for pasting to useful tools or file, or for use with [sheet](https://github.com/gerrowadat/sheet).
