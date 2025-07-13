@echo off
rem run this script as admin

if not exist catcher.exe (
    echo "file not found"
    goto :exit
)

sc create Catcher binpath= "%CD%\catcher.exe -config=config/config.yml" start= auto DisplayName= "Catcher"
net start Catcher
sc query Catcher


:exit
