# CHANGELOG

## Untagged

 - Set conn.WriteTimeout to 10s and failover to log package if we fail to write it out.
 - Wait a maximum of 10s for any locks before bailing and printing to stdout
 - Add a `Flush` method to the logger so we can wait for all messages to be sent
 - Add a timestamp to logs that timeout and get written to stdout
 - Fix `tcp write i/o timeout` causing an infinite loop in `writeToLogEntries`
 - Add concurrentWrites param to limit the maximum number of goroutines active logging at a time
 - Add `errOutput` param to allow users to direct errors