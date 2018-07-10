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

* arguments need types like the in/out channels have (for URLs and files)
* out-of-process communication for regular processes needs to be better
* out-of-process with typed info somehow (some kind of JSON exchange I guess)
* How do we resolve the split between stdin and arguments? Imagine:
    fetch http://google.com | cat
    web-ls-thing | fetch
  These should function ideally, but right now it relies on handling stdin and
  arguments together. That feels awkward.

* We need a way to query data in a rich manner. Grep doesn't cut it on its own.
* select-cols & select-props (or sc/sp?)
    ls | select-props mtime | grep 2006 | rm
  ... remove any file from 2006
  This also means that the pipeline needs some way to 'backtrack'. We need to
  find the file object, not the property we selected. Perhaps a 'reference' data
  type, which provides both the sub-value, and a way to fetch the full
  data instance. Then, provide a custom cast method on the interface, maybe?

