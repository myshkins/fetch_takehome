### Overview
This program satisfies the specifications laid out in the "Fetch Take-Home Exercise — Site
Reliability Engineering" document. 

### How to Run
Copy the source code to your local machine
```
git clone git@github.com:myshkins/fetch_takehome.git

cd fetch_takehome
```
If you have go installed on your machine, you can follow the steps in the next section. Alternatively, you may follow the steps in [With Docker](###with-docker).

##### With Go installed
Change to the `cmd/health_check` directory and then build the binary
```
cd cmd/health_check

go build -C cmd/health_check
```

You can then run the program via:
```
./health_check --config-file=path/to/config/yaml --interval=15
```
Both flags are optional. The `-config-file` flag defaults to the `sample_input.yaml` provided in the repo, and the `-interval` flag defaults to 15 seconds.

##### With Docker
To run the program with docker, first build the image. From the root of the `fetch_takehome`  repo, run:
```
docker build -t health_check
```

You can then run the program in a docker container via:
```
docker run health_check "--config-file=path/to/config/yaml" "--interval=15"
```
Again both flags are optional.

