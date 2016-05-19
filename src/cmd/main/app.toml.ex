
[web]
Debug = true
Listen = ":1333"
Pprof = ":1777"
Fasthttp = true
StaticDir = "./public"

[AccessLog]
Enable = true
# if FilePath not empty use file output else us os.std
FilePath = "log/access.log"

############## pongo2 ########
[Pongo]
# Directory to load templates. Default is "templates"
Directory = "src/templates"
# Reload to reload templates everytime.
Reload = true


############### models ###############

[models.default]
enable = true
driver = "mysql"
host = "127.0.0.1:3306"
user = "root"
passwd = "sdfsdf"
db = "test"

debug = true
logPath = ""
maxIdle = 100
maxOpen = 100
SSLMode = "disable"


############## redis ###############

[redis.default]
enable = true
addr = "localhost:6379"
passwd = ""
MaxIdle     = 50
IdleTimeout = 50


###############  logs #############


[logs.da.console]
Enable = true
Mode = "console"
Level = "Trace"

[logs.default.console]
Enable = true
Mode = "console"
Level = "Trace"

[logs.default.file]
Enable = true
Mode = "file"
Level = "Info"
#model file
FileName  = "./err.log"
LogRotate  = true
MaxLines   = 1000000
MaxSizeShift = 28
DailyRotate  = true
MaxDays     = 7


[logs.default.smtp]
Enable = false
Mode = "smtp"
Level = "Trace"
# model smtp
User      = ""
Passwd    = ""
Host      = ""
Receivers = []
Subject   = ""