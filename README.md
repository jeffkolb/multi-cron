# Multi-Cron

Multi-Cron is a cron runner written in Go, which utilizes version 2.0 of Rob Figueirdo's open source [cron Go library](https://github.com/robfig/cron/tree/v2)

It will schedule one or more cron-style jobs which run asynchronously with stdio and stderr redirected to the console.

The cron schedule format used is slicker than regular old cron in that it has a minimum resolution of 1 second, timezone support, and useful frequencies aliases.

---

## Usage

This application is tailor-made to run as the entrypoint (`PID 1`) of a Docker container and run 1 or more
applications/scripts,
configured via environmental variables, but it can also run as a stand-alone executable if you want to do that for some reason.

```bash
sudo docker run -d \
    -e CRON_SCH_1="!@daily" -e CRON_CMD_1="echo" -e CRON_ARGS_1="I ran at midnight" \
    -e CRON_SCH_2="0 0 20 * * *" -e CRON_CMD_2="/path/to/some/script" \
    -e CRON_SCH_3="@every 3h" -e CRON_CMD_3="python" -e CRON_ARGS_3="/path/to/python.py" \
    jeffkolb/multi-cron:latest

```

The command above adds 3 tasks to the multi-cron queue which:

1. Runs `echo` with the arguments `I ran at midnight` upon application start and then every day at midnight UTC

2. Runs a `script` located at `/path/to/some/script` at 8PM UTC with no arguments

3. Runs a `python` script located at `/path/to/python.py` every 3 hours

* _If you prefix the schedule with a `!` character, the command will also run at startup_

* _All commands should exit after running_

---

## Multi-Cron Expression Format

### Standard Expressions

| Field name     | Mandatory? | Allowed values  | Allowed special characters |
| :---           | :---       | :---            | :---                       |
| `Seconds`      | No         | 0-59            | * / , -                    |
| `Minutes`      | Yes        | 0-59            | * / , -                    |
| `Hours`        | Yes        | 0-23            | * / , -                    |
| `Day of month` | Yes        | 1-31            | * / , - ?                  |
| `Month`        | Yes        | 1-12 or JAN-DEC | * / , -                    |
| `Day of week`  | Yes        | 0-6 or SUN-SAT  | * / , - ?                  |

*_Note that there are 6 fields (instead of 5 that cron uses) because as multi-cron resolution goes down to 1 second!_*

***If you prefix the schedule with a `!` character, the command will also run at startup***

### Expression Aliases

| Entry                      | Description                                | Equivalent To   |
| :---                       | :---                                       | :---            |
| `@yearly` (or `@annually`) | Run once a year, midnight, Jan. 1st        | `0 0 0 1 1 *`   |
| `@monthly`                 | Run once a month, midnight, first of month | `0 0 0 1 * *`   |
| `@weekly`                  | Run once a week, midnight on Sunday        | `0 0 0 * * 0`   |
| `@daily` (or `@midnight`)  | Run once a day, midnight                   | `0 0 0 * * *`   |
| `@hourly`                  | Run once an hour, beginning of hour        | `0 0 * * * *`   |

Thus, `@every 1h30m10s` would indicate a schedule that activates every `1 hour, 30 minutes, 10 seconds`.

### Special Characters

`Asterisk` ( `*` )

* The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in the 5th field (month) would indicate every month.

`Slash` ( `/` )

* Slashes are used to describe increments of ranges. For example 3-59/15 in the 1st field (minutes) would indicate the 3rd minute of the hour and every 15 minutes thereafter.
    The form "*\/..." is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field. The form "N/..." is accepted as meaning "N-MAX/...",
    that is, starting at N, use the increment until the end of that specific range. __It does not wrap around__

`Comma` ( `,` )

* Commas are used to separate items of a list. For example, using "MON,WED,FRI" in the 5th field (day of week) would mean Mondays, Wednesdays and Fridays.

`Hyphen` ( `-` )

* Hyphens are used to define ranges. For example, 9-17 would indicate every hour between 9am and 5pm inclusive.

`Question mark` ( `?` )

* Question mark may be used instead of '*' for leaving either day-of-month or day-of-week blank.

### Intervals

You may also schedule a job to execute at fixed intervals. This is supported by formatting the cron spec like this:

`@every <duration>`
Thus, `@every 5s` would run the command every `5 seconds`

---

### Notes

**Time zones**:

* **Be aware that jobs scheduled during daylight-savings leap-ahead transitions will not be run!**
* By default, all interpretation and scheduling is done in the machine's local time zone (as provided by the [Go time package](http://www.golang.org/pkg/time).
    Docker generally runs all contains as `UTC` however it is possible to override that setting by launching the container with an environmental variable (`TZ`) set

**Love2Merge**:

* Pull requests and forks are encouraged and appreciated ;)

---

### Links

* [Dockerhub project](https://hub.docker.com/r/jeffkolb/multi-cron)
* [Github Repo](https://github.com/junkiebev/multi-cron)
* [Full library documentation](https://godoc.org/gopkg.in/robfig/cron.v2)
