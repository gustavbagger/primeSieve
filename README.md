# Sieving for sparse sets of primes
This library provides an efficient algorithm for computing all integers n below some dynamic bound satisfying
1) n has exactly omega distinct prime factors for some fixed omega
2) `n+1` is a prime number p
The dynamic bound in question is a prime sieve given by `2*r*(2 + (s-1)/delta)^r p^{1/4(1+1/r)} < p^{1/2}-2` and `delta>0` where integer s can be chosen freely in the interval `[0,omega]` and `delta = 1 - \sum_{i=1}^s(1/p_i)` for a choice of s distinct prime factors `p_i` of n. 
## Motivation
A long-standing conjecture by Grosswald states that, for all primes `p>409`, the least primitive root modulo p is less than `p^{1/2} - 2`. The conjecture is resolved under the Generalised Riemann Hypothesis, and also unconditionally for `p>10^{56}`. This algorithm computes the sparse set of potential *bad primes* for which one needs to directly compute the least primitive root modulo p to varify the conjecture. 

> A slight modification of the sieving criterion used here would allow one to compute other sparse sets of primes subject to conditions **(1)** and **(2)**. For more information on the sieving techniques used, see [The Modified Prime Sieve](https://arxiv.org/abs/2507.21515).
## Quick Start
```bash
# install using Go toolchain
go install github.com/gustavbagger/PropertyZ

# run
PropertyZ <arg1> <?arg2> <?arg3>
```
## Usage
The programme can be run with 1,2 or 3 arguments. If the programme is run with 1 or 2 arguments, it will provide an interval (dependant on omega) in which each value n must lie in order to satisfy the requirements. It will also print out the optimal choices of s and delta in order to sieve effectively. If the programme is run with 3 arguments, it will print out any integer n satisfying the criteria in the given search-space. The following provides a control flow for what happens:
### compute the intervals for omega = n
```
PropertyZ <arg1>
-> runs printIntervals(arg1,arg1)
  -> prints "${upperBound} > p > ${lowerBound} ---- omega,s,delta = ${omega},${optimal s value},${optimal delta value}"
```
### compute the intervals for n >= omega >= m
```
PropertyZ <arg1> <arg2>
-> runs printIntervals(arg1,arg2)
  -> prints "${upperBound} > p > ${lowerBound} ---- omega,s,delta = ${arg1},${optimal s value},${optimal delta value}"
  -> prints "${upperBound} > p > ${lowerBound} ---- omega,s,delta = ${arg1 - 1},${optimal s value},${optimal delta value}"
  ...
  -> prints "${upperBound} > p > ${lowerBound} ---- omega,s,delta = ${arg2},${optimal s value},${optimal delta value}"
```
### compute all admissible values n with omega = n and n+1=p <= a * 10^b
```
PropertyZ <arg1> <arg2> <arg3>
-> runs search(arg1,arg2,arg3) corresponding to omega = arg1 with p<= arg2 * 10^arg3
  -> preloads all primes p <= 10^6
  -> computes initial optimal value for s
#saves storing a huge slice with unneeded primes to loop over
  -> prunes primeList down to the smallest size needed for exhaustiveness
#generates all squarefree numbers in the interval sequentially via recursive nested loops
  -> runs recursiveLoop at initial depth 0
#squarefree numbers are iterated as index slices of primeList
    -> if depth is not maximal, recurse through indexes in current depth
    -> if the current index slice gives rise to n >= currentUpperBound bail out and backtrack to a lower loop and increment there
    -> if depth is maximal runs treeSearch with current set of primes found
#treeSearch traverses admissible exponents for the squarefree number
      -> if position is maximal, check if the current modularity conditions for a suitable choice of 'medium size' primes (relative to omega) are satisfied
        -> if they are and n+1 is prime (using Miller-Rabin, 32 loops), print the prime indexes in primeList and their exponents
      -> if position is not maximal:
        -> update the sieve bound using current values
        -> update the modularity conditions for current exponents
        -> if the exponents are too large, bail out early
        -> else, continue traversing the exponent tree at the next position
```
## Contributing

### Clone the repo

```bash
git clone https://github.com/gustavbagger/PropertyZ
cd PropertyZ
```

### Build the compiled binary

```bash
go build
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.