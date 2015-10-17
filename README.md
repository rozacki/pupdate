#pupdate

-constant measure performance of single job
-shutdown after last job finished
-support for linux signaling
-etl successfull  is server time
-how to reset/reconnect mysql connection?

-PreSet error max attempt- currently we don't differentiate what really caused error
-PostStep error max attempt - currently we don't diffirentatie what really caused error
-Recovery steps (try, catch or deffer, recover)
-change flags -config to -c,-test_confg to -tc etc
-from command line output whole configuration with comments attached
-parametrize task method: exec, query, querysingle


-parametrize task: add special value (string) to indicate that all rows should be processed instead of providing max(id)
-SQL should be using SQL parameters
-SessionController.startTask  defer should handle panic(TaskData) instead of using methods global TaskData..?
-XYZConfiguration should be read-only interface that has accerss to SQL methods
-LOAD FILES-parallel
-Partition database
-extend sql to hanlde files: use nfs
-tableau interface
-rest interface
-sheel interface

#tests
-automated tests for all add and update scripts
-mock database for add and update tests in sqlite
-how to validate test using configuration?

#logging
-logging levels: syslog
-when there is an error then show all conext if not show only sid,tid name, start time, stop time, number of rows affected

#delta
if destination database is dropped then etl could transfer all data back if data is partitioned but the second phase will be
redundand then because there is no way to find out if entity was juest added rfr and advice will need additional 2 tasks

Deplyment
Deplpoy app with monitoring/seed_date.dat


ALWAYS do transactions in small batches!! or https://dev.mysql.com/doc/refman/5.5/en/optimizing-innodb-bulk-data-loading.html
