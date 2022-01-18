# Query benchmarking tool
Query benchamrking is a command line tool that can be used to benchmark PromQL query performance across multiple workers/clients against a Promscale instance. The tool it takes as its input a CSV file and a flag to specify the number of concurrent workers. After processing all the queries specified by the parameters in the CSV file, the tool outputs a summary with the following stats:

- Total queries processed.
- Total processing time across all queries.
- The minimum query time (for a single query).
- The median query time.
- The average query time.
- The maximum query time. 

## Setup Promscale instance with sample data

1. Run `run-timescaledb` to install and run TimescaleDB with Promscale extension.


2. Run `make run-promscale` to run the Promscale container.


3. Add sample data:
```
    curl -v \
    -H "Content-Type: application/x-protobuf" \
    -H "Content-Encoding: snappy" \
    -H "X-Prometheus-Remote-Write-Version: 0.1.0" \
    --data-binary "@resources/real-dataset.sz" \
    "http://localhost:9201/write"
   ```

## How to run it?

- Run the service locally with optional `workers` and `timeout` arguments:

`make run workers=3 timeout=2000 // run the tool with 3 worker processes`

**PRO TIP**: Run `make help` to discover awesome commands!

## Considerations
- I decided to implement the `reporter` with a mutex because according to my point of view it was the simplest solution here, and I didn't want to overcomplicate the design. Having said that, the communication between the `scheduler` and `reporter` could also be implemented through a channel.
