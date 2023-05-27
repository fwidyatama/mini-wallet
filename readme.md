
## Installation
Make sure you already clone this repo by running this command and you already install docker:

```bash
 git clone https://github.com/fwidyatama/mini-wallet.git
```

Go to the project directory

```bash
  cd mini-wallet
```

Install packages

```bash
  go mod tidy
```


## How to run

Open terminal and run this command one by one sequentially :

```bash
  1. make postgres
  2. make createdb
  3. make migrate up
  4. make dev
```
