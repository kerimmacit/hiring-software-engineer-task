# Software Engineering Task: Ad Bidding Service

## Overview

A small Go project that returns the best ad for a placement and records basic tracking events.

## Deliverables

1. **Ad Selection Logic**: Implemented
2. **Tracking Endpoint**: Implemented
3. **Relevancy System**: Implemented 
4. **Input Validation**: Added
5. **Documentation**: Updated


## Storage Solutions

- **Source of truth**: Postgres or MongoDB or similar + Replication
- **Cache**: Redis or in-memory
- **Columnar**: For logging, tracking, analytics

This solution will be scalable, reliable, and performant

Also introduced repository level with interface, in order to make it easy to integrate new storage solution

## Scaling Considerations

**1. How would you scale this service to handle millions of ad requests per minute?**

- Application should be stateless. We need to use both real dbms and cache
- Load balancer should be introduced to handle networking to multiple pods. Also, autoscaling can be used for increase pod counts according to demand
- For suitable endpoints, streaming can be introduced. ex: POST /tracking endpoint can be converted to event publisher to kafka, and consumers can log these events
- If needed, gRPC can be introduced instead of REST

**2. What bottlenecks do you anticipate and how would you address them?**

- Tracking event write spikes: Memory buffer + Kafka
- Autoscaling for request spikes

**3. How would you design the system to ensure high availability and fault tolerance?**

- Multi region deployments
- Circuit breakers, timeouts, retry logics

**4. What data storage and access patterns would you recommend for different components (line items, tracking events, etc.)?**

- Tracking: kafka + columnar db
- Line Items / Ads: postgres or mongodb for source of truth + redis or in-memory for cache + streaming/messaging for invalidate cache 

**5. How would you implement caching to improve performance?**

- Local cache for speed, redis for scalable shared
- Cache invalidation logic
- Cold Start: cache should be filled when startup

## Future Improvements

**Test**: Only added unit tests on relevancy scoring for simplicity. More and various tests can be added

**Error**: Instead of generic error, specific error type can be introduced

**UUID**: In order to prevent any sort of duplication in line item creation, uuid should be created from requester side and used as a dedupe key in server side

**Budget Usage**: We might use provision logic while winning the ads, and consume the budget when there are tracking events
