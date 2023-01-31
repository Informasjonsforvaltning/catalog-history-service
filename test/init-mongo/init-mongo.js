db = db.getSiblingDB('catalog-history-service');
db.createCollection('concepts');
db.concepts.insert(
{
    "id": "123",
    "resourceId": "123456789",
    "person": {
        "id": "123",
        "email": "example@example.com",
        "name": "John Doe"
    },
    "datetime": "2019-01-01T00:00:00Z",
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
