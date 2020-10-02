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

## Behavior
1. System first checks request from ip address
    - ip address has assigned deviceId -> OK
    - ip address does not have assigned deviceId -> NOK (screenshot)
>![ipaddress_error](screens/no-ip.png)
2. System then checks if workplace has any user logged
    - user is logged -> OK
    - user is not logged -> NOK (screenshot)
>![user_error](screens/no-user.png)
3. Screen for scanning barcode order is displayed (screenshot)
    - if code exists in K2 -> OK
    - if code does not exists in K2 -> NOK (screenshot)
>![order](screens/read-code.png)
>
>![no-code](screens/no-code.png) 
4. After successfully scanning the code, main screen is displayed (screenshot)
>![main](screens/main.png) 
     
© 2020 Petr Jahoda