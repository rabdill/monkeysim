# MonkeySim Readme

MonkeySim only has one requirement (other than Golang) -- it uses **dep** for dependency management. Check out [their repository](https://github.com/golang/dep) for info on how to install it locally. If you don't want to deal with all that, the only external dependency we currently have is **[gin-gonic/gin](github.com/gin-gonic/gin)**, so running `go get github.com/gin-gonic/gin` should work just as well as running `dep ensure`, at least for now.

To run:
```
dep ensure
go build
./monkeysim
```

`./monkeysim` takes an optional command line parameter for specifying how many monkeys to spin up at once. Default is 1, but can be any number. For example: `./monkeysim 11` will make 11 monkeys typing at 11 typewriters. Once it's running, open a browser and go to `localhost:8080` to see the status of your simulation.

You will also need a text file in the same directory as the `monkeysim` binary called **target.txt** -- this is what the monkeys' success will be judged against. You can use any old file (type something yourself, grab a text file from Project Gutenberg, etc.), but check out the `getTarget()` function in [the helpers file](https://github.com/rabdill/monkeysim/blob/master/helpers.go) -- there are lots of things stripped out of the target, including line breaks, punctuation, and padding using multiple spaces. (In addition, all capital letters are shifted to lower-case.)
