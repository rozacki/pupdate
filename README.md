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
-
