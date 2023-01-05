db = db.getSiblingDB('catalog-history-service');
db.createCollection('begreper');
db.datasources.insert([
    {
        "_id": "test-id",
        "term": "someTerm",
        "def": "someDef",
    }
]);
