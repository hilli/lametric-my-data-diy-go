# Go types and http handler for LaMetric My Data DIY

[LaMetric](https://lametric.com/) makes an [app](https://apps.lametric.com/apps/my_data_diy__with_no-code_possibilities_/8942?apps_for=sky&mosaic=129) for easy access to their Time or Sky panels called `My Data DIY`. To make the best of that, they (LeMetric) has made a [data format](https://help.lametric.com/support/solutions/articles/6000225467-my-data-diy) for it. This is supported by this library.

## Usage

- Create some "Frames"
- Add the handler to your webserver.
- Point the LaMetric device to your services as a HTTP pull endpoint
- Bob's your uncle (That is: Enjoy your data on the LaMetric display)

For a code example see the included webserver in [cmd/web/main.go](https://github.com/hilli/lametric-my-data-diy-go/blob/main/cmd/web/main.go)

## Running example

```shell
go run cmd/web/main.go
# or
task run
```

## Notes

This repo uses a [taskfile](https://taskfile.dev/) to simplify things. [Install task](https://taskfile.dev/installation/) if you haven't already and run this to list available targets:

```shell
$ task --list
task: Available tasks for this project:
* build:               Build the binary
* cleanup:             Cleanup the project
* coverage:            Generate coverage report and print it to the console
* default:             Default task; Run tests, monitoring continuously for changes
* html-coverage:       Generate HTML coverage report and open it in the browser
* run:                 Run the binary
* test:                Run all tests
```