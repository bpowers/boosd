boosd - object-oriented system dynamics engine
==============================================

This is an implementation of an object-oriented system dynamics
engine.  Models have a simple, text-based representation:

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

This is (heavily) influenced by [go](http://golang.org)'s syntax, but
boosd models are declarative rather than imperitive.  Currently the
best documentation and background on boosd is in
[my thesis](http://bpowers.github.com/thesis.pdf).

The basic idea is to parse the above model into an AST, and then
convert that AST into a go AST, which is compiled into a static go
binary.

development
-----------

I'm currently testing by running a commandline like this:

    time (f=`mktemp --suffix=.go`; go install && boosd models/exp.osm >$f; go run $f; rm $f)

license
-------

gocalc is offered under the MIT license, see LICENSE for details.
