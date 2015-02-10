goutil -> itertools

[![wercker status](https://app.wercker.com/status/092b6dbc492403c29b16676d5c5d5861/m/ "wercker status")](https://app.wercker.com/project/bykey/092b6dbc492403c29b16676d5c5d5861)
======

Iterate - function which turns all typically
            iterable objects into object accepted by the range function

Range - function to create a generator
          similar to the xrange python function.

Pair - a type, which is
        in line with the c++ Pair object

Map - function which allows us to iterate over Slice, Array, chan and Map types
      and execute a function on each of the elements in those
      structures.

CMap - version of the Map function
        which will loop and run the function concurrently
        using go routines.

Filter - function which acts as a generator filtering the passed in iterable object of
      values which do equate to true according to the passed in 
      evaluation function

CFilter - concurrent version of Filter using go routines
      
FilterFalse - function which acts as a generator filtering the passed in iterable object of 
      values which do not equate to true according to the passed in evaluation function
      ( this function has the exact opposite effect as Filter )

CFilterFalse - concurrent cersion of FilterFalse using go routines

ZipLongest - Make an iterator that aggregates elements from each of the iterables. If the iterables are of uneven length,       missing values are filled-in with fillvalue. Iteration continues until the longest iterable is exhausted


Zip - Make an iterator that aggregates elements from each of the iterables. If the iterables are of uneven length, indexes with missing values are dropped.
