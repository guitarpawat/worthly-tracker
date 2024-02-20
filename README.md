# 💹 Worthly Tracker
The net worth tracker that based on my requirement.

## Before running
1. Create sqlite file for database storage.
2. Create yaml config file as below:
```yaml
datasource:
  uri: file:C:\Users\path-to-file\tracker.sqlite
server:
  port: 8080
```
3. Compile the program with CGO enabled.
4. Run the program and specified a location to config file as the argument. For example: `main.exe tracker.yaml`

## Features

### MVP 1
- [x] Record date
- [ ] Manage assets
- [x] Support for bought value offset

### MVP 2
- [ ] Delete data
- [ ] Edit data
- [ ] Data summary

### MVP 3
- [ ] Data visualization
- [ ] Life Goals

### MVP 4
- [ ] Customizable data summary and visualization

### To Do
- [ ] Move documentation from Notion to the code (probably using swagger for API and godoc for DTO)
- [ ] Refactor frontend code
- [ ] Create automated frontend regression tests (probably using selenium or chromedp)

## Why SQLite?
1. Use Dropbox, Google Drive as a free cloud backup.
2. I estimate myself to record with 50 entries/time and 1 time/month. The estimate over my lifetime will be: 
`50 entries * 12 months * 80 years = 48,000 entries`, which can be handled by sqlite without needs of other database system.

## Contribution
Since I only work on this project to suite my needs, I will only accept PR on security issues or bugs. If you want to modify the code
and make the better version for yourself, please feel free to fork this project.