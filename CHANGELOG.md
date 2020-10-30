# Display WebService Changelog

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Please note, that this project, while following numbering syntax, it DOES NOT
adhere to [Semantic Versioning](http://semver.org/spec/v2.0.0.html) rules.

## Types of changes

* ```Added``` for new features.
* ```Changed``` for changes in existing functionality.
* ```Deprecated``` for soon-to-be removed features.
* ```Removed``` for now removed features.
* ```Fixed``` for any bug fixes.
* ```Security``` in case of vulnerabilities.

## [2020.4.1.30] - 2020-10-30

### Changed
- when closing order, terminal_input_login is not closing

## [2020.4.1.26] - 2020-10-26

### Fixed
- fixed one more leaking goroutine bug

## [2020.4.1.24] - 2020-10-24

### Fixed
- fixed leaking goroutine bug when opening sql connections, the right way is this way


## [2020.4.1.15] - 2020-10-15

### Fixed
- when changing user, proper sending Pcs="0"

## [2020.4.1.12] - 2020-10-12

### Fixed
- when ending cutting, click events removed from js, they caused duplicates
- bad ip address after ending cutting

### Changed
- time is everywhere in HH:MM
- saving to K2 allowed

## [2020.4.1.11] - 2020-10-11

### Fixed
- added windows.focus() to js files, where scanning from barcode reader or rfid reader is present

### Changed
- what is scanned (including meta characters like Enter, etc.) is send to backend
- parsing is processed in backend to know, what was scanned
- updated parsed value is send back to frontend

## [2020.4.1.8] - 2020-10-08

### Added
- removing "/R" from scanned order barcode

## [2020.4.1.7] - 2020-10-07

### Added
- reading order from K2
- reading pcs from K2
- saving data to K2 database (final saving in disabled, just inputting info)
- creating terminal_input_idle in zapsi database
- closing terminal_input_idle in zapsi database

### Fixed
- some fixes in javascript and go functions

## [2020.4.1.6] - 2020-10-06

### Added
- screens and behavior for idles
- screens and behavior for end cutting
- checking device in zapsi database
- checking logged user in zapsi database
- creating order in zapsi database
- creating idle in zapsi database
- getting idles from zapsi database
- checking user in zapsi database
- closing order in zapsi database
- closing idle in zapsi database


### Changed
- proper html formatting
- proper js formatting
- move go code to proper files
    - main.go contains code for web service
    - zapsi.go contains code for communicating with zapsi database
    - k2.go contains code for communicating with k2 database
    - log.go contains code for logging
    - pages.go contains code for pages loading
    - origin.go contains code for initial checking
- updated logging
    - one goroutine periodically checks for [ipaddress+device name] from zapsi database
    - when logging is called, device name is displayed like this ```[CNC1] --INF-- Saving data to K2```

## [2020.4.1.5] - 2020-10-05

### Added
- screens and behavior for user change
- screens and behavior for user break

### Changed
- some javascript changes

## [2020.4.1.2] - 2020-10-02

### Added
- initial commit
- added checking for device from ip address
- added checking for user
- added checking order
- added screen for inserting barcode
- added main screen
- added links to buttons on home screen