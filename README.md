# Elevator Control

## Compile and Run
1. Compile by running `make`
2. Run by executing `./elevator-control`

## Flags
By default there are 4 eleveators.  This can be changed by setting the `-n` flag.

## Scheduling
This program splits scheduling into two parts.  The first part is which elevator should be sent to pickup the customer.  The second part is to ensure we a

### Which elevator should be sent?
When a request for a pickup comes in, each elevator is scored based on three factors.  In order of weighting, they are.
1. Is the elevator currently not in use?  +3 points if it is not.
2. Is the elevator the closet elevator to the pickup floor?  +1 point if it is.
3. Is the elevator going the same direction?  +1 point if it is.

Using this scoring algorithm we can reasonably expect we will get an un-used elevator or one that is close and going int he same direction.
