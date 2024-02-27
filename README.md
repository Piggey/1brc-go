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
