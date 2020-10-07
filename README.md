[![developed_using](https://img.shields.io/badge/developed%20using-Jetbrains%20Goland-lightgrey)](https://www.jetbrains.com/go/)
<br/>
![GitHub](https://img.shields.io/github/license/petrjahoda/sklabel_cutting_webservice)
[![GitHub last commit](https://img.shields.io/github/last-commit/petrjahoda/sklabel_cutting_webservice)](https://github.com/petrjahoda/sklabel_cutting_webservice/commits/master)
[![GitHub issues](https://img.shields.io/github/issues/petrjahoda/sklabel_cutting_webservice)](https://github.com/petrjahoda/sklabel_cutting_webservice/issues)
<br/>
![GitHub language count](https://img.shields.io/github/languages/count/petrjahoda/sklabel_cutting_webservice)
![GitHub top language](https://img.shields.io/github/languages/top/petrjahoda/sklabel_cutting_webservice)
![GitHub repo size](https://img.shields.io/github/repo-size/petrjahoda/sklabel_cutting_webservice)
<br/>
[![Docker Pulls](https://img.shields.io/docker/pulls/petrjahoda/sklabel_cutting_webservice)](https://hub.docker.com/r/petrjahoda/sklabel_cutting_webservice)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/petrjahoda/sklabel_cutting_webservice?sort=date)](https://hub.docker.com/r/petrjahoda/sklabel_cutting_webservice/tags)
<br/>
[![developed_using](https://img.shields.io/badge/database-MySQL-red)](https://www.mysql.com) [![developed_using](https://img.shields.io/badge/runtime-Docker-red)](https://www.docker.com)

# SK Label Cutting Webservice

## Description
Go web service for SK Label Cutting Workplaces

## Initial screens
1. System first checks request from ip address
    - ip address has assigned deviceId in ```device``` table -> OK
    - ip address does not have assigned deviceId in ```device``` table -> NOK (screenshot)
>![ipaddress_error](screens/no-ip.png)
2. System then checks if workplace has any user ```logged in terminal_input_login``` table
    - user is logged -> OK
    - user is not logged -> NOK (screenshot)
>![user_error](screens/no-user.png)
3. Screen for scanning barcode order is displayed (screenshot)
    - if code exists in K2 -> OK
    - if code does not exists in K2 -> NOK (screenshot)
>![order](screens/read-code.png)
>
>![no-code](screens/no-code.png) 
4. After successfully scanning the code, main screen is displayed (screenshot) with 4 buttons
    - idle button
    - end order button
    - user change button
    - user break button
    
>![main](screens/main.png)
 
## User change button
1. Screen for scanning rfid is displayed (screenshot)
    - if rfid code exists
        - current order is closed
        - K105 is saved to K2
        - new order with new user is created
        - home screen is updated with new user (screenshot)
    - if rfid code does noe exist, user is informed on screen (screenshot)
    - if button for going back is displayed, user gets back to home screen
    
>![order](screens/user_change.png)
>
>![home-new](screens/home_new_user.png)
>
>![user-error](screens/user_error.png)

## User break button
 1. Initial is processed
    - current order is closed
    - K219 is saved to K2
    - 0004 is saved to K2
    - new order with NO user is created
 2. Screen for scanning rfid is displayed (screenshot)
     - if rfid code exists
         - current order is closed
         - new order with new user is created
         - home screen is updated with new user
     - if rfid code does not exist, user is informed on screen
     
 >![order](screens/user_break.png)

## Idle button
1. Screen(s) for choosing idle is selected 
2. After idle is selected
    - K219 is saved to K2
    - idle code is saved to K2
    - idle is created in Zapsi
    - new screen is displayed
3. After end button is selected
    - K119 is saved to K2
    - idle is closed in Zapsi
    - home screen is displayed again 
>![idleSelect](screens/idle-select.png) 
>
>![idleRunning](screens/idle-running.png) 

## Cutting end button
1. Screen(s) for adding number of pcs is displayed 
2. After ok button is selected
    - K302 is saved to K2
    - proper amount is saved to K2
    - order is closed in Zapsi
    - original link is opened in format ```http://localhost:81/terminal/www/default/{DeviceId}```
>![cuttingEnd](screens/cutting-end.png)   
    
© 2020 Petr Jahoda