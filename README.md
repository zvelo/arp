golang library for looking up MAC address by IP address. It currently supports
different methods for retrieving the ARP cache from the system depending if running
on Linux, OSX or Windows. 

I've only tested this library on OSX and Linux. 

For Linux, `/proc/net/arp` is used. For OSX (or other unix systems) and Windows,
`exec.Command` is used to call the `arp` utility and parse its output. 

-----

Jonathan Briggs forked this for zvelo.
I need it to return the device information and also information for the
local machine. Some of that is only available on Linux. I tried not to
break the older interfaces.
