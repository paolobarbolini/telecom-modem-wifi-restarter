## Telecom Italia Modem Wifi Restarter [![MIT License](https://img.shields.io/github/license/paolobarbolini/telecom-modem-wifi-restarter.svg?maxAge=86400)](LICENSE) [![Build Status](https://travis-ci.org/paolobarbolini/telecom-modem-wifi-restarter.svg?branch=master)](https://travis-ci.org/paolobarbolini/telecom-modem-wifi-restarter)
_Automatically restart the wifi on your telecom italia modem!_

This program was created after having to frequently restart the wifi on my telecom italia modem to make it work again.

It can also be used to automatically rotate the wifi password on your modem.

You can download the precompiled binaries for windows, mac, linux and arm from the github [releases page](https://github.com/paolobarbolini/telecom-modem-wifi-restarter/releases).

### Usage
After downloading it, execute the following command from the terminal to restart the wifi:

```
$ /path/to/executable --modem-password MODEM_PASSWORD --wifi-password THE_WIFI_PASSWORD
```

If the ip of your modem isn't ``192.168.1.1`` you must also indicate it by adding ``--url "http://your.modem.ip.address"`` flag.

If you want to periodically restart your wifi you can create a cron job

```
$ crontab -e
0 2 * * * /path/to/executable --modem-password MODEM_PASSWORD --wifi-password THE_WIFI_PASSWORD
```

which will run every night at 2am.

### Limitations
The device running this program must be connected to the modem via ethernet, as it wouldn't be able to send an enable request to the wifi after being disconnected.
