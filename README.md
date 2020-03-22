# vulpes
[![CircleCI](https://circleci.com/gh/terakoya76/vulpes.svg?style=svg)](https://circleci.com/gh/terakoya76/vulpes)

mysql metrics reporter.

vulpes simply exports mysql status information to stdout.

- SHOW GLOBAL STATUS
- SHOW GLOBAL VARIABLES
- SHOW ENGINE INNODB STATUS

vuples is not implemented as a long-lived daemon process, but also as a simple oneshot job process.
You can use it as a part of your monitoring pipeline input by providing a scheduled execution setting outside of vulpes.

## Configuration
### Required Environment Variables
```bash
export DB_DRIVER=mysql
export DB_HOSTNAME=127.0.0.1
export DB_PORT=3306
export DB_USERNAME=mysql
export DB_PASSWORD=root
export DB_NAME=testdb
```
