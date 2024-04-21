# rules
- reads the file
- calculates the min, mean, and max temperature value per weather station
- emits the results on stdout like this

```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
```

# results
- 1-dumb-solution
```bash
./1brc-go  266,48s user 5,64s system 101% cpu 4:28,38 total
```

- 2-read-chunks
```bash
./1brc-go  227,19s user 5,11s system 101% cpu 3:48,61 total
```

- 3-cache-stations
```bash
./1brc-go  195,46s user 2,65s system 97% cpu 3:22,21 total
```

- 4-one-map
```bash
./1brc-go  159,93s user 2,96s system 95% cpu 2:50,85 total
```

- 5-parallel-map-reduce
```bash
./1brc-go  210,88s user 7,24s system 879% cpu 24,794 total
```

- 6-fast-float-parse
```bash
./1brc-go  96,24s user 5,51s system 504% cpu 20,178 total
```

- 8-memory-mapped-file
```bash
./1brc-go  133,24s user 10,18s system 962% cpu 14,896 total
```
