#!/bin/bash

record_count=${RECORD_COUNT:-100}

read_capacity=${READ_CAPACITY_UNITS:-5}
write_capacity=${WRITE_CAPACITY_UNITS:-5}

echo "Creating table"
# The || true is to prevent the script from failing if the table already exists. More complicated
# logic can be done to assess if the table already exists and if so, do not run the command,
# however for sake of simplicity, this has not been implemented.
aws dynamodb create-table --endpoint-url ${DYNAMODB_ENDPOINT} \
    --table-name ${ITEM_TABLE} \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=${read_capacity},WriteCapacityUnits=${write_capacity} \
    || true

echo "Populating table"
# Required fields: id, name, description
base_id="fadd2e37-3c27-44e8-b5e3-821b21b" # Random UUID minus the last 5 characters
base_title="Item_"
base_desc="Desc_"
for i in $(seq 1 1 ${record_count}); do
    id=$(printf "%s%05d" ${base_id} ${i})
    title=$(printf "%s%05d" ${base_title} ${i})
    desc=$(printf "%s%05d" ${base_desc} ${i})
    echo "${id} - ${title} - ${desc}"
    aws dynamodb put-item --endpoint-url ${DYNAMODB_ENDPOINT} \
        --table-name ${ITEM_TABLE} \
        --item "{\"id\": {\"S\": \"${id}\"}, \"title\": {\"S\": \"${title}\"}, \"description\": {\"S\": \"${desc}\"}}"
done

echo "Finished populating all data"
