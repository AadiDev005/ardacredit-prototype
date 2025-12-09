
# ArdaCredit Backend Prototype 


##  LIVE DEMO
```
curl -X POST http://localhost:8080/api/v1/apply \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"user123","amount":5000,"income":60000}'
```
** Response:**
```
{"user_id":"user123","score":850,"approved":true,"reason":"Approved based on income-to-loan ratio"}
```

##  Architecture (4 Fundamentals)

| Scaling Reads/Writes | Data Storage | Communication | Failure Handling |
|---------------------|--------------|---------------|------------------|
| Consistent hashing  | ACID txns (CockroachDB ready) | Kafka sagas | Async audit + retries |

##  Quick Start
```
# 1. Start infra (Kafka + Redis)
docker-compose up -d

# 2. Run API
go run main.go

# 3. Test credit application
curl -X POST http://localhost:8080/api/v1/apply \
  -d '{"user_id":"user123","amount":5000,"income":60000}'
```

##  Scales to 10M Users
- **MVP:** Single node (current)
- **10k DAU:** Redis caching
- **1M DAU:** CockroachDB sharding
- **10M DAU:** Multi-region + Istio

## 🛠️ Tech Stack
```
Go → Gin API → Kafka (audit) + Redis (cache)
↓
CockroachDB (production) + Kubernetes (EKS)
```

