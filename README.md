# BBQueue

BBQueue is an [open hardware](https://github.com/minor-industries/hardware/tree/main/bbqueue), open source wireless
temperature probe device for cooking.

Features:

- Low-power radio (RFM69): This can operate for weeks on battery power and has much better range than WIFI.
- Web-based dashboards: temperature trends are viewable in any web browser.
- High accuracy: Uses commercial temperature probes (e.g. [thermoworks](https://www.thermoworks.com/shop/products/probes/pro-series/)) and an ADS1115 analog-to-digial converter.
- Two temperature channels are available to monitor both food and bbq temperature.
- Battery level monitoring.

## Photos

![outside shot](https://minor-industries.sfo2.digitaloceanspaces.com/sw/bbqueue_outside.jpg)

## Screenshot

![screenshot](https://minor-industries.sfo2.digitaloceanspaces.com/sw/bbqueue_screenshot_01.png)

## Hardware

![bbqueue board](https://minor-industries.sfo2.digitaloceanspaces.com/hw/bbqueue_board.jpg)

![bbqueue with cables](https://minor-industries.sfo2.digitaloceanspaces.com/hw/bbqueue_with_cables.jpg)

Hardware schematics are available [here](https://github.com/minor-industries/hardware/tree/main/bbqueue).

## Known issues

- The analog to digital converter tends to pick up noise when the power supply is plugged in. Better to run
  on batteries until I can figure this out and create a new hardware revision. 