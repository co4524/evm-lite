# EVM-LITE
## A lean Ethereum node with interchangeable consensus.

We took the [Go-Ethereum](https://github.com/ethereum/go-ethereum)
implementation (Geth) and extracted the EVM and Trie components to create a
lean and modular version with interchangeable consensus.

The EVM is a virtual machine specifically designed to run untrusted code on a
network of computers. Every transaction applied to the EVM modifies the State
which is persisted in a Merkle Patricia tree. This data structure allows to
simply check if a given transaction was actually applied to the VM and can
reduce the entire State to a single hash (merkle root) rather analogous to a
fingerprint.

The EVM is meant to be used in conjunction with a system that broadcasts
transactions across network participants and ensures that everyone executes the
same transactions in the same order. Ethereum uses a Blockchain and a Proof of
Work consensus algorithm. EVM-Lite makes it easy to use any consensus system,
including [Babble](https://github.com/mosaicnetworks/babble) .

## ARCHITECTURE

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State       |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  | -Keystore   |         | -Trie       |  |
                |  |             |         | -Database   |  |
                |  +-------------+         +-------------+  |
                |         |                       ^         |     
                |         v                       |         |
                |  +-------------------------------------+  |
                |  | Engine                              |  |
                |  |                                     |  |
                |  |       +----------------------+      |  |
                |  |       | Consensus            |      |  |
                |  |       +----------------------+      |  |
                |  |                                     |  |
                |  +-------------------------------------+  |
                |                                           |
                +-------------------------------------------+

```

## Consensus Implementations:

- **SOLO**: No Consensus. Transactions are relayed directly from Service to
  State.

- **[BABBLE](https://github.com/mosaicnetworks/babble)**: Inmemory Babble node.
  EVM-Lite does not support Babble's FastSync and Dynamic Membership protocols
  yet, so it is important to set the `--store` flag, and a high `--sync-limit` 
  value. 

- **[RAFT](https://github.com/hashicorp/raft)**: Hashicorp implementation of
  Raft (limited).

more to come...

## USAGE

Each consensus has its own subcommand `evml [consensus]`, and its own
configuration flags.

```
EVM-Lite node

Usage:
  evml [command]

Available Commands:
  babble      Run the evm-lite node with Babble consensus
  help        Help about any command
  raft        Run the evm-lite node with Raft consensus
  solo        Run the evm-lite node with Solo consensus (no consensus)
  version     Show version info

Flags:
  -d, --datadir string        Top-level directory for configuration and data (default "/home/user/.evm-lite")
      --eth.cache int         Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string         Eth database file (default "/home/user/.evm-lite/eth/chaindata")
      --eth.genesis string    Location of genesis file (default "/home/user/.evm-lite/eth/genesis.json")
      --eth.keystore string   Location of Ethereum account keys (default "/home/user/.evm-lite/eth/keystore")
      --eth.listen string     Address of HTTP API service (default ":8080")
      --eth.pwd string        Password file to unlock accounts (default "/home/user/.evm-lite/eth/pwd.txt")
  -h, --help                  help for evml
      --log string            debug, info, warn, error, fatal, panic (default "debug")

Use "evml [command] --help" for more information about a command.

```

Options can also be specified in a `evml.toml` file in the `datadir`.

ex (evml.toml):
``` toml
log=info
[eth]
db = "/eth.db"
[babble]
listen="127.0.0.1:1337"
```

## Configuration

The application writes data and reads configuration from the directory specified
by the --datadir flag. The directory structure must respect the following
stucture:

```
host:~/.evm-lite$ tree
├── babble
│   ├── peers.json
│   └── priv_key.pem
├── eth
│   ├── genesis.json
│   ├── keystore
│   │   └── UTC--2018-10-14T11-12-24.412349157Z--633139fa62d5c27f454259ba59fc34773bd19457
│   └── pwd.txt
└── evml.toml
```

The above example shows a `babble` folder, but the general idea is that
consensus  configuration goes in a separate folder from the Ethereum
configuration.

The Ethereum genesis file defines Ethereum accounts and is stripped of all   
the Ethereum POW stuff. This file is useful to predefine a set of accounts  
that own all the initial Ether at the inception of the network.  

Example Ethereum genesis.json defining two account:
```json
{
   "alloc": {
        "629007eb99ff5c3539ada8a5800847eacfc25727": {
            "balance": "1337000000000000000000"
        },
        "e32e14de8b81d8d3aedacb1868619c74a68feab0": {
            "balance": "1337000000000000000000"
        }
   }
}
```

It is possible to enable evm-lite to control certain accounts by providing a  
list of encrypted private keys in the keystore directory. With these private
keys, evm-lite will be able to sign transactions on behalf of the accounts
associated with the keys.  

```
host:~/.evm-lite/eth/keystore$ tree
.
├── UTC--2016-02-01T16-52-27.910165812Z--629007eb99ff5c3539ada8a5800847eacfc25727
├── UTC--2016-02-01T16-52-28.021010343Z--e32e14de8b81d8d3aedacb1868619c74a68feab0
```

These keys are protected by a password. Use the `eth.pwd` flag to specify the
location of the password file.

**Needless to say you should not reuse these addresses and private keys**

## Database

EVM-Lite will use a LevelDB database to persist state objects. The file of the  
database can be specified with the `eth.db` flag which defaults to
`<datadir>/eth/chaindata`.  

## API
The Service exposes an API at the address specified by the --eth.listen flag for
clients to interact with Ethereum.  

### Get controlled accounts

This endpoint returns all the accounts that are controlled by the evm-lite
instance. These are the accounts whose private keys are present in the keystore.

example:
```bash
host:~$ curl http://[api_addr]/accounts -s | json_pp
{
   "accounts" : [
      {
         "address" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
         "balance" : 1337000000000000000000,
         "nonce": 0
      },
      {
         "address" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
         "balance" : 1337000000000000000000,
         "nonce": 0
      }
   ]
}
```
### Get any account

This method allows retrieving the information about any account, not just the
ones whose keys are included in the keystore.  

```bash
host:~$ curl http://[api_addr]/account/0x629007eb99ff5c3539ada8a5800847eacfc25727 -s | json_pp
{
    "address":"0x629007eb99ff5c3539ada8a5800847eacfc25727",
    "balance":1337000000000000000000,
    "nonce":0
}
```

### Send transactions from controlled accounts

Send a transaction from an account controlled by the evm-lite instance. The
transaction will be signed by the service since the corresponding private key is
present in the keystore.

example: Send Ether between accounts  
```bash
host:~$ curl -X POST http://[api_addr]/tx -d '{"from":"0x629007eb99ff5c3539ada8a5800847eacfc25727","to":"0xe32e14de8b81d8d3aedacb1868619c74a68feab0","value":6666}' -s | json_pp
{
   "txHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf"
}
```

### Get Transaction receipt
example:
```bash
host:~$ curl http://[api_addr]/tx/0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf -s | json_pp
{
   "to" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
   "root" : "0xc8f90911c9280651a0cd84116826d31773e902e48cb9a15b7bb1e7a6abc850c5",
   "gasUsed" : "0x5208",
   "from" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
   "transactionHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf",
   "logs" : [],
   "cumulativeGasUsed" : "0x5208",
   "contractAddress" : null,
   "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
}

```

Then check accounts again to see that the balances have changed:
```bash
{
   "accounts" : [
      {
         "address" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
         "balance" : 1336999999999999993334,
         "nonce":1
      },
      {
         "address" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
         "balance" : 1337000000000000006666,
         "nonce":0
      }
   ]
}
```
### Send raw signed transactions

Most of the time, one will require to send transactions from accounts that are
not  controlled by the evm-lite instance. The transaction will be assembled,
signed  and encoded on the client side. The resulting raw signed transaction
bytes can be submitted to evm-lite through the ```/rawtx``` endpoint.  

example:
```bash
host:~$ curl -X POST http://[api_addr]/rawtx -d '0xf8628080830f424094564686380e267d1572ee409368e1d42081562a8e8201f48026a022b4f68bfbd4f4c309524ebdbf4bac858e0ad65fd06108c934b45a6da88b92f7a046433c388997fd7b02eb7128f4d2401ef2d10d574c42edf15875a43ee51a1993' -s | json_pp
{
    "txHash":"0x5496489c606d74ad7435568393fa2c4619e64497267f80864109277631aa849d"
}
```

### Call contract without state change
The ```/call``` endpoint allows calling SmartContract code for READONLY operations.
these calls will NOT modify the EVM state.

The API input format is just like ```/tx``` endpoint, but will return the return data.

## Get consensus info

The ```/info``` endpoint exposes a map of information provided by the consensus
system.

example (with Babble consensus):
```bash
host:-$ curl http://[api_addr]/info | json_pp
{
   "rounds_per_second" : "0.00",
   "type" : "babble",
   "consensus_transactions" : "10",
   "num_peers" : "4",
   "consensus_events" : "10",
   "sync_rate" : "1.00",
   "transaction_pool" : "0",
   "state" : "Babbling",
   "events_per_second" : "0.00",
   "undetermined_events" : "22",
   "id" : "1785923847",
   "last_consensus_round" : "1",
   "last_block_index" : "0",
   "round_events" : "0"
}

```

## CLIENT

Please refer to [EVM-Lite Client](https://github.com/mosaicnetworks/evm-lite-client)
for Javascript utilities and a CLI to interact with the API.

## DEV

DEPENDENCIES

We use glide to manage dependencies:

```bash
[...]/evm-lite$ curl https://glide.sh/get | sh
[...]/evm-lite$ glide install
```
This will download all dependencies and put them in the **vendor** folder; it
could take a few minutes.

CONSENSUS

To add a new consensus system:

- implement the consensus interface (consensus/consensus.go)
- add a property to the the global configuration object (config/config.go)
- create the corresponding CLI subcommand in cmd/evml/commands/
- register that command to the root command

## DEPLOY

We provide a set of scripts to automate the deployment of testnets. This
requires [terraform](https://www.terraform.io/) and
[docker](https://www.docker.com/).

Support for AWS is also available (cf. deploy/)
