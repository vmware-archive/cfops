uggo
=====

Ungraceful Gnu-ification for Go.

This helps Go commandline apps to behave a bit more 'coreutils'-like, or 'Gnu-ish'. 

Yes, there are many other flagset libraries out there.
This one provides a wrapper around the existing `flag.FlagSet`, embedding and embellishing it.
You can use uggo as a drop-in replacement for flag.FlagSet, or you can specify a couple of different behaviours. Up to you.

The reason why it's called 'ungraceful' is because it's an approximation of gnu-like behaviour, as close as I could get while still wrapping `flag.FlagSet`. 
It's good enough for me anyway. 
Also, the STDIN-pipe-detection is platform-specific, relying on an undocumented behaviour in Windows.

## Initial features

 * 'Gnuify' options such as `-lah` so they are treated as `-l -a -h` (whereas `--lah` is treated as-is, as a single option `--lah`)
 * detect whether STDIN is being piped from another process
 * 'aliased flags' (wrapper methods making it easier to define short and long options, or just multiple options)
 * Re-jigged format of 'Usage' to incroporate 'aliased flags'

## See Also

These programs use uggo:

 * [https://github.com/laher/someutils](https://github.com/laher/someutils)
 * [https://github.com/laher/wget-go](https://github.com/laher/wget-go)
 * [https://github.com/laher/scp-go](https://github.com/laher/scp-go)
