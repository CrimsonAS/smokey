# a shell experiment

What if shells were a bit more powerful? What if they had more knowledge of what
was passing between them than just text?

Some thoughts:

    ls filename.txt | cat

... stream a file to cat

    fetch http://google.com | cat

... stream a HTTP GET to cat

    ls | grep .txt | cat

... cat all files matching .txt

# TODO

* actually implement raw terminal mode
* fix standard process interaction (and make it work for e.g. vim)
* arguments need types like the in/out channels have (for URLs and files)?
* out-of-process communication for regular processes needs to be better
* out-of-process with typed info somehow (some kind of JSON exchange I guess)
* How do we resolve the split between stdin and arguments? Imagine:
    fetch http://google.com | cat
    web-ls-thing | fetch
  These should function ideally, but right now it relies on handling stdin and
  arguments together. That feels awkward.

