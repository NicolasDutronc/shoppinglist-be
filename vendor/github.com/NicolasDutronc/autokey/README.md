# autokey
autokey is a Go library that automatically generates and rotates a base 64 encoded key which can be used to sign JWTs.

This package provides a Manager struct that has only two methods:
 * `Start(interrupt chan struct{})` that starts the manager loop (it should be run this in its own goroutine). The interrept channel is used to signal to the listener goroutine that the manager has stops.
 * `Stop()` is used to stops the manager loop by closing the timer and its stop channel.

Once started, the manager will try to find if a key has already been stored and create a new one if there is no key stored. Then it will start a timer and automatically renew the key periodically.

The storage of the key is defined by an interface so it is pluggable. A environment storage is provided.

This package also provides the consumer interface that defines how to retrieve the key. So if your struct or function needs to access the key, you can use this interface. It is helpful for tests and decoupling.