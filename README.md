## About this library

This library is incomplete. I created this as a programming exercise to help me understand how CNI works and how it can be implemented. It currently implements CNI `0.3.1`.

## Example application

Compile with:
~~~
go build
~~~

### Testing manually

Test with the following commands:

#### Version command

~~~
export CNI_COMMAND=VERSION
./easycni
~~~

#### Add command

~~~
export CNI_PATH=/opt/cni/bin/
export CNI_CONTAINERID=cnitool-77383ca0a0715733ca6f
export CNI_IFNAME=eth0
export CNI_COMMAND=DEL
export CNI_NETNS=/var/run/netns/testing
cat 150-easycni.conf | ./easycni
~~~

#### Del command

~~~
export CNI_PATH=/opt/cni/bin/
export CNI_CONTAINERID=cnitool-77383ca0a0715733ca6f
export CNI_IFNAME=eth0
export CNI_COMMAND=DEL
export CNI_NETNS=/var/run/netns/testing
cat 150-easycni.conf | ./easycni
~~~

### Testing with cnitool

You can also test this with [https://github.com/containernetworking/cni/tree/master/cnitool](https://github.com/containernetworking/cni/tree/master/cnitool)
~~~
cp 150-easycni.conf /etc/cni/net.d/150-easycni.conf
cp easycni /opt/cni/bin/
ip netns add testing
CNI_PATH=/opt/cni/bin/ cnitool add easycni /var/run/netns/testing
CNI_PATH=/opt/cni/bin/ cnitool del easycni /var/run/netns/testing
~~~~

### Testing with kubernetes

You can also test this with kubernetes. The example application will do a noop and it will pretend that an interface with a single IP (172.17.0.10) and MAC address (aa:aa:aa:aa:aa:aa) were plugged in.

This is obviously very limited, but it's enough to create a single pod and to see that it receives an IP address.

~~~
cp 150-easycni.conf /etc/cni/net.d/10-easycni.conf
cp easycni /opt/cni/bin/
~~~

Example output:
~~~
[root@node1 easycni]# kubectl get pods -A -o wide | grep coredns
kube-system   coredns-78fcd69978-mzsvh        1/1     Running       0                117s    10.85.0.20        node1   <none>           <none>
kube-system   coredns-78fcd69978-n5khq        1/1     Running       0                3h30m   10.85.0.17        node1   <none>           <none>
[root@node1 easycni]# kubectl delete pod -n kube-system  coredns-78fcd69978-mzsvh
pod "coredns-78fcd69978-mzsvh" deleted
[root@node1 easycni]# kubectl get pods -A -o wide | grep coredns
kube-system   coredns-78fcd69978-bxn4h        0/1     Running       0                9s      172.17.0.10       node1   <none>           <none>
kube-system   coredns-78fcd69978-n5khq        1/1     Running       0                3h30m   10.85.0.17        node1   <none>           <none>
~~~

