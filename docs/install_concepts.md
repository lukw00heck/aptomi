# Aptomi Install / Concepts Mode
* Aptomi will be installed on a local machine and pre-populated with an example
* The part responsible for application deployment **will be disabled**. It means you can explore Aptomi without configuring a k8s cluster, but you won't be able to deploy any apps either

# Installation

## Option #1: docker container
You can run Aptomi **server**, pre-populated with an example, in a Docker container: 
```bash
docker run -it --rm -p 27866:27866 aptomi/aptomi-test-install:xenial sh -c 'curl https://raw.githubusercontent.com/Aptomi/aptomi/master/scripts/aptomi_install.sh | bash /dev/stdin --with-example && aptomi server'
```

And install **client** locally (optional):
```
curl https://raw.githubusercontent.com/Aptomi/aptomi/master/scripts/aptomi_install.sh | bash /dev/stdin --client-only
```

## Option #2: local binaries
Alternatively, you can install Aptomi, pre-populated with an example, on a local machine:
```bash
curl https://raw.githubusercontent.com/Aptomi/aptomi/master/scripts/aptomi_install.sh | bash /dev/stdin --with-example && aptomi server
```

It will install:
* Aptomi binaries in `/usr/local/bin/`
* Aptomi server config in `/etc/aptomi/` and use `/var/lib/aptomi` as persistent data store
* Aptomi client config and examples in `~/.aptomi/`

# Accessing UI
Open UI at [http://localhost:27866/](http://localhost:27866/) and log in as **'admin/admin'**.

Example is already loaded into Aptomi, so you don't need to load it separately. 

# Useful Commands

## Cleaning up local installation
To delete all Aptomi binaries installed locally as well as Aptomi data, run:
```
curl https://raw.githubusercontent.com/Aptomi/aptomi/master/scripts/aptomi_uninstall_and_clean.sh | bash
```
