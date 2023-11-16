# timecop
A time blocking service that can be run locally or remotely as a server


## TODO
- [ ] Setup wizard
    - Decide when to prompt for next day's schedule
    - Provide API key?
- [ ] Move from JSON/REST to gRPC/websocket?
- [ ] Display schedule using markdown renderer (Charm's Glamour)
- [ ] Calendar sync (calendly, cal.com, and/or google calendar, apple calendar, etc.)
- [ ] GUI to edit schedules
- [ ] tomorrow.csv file to load up for the next day
- [ ] allow user to give list of tasks with priority and length of time, and the program will perform topological sort to create a schedule
    - [ ] Allow user to cycle through multiple options and pick the desired one
- [ ] Switch schedule-builder files from csv to toml
- [ ] Todoist integration
    - [ ] Export to Todoist account upon creation and completion
    - [ ] Import today's tasks from Todoist into a toml file to be added to?

## Notes To Self
- consider using a log rather than crossing out tasks/blocks (easier to implement and better for post-analysis)
