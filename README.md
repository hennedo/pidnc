# PiDNC
[![Build Status](https://drone.dre.li/api/badges/henne/pidnc/status.svg)](https://drone.dre.li/henne/pidnc)

[Trello Roadmap](https://trello.com/b/Z3vMX6jr/dnc-shit)

## Configuration
Konfiguration erfolgt aktuell über environment Variablen.

| Var          | Default | Required | Beschreibung                                  |
| ------------ | ------- | -------- | --------------------------------------------- |
| HOST         |         |          | Host auf den das Webinterface binded          | 
| PORT         | 8000    |          | Port auf dem das Webinterface binded          |
| GCODE_FOLDER |         | x        | Ordner wo die gcodes landen. Muss existieren! |
| SERIAL_PORT  |         | x        | Pfad zum Serial Port. Muss existieren!        |