db = db.getSiblingDB('catalog-history-service');
db.createCollection('concepts');
db.concepts.insert(
{
    "person": {
        "id": "123",
        "email": "example@example.com",
        "name": "John Doe"
    },
    "operations": [
        {
            "op": "replace",
            "path": "/name",
            "value": "Jane"
        },
        {
            "op": "remove",
            "path": "/height"
        }
    ]
}
);
