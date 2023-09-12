# simple-store
This is a simple store API that simulates a purchase and retrieve this all purchases or just one purchase
Here are some examples:
GET: http://localhost:8001/purcharse/64ffba95530ac4ccc8d2e63d/currency/Real
GET: http://localhost:8001/purcharse
POST: http://localhost:8001/purcharse
BODY:
{
    "Description":"first purchase",
    "Date":"2023-09-11T23:22:21Z",
    "Amount": 1101.12
}
