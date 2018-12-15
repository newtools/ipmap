IPMap
-------
IPMap is a trie structure that has very fast O(1) lookups and O(N) space complexity.

It is not threadsafe.
One way to get around its lack of thread safety is to replace it on changes. More dynamic use cases that require thread safe use cases
should probably refactor the code with a read-write mutex.

Performance Details
-------------------
IPMap's performance is actually quite impressive. On x86 it averages ~450 clock cycles for Setting an ip address, and ~150 clock cycles to do one lookup (IpV4, Ipv6 is ~4x...obviously).
A couple of reasons the performance characteristics are good is that there is no hash function, and the physical memory layout tends to group the actual ip octets next to one another, so memory access
tends to only do only one L1 cache line load per lookup. This is better than a traditional hash table which will likely have two line cache loads per look up (one for the hashtable, and one
for the pointed to datastructure).