# a shell experiment

What if shells were a bit more powerful? What if they had more knowledge of what
was passing between them than just text?

Some thoughts:

    file://filename.txt | cat

... stream a file to cat

    http://google.com | cat

... stream a HTTP GET to cat

    ls | cat

ls will produce file objects, and give them to cat

# TODO

* arguments need types like the in/out channels have (for URLs and files)
* out-of-process communication for regular processes needs to be better
* out-of-process with typed info somehow (some kind of JSON exchange I guess)
