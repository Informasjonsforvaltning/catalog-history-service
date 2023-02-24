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
        },
        {
            "op": "add",
            "path": "/name",
            "value": "Jane Test"
        }
    ]
}
);

db.concepts.insert(
    {
        "id": "789",
        "resourceId": "123456789",
        "person": {
            "id": "789",
            "email": "example3@example.com",
            "name": "Joe Doe"
        },
        "datetime": "2019-01-03T00:00:00Z",
        "operations": [
            {
                "op": "add",
                "path": "/name",
                "value": "Joe"
            }
        ]
    }
    );

    db.concepts.insert({
        "id": "456",
        "resourceId": "123456789",
        "person": {
            "id": "456",
            "email": "example2@example.com",
            "name": "Sarah Doe"
        },
        "datetime": "2019-01-02T00:00:00Z",
        "operations": [
            {
                "op": "replace",
                "path": "/name",
                "value": "Sarah"
            }
        ]
    });
    
    db.concepts.insert({
        "id": "012",
        "resourceId": "123456789",
        "person": {
            "id": "012",
            "email": "example4@example.com",
            "name": "Bob Doe"
        },
        "datetime": "2019-01-04T00:00:00Z",
        "operations": [
            {
                "op": "replace",
                "path": "/name",
                "value": "Bob"
            }
        ]
    });
    
