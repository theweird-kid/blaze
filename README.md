# Blaze

Blaze is a simple distributed task scheduler built to reliably execute background jobs in a failure-prone environment.

It focuses on correctness first, using leases and idempotency to safely run tasks across multiple workers.

---

## Features

- Distributed job execution
- Lease-based worker coordination
- Idempotent task execution
- Retry on failure
- HTTP-based jobs
- Horizontally scalable workers

---

## How It Works

1. A job is created via the API
2. Workers poll for available jobs
3. A worker acquires a lease on a job
4. The job is executed
5. The job is marked complete or retried on failure

Leases ensure that only one worker executes a job at a time.  
Idempotency ensures retries do not cause duplicate side effects.

---

## Job Types

### HTTP Job

```json
{
  "type": "HTTP",
  "http": {
    "method": "POST",
    "url": "https://example.com/task",
    "headers": {
      "Authorization": "Bearer token"
    },
    "body": {
      "id": 123
    }
  }
}
