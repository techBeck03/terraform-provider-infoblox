# go-ipmath

## Simple IP math functions for golang<br/><br/>

### Base object<br/><br/>
```
type IP struct {
	Address net.IP
	Network *net.IPNet
}
```

### Supported Operations<br/><br/>

- `Inc`: Increments IP address by 1
- `Dec`: Decrements IP address by 1
- `Add`: Increments IP address by provided increment
- `Subtract`: Decrements IP address by provided increment
- `Difference`: Returns the signed int difference of object IP and provided IP
- `EQ`: Checks if base object IP and provided IP are equal
- `GT`: Checks if base object IP is greater than provided IP
- `LT`: Checks if base object IP is less than provided IP
- `GTE`: Checks if base object IP is greater than or equal to provided IP
- `LTE`: Checks if base object IP is less than or equal to provided IP

## Example Usage

### Addition

```golang
ip, err := ipmath.NewIP("172.19.4.10/23")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Add(256)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
172.19.5.10
```

### Addition with Errors

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Add(256)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
2021/03/04 15:52:22 172.19.5.10 is not in CIDR network 172.19.4.0/24
```

### Subtraction

```golang
ip, err := ipmath.NewIP("172.19.4.10/23")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Subtract(4)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
172.19.4.6
```

### Subtraction with Errors

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Subtract(30)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
2021/03/04 15:57:46 172.19.3.236 is not in CIDR network 172.19.4.0/24
```

### Difference

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.Difference(net.ParseIP("172.19.4.5")))
fmt.Println(ip.Difference(net.ParseIP("172.19.4.25")))
```

#### Output

```bash
-5
15
```

### Increment

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Inc()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
172.19.4.11
```

### Decrement

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
newIP, err := ip.Dec()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(newIP.String())
```

#### Output

```bash
172.19.4.9
```

### EQ Comparison

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.EQ(net.ParseIP("172.19.4.10")))
fmt.Println(ip.EQ(net.ParseIP("172.19.4.11")))
```

#### Output

```bash
true
false
```

### GT Comparison

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.GT(net.ParseIP("172.19.4.10")))
fmt.Println(ip.GT(net.ParseIP("172.19.4.9")))
fmt.Println(ip.GT(net.ParseIP("172.19.4.12")))
```

#### Output

```bash
false
true
false
```

### LT Comparison

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.LT(net.ParseIP("172.19.4.10")))
fmt.Println(ip.LT(net.ParseIP("172.19.4.9")))
fmt.Println(ip.LT(net.ParseIP("172.19.4.12")))
```

#### Output

```bash
false
false
true
```

### GTE Comparison

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.GTE(net.ParseIP("172.19.4.10")))
fmt.Println(ip.GTE(net.ParseIP("172.19.4.9")))
fmt.Println(ip.GTE(net.ParseIP("172.19.4.12")))
```

#### Output

```bash
true
true
false
```

### LTE Comparison

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.LTE(net.ParseIP("172.19.4.10")))
fmt.Println(ip.LTE(net.ParseIP("172.19.4.9")))
fmt.Println(ip.LTE(net.ParseIP("172.19.4.12")))
```

#### Output

```bash
true
false
true
```
