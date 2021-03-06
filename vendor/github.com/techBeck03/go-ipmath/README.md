# go-ipmath

## Simple IP math functions for golang<br/><br/>

### Base object<br/><br/>
```golang
type IP struct {
	Address net.IP
	Network *net.IPNet
}
```

### Creating base object<br/><br/>

#### New from CIDR string

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
```

#### Existing IP and Network object

```golang
ip, net, _ := net.ParseCIDR("172.19.4.10/24")

ipObj := ipmath.IP{
    Address: ip,
    Network: net,
}
```

#### Existing IP object only

```golang
net := net.ParseIP("172.19.4.10")

ipObj := ipmath.IP{
    Address: ip,
}
```


### Supported Operations<br/><br/>

- `Add`: Increments IP address by provided increment
- `Subtract`: Decrements IP address by provided increment
- `Inc`: Increments IP address by 1
- `Dec`: Decrements IP address by 1
- `Difference`: Returns the signed int difference of object IP and provided IP
- `EQ`: Checks if base object IP and provided IP are equal
- `GT`: Checks if base object IP is greater than provided IP
- `LT`: Checks if base object IP is less than provided IP
- `GTE`: Checks if base object IP is greater than or equal to provided IP
- `LTE`: Checks if base object IP is less than or equal to provided IP
- `Clone`: Clones the base object into a new base object

## Detailed Examples

### Addition

```golang
ip, err := ipmath.NewIP("172.19.4.10/23")
if err != nil {
    log.Println(err)
    return
}
err = ip.Add(256)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.5.10
172.19.5.10/23
```

### Addition (without Network)

```golang
ip := ipmath.IP{
    Address: net.ParseIP("172.19.4.10"),
}
err := ip.Add(256)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.5.10
Unable to create cidr string because `Network` is undefined
```

### Addition with Errors (Requires `Network` be defined)

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
err = ip.Add(256)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
```

#### Output

```bash
2021/03/04 15:52:22 172.19.5.10 is not in CIDR network 172.19.4.0/24
```

### Subtraction

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
err = ip.Subtract(4)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.4.6
172.19.4.6/24
```

### Subtraction (without Network)

```golang
ip := ipmath.IP{
    Address: net.ParseIP("172.19.4.10"),
}
err := ip.Subtract(4)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.4.6
Unable to create cidr string because `Network` is undefined
```

### Subtraction with Errors (Requires `Network` be defined)

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
err = ip.Subtract(30)
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
```

#### Output

```bash
2021/03/04 15:57:46 172.19.3.236 is not in CIDR network 172.19.4.0/24
```

### Increment

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
err = ip.Inc()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.4.11
172.19.4.11/24
```

### Increment (without Network)

```golang
ip := ipmath.IP{
    Address: net.ParseIP("172.19.4.10"),
}
err := ip.Inc()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.4.11
Unable to create cidr string because `Network` is undefined
```

### Decrement

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
err = ip.Dec()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

### Decrement (without Network)

```golang
ip := ipmath.IP{
    Address: net.ParseIP("172.19.4.10"),
}
err := ip.Dec()
if err != nil {
    log.Println(err)
    return
}
fmt.Println(ip.ToIPString())
cidrString, err := ip.ToCIDRString()
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(cidrString)
```

#### Output

```bash
172.19.4.9
Unable to create cidr string because `Network` is undefined
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

### Clone

```golang
ip, err := ipmath.NewIP("172.19.4.10/24")
if err != nil {
    log.Println(err)
    return
}
clonedIP := ip.Clone()
err = clonedIP.Inc()
if err != nil {
    log.Println(err)
    return
}
fmt.Printf("ip       = %s\n", ip.ToIPString())
fmt.Printf("clonedIP = %s", clonedIP.ToIPString())
```

#### Output

```bash
ip       = 172.19.4.10
clonedIP = 172.19.4.11
```