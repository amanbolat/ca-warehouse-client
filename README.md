# CrossAsia Warehouse client service

This service is a connector to the main CrossAsia ERP system and only works with FileMaker Database.

### Deployment
Env variables below should be provided before running the service
```
FM_PASS=password
FM_HOST=ip_address
FM_USER=fm_username
FM_DATABASE_NAME=db_name
KDN_BUSINESS_ID=快递鸟business_id
KDN_API_SECRET=快递鸟_secret
PRINTER=printer_name
PORT=server_port
FONT_PATH=font_path
DEBUG=bool
BOLT_DB_PATH=badger_db_path
```

### Installation
1. Download the latest release from this repository
2. Prepare the file with env vars
3. Prepare proxy such as Caddy
4. Prepare daemon description file to run the service on login or just run it:
`whcleint -c <path_to_config_file>` 

Daemon property list example:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>EnvironmentVariables</key>
        <dict>
                <key>PATH</key>
                <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/sbin</string>
        </dict>
        <key>KeepAlive</key>
        <dict>
                <key>SuccessfulExit</key>
                <false/>
        </dict>
        <key>Label</key>
        <string>whsrv</string>
        <key>ProgramArguments</key>
        <array>
                <string>/opt/srv/whclient/whclient</string>
                <string>-c</string>
                <string>/opt/srv/whclient/config.env</string>
                <string>run</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>SessionCreate</key>
        <false/>
        <key>StandardErrorPath</key>
        <string>/var/log/wh.err</string>
        <key>StandardOutPath</key>
        <string>/var/log/wh.out</string>
</dict>
</plist>
```


### Main functions
- List the shipments
- List warehouse entries for
- Create and edit entries

### Schedulers
Every few seconds fetches shipments and checks for new ones with status "preparation". If found any 
prints the preparation information and caches them in order to not print the same records multiple times.

### Web client and authentication
No authentication is required because it should only be run on the local machine with local web client.
ABAC rules are forced on the database side.

### Printing
This service have functionality to create PDF file using mono font and then printing it using 
CUPS printer.