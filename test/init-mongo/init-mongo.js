db = db.getSiblingDB('catalog-history-service');
db.createCollection('datasources');
db.datasources.insert([
    {
        "_id": "test-id",
        "dataSourceType": "DCAT-AP-NO",
        "dataType": "dataset",
        "url": "http://url.com",
        "acceptHeaderValue": "text/turtle",
        "publisherId": "123456789",
        "description": "test source"
    }
]);
