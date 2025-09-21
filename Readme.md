# Product
This project handles social service 

## Development

### Prerequisites
To install development tools run this (sudo is for installing binaries in system bin directory):
```bash
make init
```
And create .env file:
```bash
make env
```

### Get Packages
Create a file named `.netrc` in your home directory (`/home/`) and add the following content, replacing `[username]` and `[access-token...]` with your Git username and access token:

```text
machine git.oceantim.com
login [username]
password [access-token-with-read_repository-and-read_api-access]
```
Run the following command to configure Go environment variables:
```shell
go env -w GOPRIVATE="git.oceantim.com/*"
```

[open swagger](http://localhost:8000/docs/index.html)
