# Simple CLI client for jidelna.cz

## Build
1. Clone this repository using
```
git clone https://github.com/matejp0/jidelna
```
2. Go to the cloned repository
```
cd jidelna
```
3. Add a `vars.go` file with the following contents:
```
package main

const EMAIL = "<your email>"
const PASSWORD = "<your password>"
```
4. Build it
```
go build .
```
5. Install it
```
go install .
```

## Usage
You can get the usage by running
```
jidelna --help
jidelna list --help
jidelna order --help
```

### Example
```
jidelna list
jidelna order -d 2022-10-20 18015
jidelna info
```
