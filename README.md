# dns-tester

Project that allows to test how fast a dns server responds to multiple queries from multiple clients

## How it works

The tester works by creating a certain number of threads (also called clients) and giving them a list of domains to resolve. Each client have a unique list of domains (without intersections with the lists of the other clients) and must resolve every domain name in it. The number of requests per clients is the total number of requests divided by the number of clients.
At the end of the program, the tester prints stats about the time needed by each client to resolve its entire list of domains.

## Build instructions

- Install the latest version of golang and make sure it's accessible through path
- Clone project
- cd into the repo directory
- Type ```go mod tidy```
- Type ```go build```

## Usage instructions

- Build the project (cf. Build Instruction)
- Type ```./dns-tester -dns <dns IP> -nr <total number of requests> -nc <number of different clients>```

## Domain list generation

The file ```data/top-1m.csv``` contains an outdated list of the top 1 millions domains requested. Some domains in this list do not exists anymore so we must verify them before using it.

The domain list used by the tester is ```data/top-1m-corrected.csv``` and can be generated using the script ```data/list-generator.py``` that tests all domains in ```data/top-1m.csv``` and put existing domains in ```data/top-1m-corrected.csv```