# celltomaton-service

## About
used to deploy celltomaton as a service using docker containers

may be reachable on [celltomaton.eah.space](http://celltomaton.eah.space)

## Usage
Make a post request with a json object containing an initial vector, height and rule.

```
{
	"array": [0, 1, 1],
	"height": 10,
	"rule": 3
}
```

This will return a matrix with the width of the initial vector and the given height.
