# agromeration
Like agglomeration but is agromeration

## Content
1. [ Requirements ](#requirements)
2. [ Local installation ](#local-install)


<a name="requirements"></a>
## Requirements
For development needs you must install following packages:
- [wget](https://www.gnu.org/software/wget/) package.
- [make](https://www.gnu.org/software/make/) package.


<a name="local-install"></a>
## Local installation
After installing wget and make packages yous should install Go.


Downloading source:
```shell
wget install https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
```

Remove any previous Go installation by deleting the /usr/local/go folder (if it exists), then extract the archive you just downloaded into /usr/local, creating a fresh Go tree in /usr/local/go:
```shell
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.23.1.linux-amd64.tar.gz
```

Add /usr/local/go/bin to the PATH environment variable. You can do this by adding the following line to your $HOME/.profile or /etc/profile (for a system-wide installation):
```shell
export PATH=$PATH:/usr/local/go/bin
```


Changes made to a profile file may not apply until the next time you log into your computer. To apply the changes immediately, just run the shell commands directly or execute them from the profile using a command such as:
```shell
source $HOME/.profile.
```

Verify that you've installed Go by opening a command prompt and typing the following command:
```shell
go version
```

Before launching, we should build the project:
```shell
make build-local
```

To launch HTTP-server:
```shell
make up-local
```