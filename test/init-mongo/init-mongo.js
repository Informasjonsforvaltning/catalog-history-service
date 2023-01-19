db = db.getSiblingDB('catalog-history-service');
db.createCollection('concepts');
db.concepts.insert([
    {
        "_id": "test-id",
        "term": "someTerm",
        "def": "someDef",
    }
]);
