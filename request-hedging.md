# 🛡️ Request Hedging Behind an ALB – Summary

## What is Request Hedging?
Request hedging is a technique to reduce **tail latency** (e.g. 95th/99th percentile) by sending **duplicate requests** and using the **first successful response**.

---

## Scenario
- **Client sends requests** to a service behind a **single DNS** entry.
- The DNS points to an **Application Load Balancer (ALB)**.
- The ALB distributes traffic to **multiple backend instances** (e.g. 3).

---

## Key Insights

### ✅ When does hedging help?
- The ALB must **distribute hedged requests to different instances**.
- **Sticky sessions must be disabled** (e.g., no IP affinity, cookies).
- The backend pool must be **large enough** to avoid repeated collisions.
- The ALB's algorithm (e.g., round-robin or least-connections) must support **non-deterministic routing**.

### ❌ When does hedging hurt?
- If both requests go to the **same backend**, it adds load but no benefit.
- Hedging increases **backend load**, **logs**, and **wasted compute**.
- In high-load systems, it can **worsen performance**.

---

## Best Practices

- ✅ Only hedge requests on endpoints with **high latency variance**.
- ✅ Use **adaptive hedging** based on load or observed latencies.
- ✅ Set a **delay threshold** (e.g., 50ms) before sending a second request.
- ✅ Discard the slower response, but handle cleanup gracefully if needed.

---

## Summary Formula

> For request hedging to reduce tail latency:  
> the **ALB must spread requests** across **many backend instances**,  
> with **no correlation** between hedged copies.


```go 
func hedge() string {
    ch := make(chan string)
    ctx, cancel := context.WithCancel(context.Background())

    for i := 0; i < 5; i++ {
        go func(ctx *context.Context, ch chan string, i int) {
            log.Println("gogo: ", i)
            if request(ctx, "http://localhost:8030", i) {
                ch <- fmt.Sprintf("finishedfrom %v", i)
                log.Println("completed: ", i)
            }
        }(&ctx, ch, i)
    }
    select {
    case s := <-ch:
        cancel()
        return s
    case <-time.After(5 * time.Second):
        cancel()
        return "requests took more than 5 seconds"
    }
}
```

