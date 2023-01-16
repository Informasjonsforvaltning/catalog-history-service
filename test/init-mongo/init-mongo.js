db = db.getSiblingDB('catalog-history-service');
db.createCollection('concepts');
db.concepts.insert([
    {
        "_id": "test-id",
        "term": "someTerm",
        "def": "someDef",
    }
]);

// Create a variable to hold the update struct
var update = {
    patches: [
        {
            "op": "replace",
            "path": "/term",
            "value": "newTerm"
        },
        {
            "op": "replace",
            "path": "/def",
            "value": "newDef"
        }
    ]
};
