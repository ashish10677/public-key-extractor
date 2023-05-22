# Public Key Extractor

This is a small project that is designed to pull the public key of the sender from a transaction hash.

## How to run

1. Install all dependencies by running `go mod download`.
2. In your environment set the value of INFURA_KEY as the infura project id.
3. Build the project using `go build`.
4. Run the project using `./public-key-extractor --txnHash <transaction hash>`.
