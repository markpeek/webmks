Test to connect to a vSphere webmks using VNC over websockets
======

To run, set these environment variables to your settings
```
export GOVC_URL="10.1.10.51"
export GOVC_USERNAME='administrator@vsphere.local'
export GOVC_PASSWORD='vmware'
export GOVC_INSECURE=1
```

Build the binary
```
$ go build
```

Open up a web console to view the test and then run it with a path to the VM:

```
$ ./webmks /datacenter/vm/photon-3.0
Found machine: photon-3.0
sending test command to console
done
```
