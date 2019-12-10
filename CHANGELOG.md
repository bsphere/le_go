# CHANGELOG

## Untagged

 - Set conn.WriteTimeout to 10s and failover to log package if we fail to write it out.
 - Wait a maximum of 10s for any locks before bailing and printing to stdout