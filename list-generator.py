from os import error
import socket

def get_ipv4_by_hostname(hostname):
    # see `man getent` `/ hosts `
    # see `man getaddrinfo`

    return list(
        i        # raw socket structure
            [4]  # internet protocol info
            [0]  # address
        for i in 
        socket.getaddrinfo(
            hostname,
            0  # port, required
        )
        if i[0] is socket.AddressFamily.AF_INET  # ipv4

        # ignore duplicate addresses with other socket types
        and i[1] is socket.SocketKind.SOCK_RAW  
    )

with open('top-1m.csv') as f:
        for lignes in f:
            l = lignes.rstrip().split(",")
            try:
                a = get_ipv4_by_hostname(l[1])
                print(a)
                with open("top-1m-corrected.csv", "a") as f2:
                    f2.write(l[1]+"\n")
            except error:
                pass