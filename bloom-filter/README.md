# Bloom Filter in Go 🚀

This is a simple implementation of a **Bloom Filter** in Go — using the correct technique of **double hashing** to simulate multiple hash functions efficiently.

---

## 🌸 What is a Bloom Filter?

A **Bloom Filter** is a space-efficient probabilistic data structure used to test whether an element is a member of a set.

### Characteristics:

✅ Fast insert and lookup  
✅ Very low memory usage  
❌ Allows false positives (may say item is present when it's not)  
❌ No false negatives (if it says "no", the item was never added)  

### Typical use cases:

- Web crawlers → avoid revisiting same URL
- Distributed systems → avoid duplicate processing
- Databases → test existence of keys without disk access
- Caches → avoid cache misses

---

## ⚙️ How does it work?

- A **bit array** of size `m` is initialized to 0.
- There are **k independent hash functions**.
- To insert an item:
  - Hash it with k functions → get k positions → set those bits to 1.
- To test an item:
  - Hash it with k functions → check those k bits → if any bit is 0 → definitely not present, otherwise → possibly present.

---

## 🔍 Why multiple hash functions?

✅ To reduce **collisions** → fewer false positives.  
✅ To distribute bits more uniformly.  
✅ Each hash should act "independently".

---

## 🎩 Why double hashing trick?

Using **k truly independent hash functions** (FNV, SHA, Murmur, MD5, etc.) is expensive.

In practice, it's common to:

1️⃣ Compute one strong hash `h1(x)`  
2️⃣ Compute a second strong hash `h2(x)`  
3️⃣ Derive k hash functions as:

```go
hash_i(x) = h1(x) + i * h2(x)

👉 This gives good practical independence, is fast, and works well.

This is called the double hashing trick — its used by:

Redis Bloom Filter

Cassandra Bloom Filter

RocksDB Bloom Filter

ScyllaDB Bloom Filter