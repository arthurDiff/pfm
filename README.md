# Port Forward Me

This is small cli tool to start port forwarding and open up provided port for both router and firewall.

Currently only supports TCP and UDP for protocols and iptables for firewall.

This also requires UPNP is on for your router.

NOTE: Due to firewall access this needs to be ran with `sudo`

### Getting Started

#### Install:

```bash
go install github.com/artrctx/pfm
```

#### Run:

```bash
pfm --port 25565 --protocol tcp --firewall iptables
```

#### Example Output:

```
2025/12/07 19:43:18 Initializing firewall client for iptables
2025/12/07 19:43:18 Adding firewall ruleset for protocol tcp | port 8080
2025/12/07 19:43:18 Initializing UPNP
2025/12/07 19:43:20 Start UPnP Port Mapping
2025/12/07 19:43:20 Port Forwarding Mapped Successfully!
Extern: XXX.XXX.X.XX:8080 | Local: 192.168.4.24:8080
Press q to quit
```

You can verify config is working by using `https://portchecker.co/` or pinging your extern addr.
