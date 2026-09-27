# The Specification Test

## Contents

- The Specification Principle
- Quick Test-Outline Diagnostic
- Layer Responsibilities
- Example: Handler Specification
- Example: Service Specification
- Example: Model Specification
- Cost Consideration
- Detailed Code Audit (3-Step Process)
- Quick Reference

## The Specification Principle

> If the specification of an object describes features beyond the primary responsibility of its abstraction layer, such features should be extracted into lower layers.

The specification test helps identify code that belongs in a different layer by examining what tests would verify.

## Quick Test-Outline Diagnostic

1. **Write test structure** (test names and table cases) without implementation
2. **Examine what the tests verify**
3. **Ask:** Does this test verify the object's primary responsibility?

If a test verifies something outside the layer's primary concern, that code should be extracted.

## Layer Responsibilities

| Layer | Primary Responsibility | Tests Should Verify |
|-------|------------------------|---------------------|
| Presentation (Handler) | HTTP handling | Authentication, authorization, response codes, redirects |
| Application (Service) | Use-case orchestration | Correct domain objects called, transaction boundaries |
| Domain (Model) | Business rules | Validations, state transitions, calculations |
| Infrastructure | Technical implementation | Data persistence, API calls |

## Example: Handler Specification

```text
HTTP test outline:
  missing signature → reject request
  invalid signature → reject request
  pull-request business behavior → move detailed cases to application tests
  issue business behavior → move detailed cases to application tests
```

The business logic tests indicate code that should move to the application or domain layer.

**After extraction:**

```text
HTTP tests: authentication, status, request/response contracts.
Application tests: pull-request event, issue event, unknown user.
Keep a small integration test proving that the entry point wires the operation correctly.
```

## Example: Service Specification

```text
ProcessOrder tests:
  valid order → operation commits
  payment failure → no local success
  inventory failure → transaction rolls back
Order tests:
  discount limit and minimum order total → domain rules
```

Discount and minimum order rules are domain logic—they belong in the Order model.

## Example: Model Specification

```text
Order tests: invariants, total calculation, discount, legal transitions.
Notification operation tests: delivery intent, transport failure, idempotency.
Warehouse adapter tests: outbound contract, timeout, response parsing.
```

Notification and external sync are not domain responsibilities.

## Cost Consideration

Higher-layer tests are:
- **Harder to write** (more context setup)
- **Slower to execute** (HTTP requests, full stack)
- **More brittle** (depend on more components)

Moving logic to lower layers enables faster, simpler, more focused tests.

| Test Type | Speed | Setup Complexity | Brittleness |
|-----------|-------|------------------|-------------|
| Model/unit | Fast | Low | Low |
| Service | Medium | Medium | Medium |
| Handler/request | Slow | High | High |
| System/integration | Slowest | Highest | Highest |

## Detailed Code Audit (3-Step Process)

### Step 1: List Responsibilities

For the code you're examining, list every responsibility it handles.

Example for `Handler.CreateOrder`:
- Parse order parameters
- Authenticate user
- Authorize order creation
- Validate inventory
- Calculate pricing
- Apply discounts
- Create order record
- Send confirmation email
- Sync to warehouse API
- Return JSON response

### Step 2: Categorize by Layer

| Responsibility | Layer | Belongs in Handler? |
|----------------|-------|------------------------|
| Parse parameters | Presentation | ✓ Yes |
| Authenticate user | Presentation | ✓ Yes |
| Authorize creation | Application | ✓ Yes (or policy) |
| Validate inventory | Domain | ✗ No |
| Calculate pricing | Domain | ✗ No |
| Apply discounts | Domain | ✗ No |
| Persist record | Infrastructure | ✗ No |
| Orchestrate email | Application | ✗ No (not handler's job) |
| Sync to API | Infrastructure | ✗ No |
| Return JSON | Presentation | ✓ Yes |

### Step 3: Extract

Move misplaced responsibilities to appropriate layers:

```go
func (s *Orders) Create(ctx context.Context, actor User, input CreateOrder) (Order, error) {
    order, err := NewOrder(actor.ID, input.Items)
    if err != nil { return Order{}, err }
    // The store transaction persists the order and notification intent together.
    if err := s.store.CreateWithOutbox(ctx, order); err != nil { return Order{}, err }
    return order, nil
}
// The HTTP handler parses input and maps errors to responses.
// Domain tests exercise NewOrder without HTTP or database setup.
```

## Quick Reference

**Handler should test:**
- HTTP status codes
- Redirects
- Authentication/authorization
- Parameter handling
- Response format

**Handler should NOT test:**
- Business rules
- Calculations
- State transitions
- External service behavior

**Service should test:**
- Correct objects orchestrated
- Transaction success/failure
- Error handling

**Service should NOT test:**
- Domain validation rules
- Business calculations
- HTTP concerns

**Model should test:**
- Validations
- Business rules
- Calculations
- State transitions

**Model should NOT test:**
- HTTP concerns
- Notification delivery
- External API calls
