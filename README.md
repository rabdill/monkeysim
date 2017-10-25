# MonkeySim Readme

To run:
```
go build

./monkeysim
```

`./monkeysim` takes an optional command line parameter for specifying how many monkeys to spin up at once. Default is 1, but can be any number. For example: `./monkeysim 11` will make 11 monkeys typing at 11 typewriters.

You will also need a text file in the same directory as the `monkeysim` binary called **target.txt** -- this is what the monkeys' success will be judged against. You can use any old file (type something yourself, grab a text file from Project Gutenberg, etc.), but check out the `getTarget()` function in [the helpers file](https://github.com/rabdill/monkeysim/blob/master/helpers.go) -- there are lots of things stripped out of the target, including line breaks, punctuation, and padding using multiple spaces. (In addition, all capital letters are shifted to lower-case.)