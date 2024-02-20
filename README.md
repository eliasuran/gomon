# gomon

<h3 align="center">Make changes to your go http server without having to shutdown and rerun the program each time</h3>

---

<div align="center">

  <img height="180" src="https://github.com/eliasuran/gomon/assets/118540201/d1c78de7-1835-4106-af8f-0a053055a048" alt="gopher" />

  <img height="160" src="https://github.com/eliasuran/gomon/assets/118540201/6b8b7188-e60d-44ff-b6c7-a9a6fe987c69" alt="database" />
  
</div>

---

## Preview

* Automatic updates when making changes to your api 🚀

* Built for the net/http standard package 🔥



## Installation

### Homebrew (in progress)

### Manual installation

```sh
git clone https://github/eliasuran/gomon
cd gomon
go install
```

Run the program

```sh
gomon "path to dir with http server"
```

## Usage

```shell
gomon "path/to/dir/with/main.go/file"
```


* Run gomon and point to a directory containing a main.go file (currently only works with main.go files, working on fix)

* Gomon will build a binary and run it. When you make changes, it will automatically build a new binary with the latest changes.

* If it encounters an error, it will stay on the latest working version until the error is resolved.
