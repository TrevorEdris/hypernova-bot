# Performance Benchmarks using Local Development

## Local Storage

### Create record

```
❯ hey -z 1m -H "Content-Type: application/json" -m POST -d '{"name":"Item_HeyTest","description":"lol desc","price":1.337}' http://localhost:8000/item

Summary:
  Total:        60.0020 secs
  Slowest:      0.0301 secs
  Fastest:      0.0003 secs
  Average:      0.0030 secs
  Requests/sec: 27101.4270


Response time histogram:
  0.000 [1]     |
  0.003 [937950]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.006 [57346] |■■
  0.009 [4137]  |
  0.012 [396]   |
  0.015 [100]   |
  0.018 [32]    |
  0.021 [14]    |
  0.024 [10]    |
  0.027 [7]     |
  0.030 [7]     |


Latency distribution:
  10% in 0.0010 secs
  25% in 0.0013 secs
  50% in 0.0017 secs
  75% in 0.0021 secs
  90% in 0.0029 secs
  95% in 0.0035 secs
  99% in 0.0051 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.0301 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0002 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0081 secs
  resp wait:    0.0027 secs, 0.0003 secs, 0.0261 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0168 secs

Status code distribution:
  [201] 1000000 responses
```

### Get record

```
❯ hey -z 1m -H "Content-Type: application/json" http://localhost:8000/item/some_unique_id

Summary:
  Total:        60.0029 secs
  Slowest:      0.0231 secs
  Fastest:      0.0003 secs
  Average:      0.0030 secs
  Requests/sec: 28354.4751


Response time histogram:
  0.000 [1]     |
  0.003 [881439]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [108757]        |■■■■■
  0.007 [7725]  |
  0.009 [1728]  |
  0.012 [248]   |
  0.014 [51]    |
  0.016 [25]    |
  0.019 [5]     |
  0.021 [16]    |
  0.023 [5]     |


Latency distribution:
  10% in 0.0010 secs
  25% in 0.0012 secs
  50% in 0.0016 secs
  75% in 0.0020 secs
  90% in 0.0027 secs
  95% in 0.0034 secs
  99% in 0.0048 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.0231 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0123 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0108 secs
  resp wait:    0.0026 secs, 0.0002 secs, 0.0230 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0150 secs

Status code distribution:
  [200] 1000000 responses
```

### Update record

```
❯ hey -z 1m -H "Content-Type: application/json" -m PUT -d '{"price":4.2069}' http://localhost:8000/item/some_unique_id

Summary:
  Total:        60.0017 secs
  Slowest:      0.0261 secs
  Fastest:      0.0003 secs
  Average:      0.0030 secs
  Requests/sec: 27336.5390


Response time histogram:
  0.000 [1]     |
  0.003 [907924]        |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.005 [85663] |■■■■
  0.008 [5224]  |
  0.011 [944]   |
  0.013 [119]   |
  0.016 [70]    |
  0.018 [41]    |
  0.021 [8]     |
  0.024 [2]     |
  0.026 [4]     |


Latency distribution:
  10% in 0.0010 secs
  25% in 0.0013 secs
  50% in 0.0016 secs
  75% in 0.0021 secs
  90% in 0.0028 secs
  95% in 0.0035 secs
  99% in 0.0049 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0003 secs, 0.0261 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0019 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0112 secs
  resp wait:    0.0027 secs, 0.0002 secs, 0.0260 secs
  resp read:    0.0002 secs, 0.0000 secs, 0.0143 secs

Status code distribution:
  [201] 1000000 responses
```

## DynamoDB Storage

Obviously, incorporating network communications into the mix is going to reduce performance. Even so,
the one replica of this API was able to support between 100 and 200 RPS on a Ryzen 9 3900x.

### Create record

```
❯ hey -z 1m -H "Content-Type: application/json" -m POST -d '{"name":"Item_HeyTest","description":"lol desc","price":1.337}' http://localhost:8000/item

Summary:
  Total:        60.0300 secs
  Slowest:      0.4034 secs
  Fastest:      0.0402 secs
  Average:      0.2363 secs
  Requests/sec: 211.5607


Response time histogram:
  0.040 [1]     |
  0.077 [2]     |
  0.113 [0]     |
  0.149 [0]     |
  0.186 [0]     |
  0.222 [2458]  |■■■■■■■■■■■■
  0.258 [8392]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.294 [1761]  |■■■■■■■■
  0.331 [84]    |
  0.367 [0]     |
  0.403 [2]     |


Latency distribution:
  10% in 0.2184 secs
  25% in 0.2231 secs
  50% in 0.2301 secs
  75% in 0.2436 secs
  90% in 0.2689 secs
  95% in 0.2779 secs
  99% in 0.2899 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.0402 secs, 0.4034 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0025 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0012 secs
  resp wait:    0.2361 secs, 0.0372 secs, 0.4033 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0011 secs

Status code distribution:
  [201] 12700 responses
```

### Get record

```
❯ hey -z 1m -H "Content-Type: application/json" http://localhost:8000/item/some_unique_key

Summary:
  Total:        60.0976 secs
  Slowest:      0.4204 secs
  Fastest:      0.1006 secs
  Average:      0.2263 secs
  Requests/sec: 220.7075


Response time histogram:
  0.101 [1]     |
  0.133 [21]    |
  0.165 [19]    |
  0.197 [33]    |
  0.228 [9565]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.260 [2634]  |■■■■■■■■■■■
  0.292 [976]   |■■■■
  0.324 [12]    |
  0.356 [0]     |
  0.388 [0]     |
  0.420 [3]     |


Latency distribution:
  10% in 0.2124 secs
  25% in 0.2157 secs
  50% in 0.2211 secs
  75% in 0.2296 secs
  90% in 0.2569 secs
  95% in 0.2643 secs
  99% in 0.2795 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.1006 secs, 0.4204 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0049 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0053 secs
  resp wait:    0.2262 secs, 0.1004 secs, 0.4202 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0009 secs

Status code distribution:
  [200] 13264 responses
```

### Update record

```
❯ hey -z 1m -H "Content-Type: application/json" -m PUT -d '{"price":4.2069}' http://localhost:8000/item/some_unique_key

Summary:
  Total:        60.2631 secs
  Slowest:      0.6258 secs
  Fastest:      0.4170 secs
  Average:      0.4600 secs
  Requests/sec: 108.6237


Response time histogram:
  0.417 [1]     |
  0.438 [866]   |■■■■■■■■■■■
  0.459 [3160]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.480 [979]   |■■■■■■■■■■■■
  0.501 [1114]  |■■■■■■■■■■■■■■
  0.521 [374]   |■■■■■
  0.542 [48]    |■
  0.563 [0]     |
  0.584 [1]     |
  0.605 [0]     |
  0.626 [3]     |


Latency distribution:
  10% in 0.4360 secs
  25% in 0.4431 secs
  50% in 0.4534 secs
  75% in 0.4776 secs
  90% in 0.4945 secs
  95% in 0.5031 secs
  99% in 0.5193 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0000 secs, 0.4170 secs, 0.6258 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0005 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0024 secs
  resp wait:    0.4598 secs, 0.4169 secs, 0.6257 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0010 secs

Status code distribution:
  [201] 6546 responses
```
