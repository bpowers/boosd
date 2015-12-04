boosd - object-oriented system dynamics engine
==============================================

Boosd is an implementation of an object-oriented system dynamics
language.  Models have a simple, text-based representation:

    main model {
            timespec = {
                    start:     0
                    end:       50
                    dt:        .5
                    save_step: 1
            }

            rate = .07
            in flow = rate * accum

            accum stock = {
                    inital: 200
                    inflow: in
            }
    }

This is influenced by [go](http://golang.org)'s syntax, but boosd
models are declarative rather than imperitive.  Currently the best
documentation and background on boosd is [my
thesis](https://bpowers.net/thesis.pdf).

The advantage of boosd is the formal grammar and semantics, and that
the model equations are cleanly presented in a compact representation.


development
-----------

This is in a very rough stage.  I currently test by running a pipeline like this:

    time (f=`mktemp --suffix=.go`; go install && boosd models/exp.osm >$f; go run $f; rm $f)

license
-------

boosd is offered under the MIT license, see LICENSE for details.
