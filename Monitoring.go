//Monitoring is not existingt interface yet.
//If exists it will monitor session, tasks, jobs and persisting their state in strctured form so that it can be restored and act upon in case of failure.
package main

import (
	_"sort"

)

// Depending on configuration we can support db monitoring.
// Currently we suport only log-based monitoring
type MonitoringConfiguration struct{
	Dsn string
	Debug bool
}