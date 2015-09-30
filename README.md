# pupdate
-routine should send back report and controller should retry-done
-pupdate have to log issues to the log file
-support multiple sqls/commands before processing
-support after and multiple sqls/commands processing
-measure performance of single job and report
-store task and job progress into MySQL database (the progress table)
-provide a way of continuing terminated task based on the progress table
-parametrize task storage type: mysql, mongodb, elastic elastic
-parametrize task method: update, select elastic
-parametrize task: add special value (string) to indicate that all rows should be processed instead of providing max(id)


#pupdate
-routine should send back report and controller should retry
-pupdate have to log issues to the log file
-support multiple "sets" for before processing- Tab
-support "set" after processing- Tab
-measure performance of single job
-shutdown after last job finished
-support for signaling
-etlsuccessfull  is server time
-how to reset/reconnect mysql connection?
-load config only flag
-PreSet error max attempt- currently we don't differentiate
-PostStep error max attempt - currently we don't diffirentatie
-Recovery steps (try, catch or deffer, recover)
-change flags -config to -c,-test_confg to -tc etc

#tests
-automated tests for all add and update scripts
-mock database for add and update tests in sqlite
-how to validate test using configuration?

#logging
-logging levels:debug (log), concole, session, syslog
-prefixes for example: executing pre-steps
-when there is an error then show all conext if not show only sid,tid name, start time, stop time, number of rows affected

#delta
if moth is dropped then etl could transfer all data back if data is partitioned but the second phase will be
redundand then because there is no way to find out if mot2.entity was juest added rfr and advice will need additional 2 tasks

Deplyment
Deplpoy app with monitoring/seed_date.dat

*RestDB - ceated subdomain base on server name taken from http.conf
    <VirtualHost *:80>
    ServerName restdb.dat-api06.data.supp.mot.dvsa

ALWAYS do transactions in small batches!! or https://dev.mysql.com/doc/refman/5.5/en/optimizing-innodb-bulk-data-loading.html
