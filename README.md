# SpotLink
A Winlink service that responds with WSPR spot reports

## Why create this service, and what does it do?

In the internet age, amateur radio operators have rich data at their
disposal for selecting which radio bands have the best propagation 
characteristics. Many folks around the world operate WSPR beacons
and listening stations that build a picture of real-time conditions.
However, when you're in the field without access to the internet
operators don't have access to these resources.

If you have access to [Winlink](https://winlink.org/), this
application makes it possible to query WSPR propagation data. This
makes it possible to:

   * Find the best frequencies to work
   * Find the best directions to direct an antenna with gain
   * Receive data on which Winlink gateways may yield higher-
     quality connections

## How to use it:

Write to `SpotLink@k0jrh.husney.com`. The subject line does not
matter.

In the body of your message, place your commands:

```
set interval="2 hours ago"
byCallsign w0asw
```

```
time               , rx_sign , rx_loc, rx_lat   , rx_lon     , distance, azimuth, frequency      , snr, rx_azimuth
2024-02-05 19:42:00, W0VI    , EN35gc, 45.104000, -93.458000 , 23      , 325    , 3570076.000000 , 1  , 144
2024-02-05 18:50:00, N2HQI   , FN13sa, 43.021000, -76.458000 , 1361    , 93     , 14097055.000000, 14 , 284
2024-02-05 18:50:00, AC0G    , EM38ww, 38.938000, -92.125000 , 674     , 171    , 14097055.000000, 10 , 352
2024-02-05 18:50:00, KV0S    , EM38tv, 38.896000, -92.375000 , 676     , 173    , 14097055.000000, 9  , 353
2024-02-05 18:50:00, K6RFT   , EM47bg, 37.271000, -91.875000 , 861     , 172    , 14097055.000000, 9  , 352
2024-02-05 18:50:00, K9AN    , EN50wc, 40.104000, -88.125000 , 684     , 140    , 14097054.000000, 5  , 323
2024-02-05 18:50:00, N0TOI   , EM48dd, 38.146000, -91.708000 , 767     , 170    , 14097067.000000, 4  , 350
2024-02-05 18:50:00, WA3VPZ  , FM19kn, 39.563000, -77.125000 , 1455    , 109    , 14097052.000000, 4  , 299
2024-02-05 18:50:00, KA7OEI-1, DN31uo, 41.604000, -112.292000, 1578    , 263    , 14097055.000000, 3  , 70
2024-02-05 18:50:00, AJ8S    , EM89bt, 39.813000, -83.875000 , 960     , 123    , 14097055.000000, 2  , 309
```

