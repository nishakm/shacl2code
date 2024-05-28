package v3_0_test

import (
	spdx "example.com/shacl2code/src/shacl2code/lang/exported_flat"
	"strings"
	"testing"
)

func Test_json(t *testing.T) {
	// Generate using:
	// .venv/bin/python -m shacl2code generate -i https://spdx.org/rdf/3.0.0/spdx-model.ttl -i https://spdx.org/rdf/3.0.0/spdx-json-serialize-annotations.ttl -x https://spdx.org/rdf/3.0.0/spdx-context.jsonld golang --output src/shacl2code/lang/exported_flat/code.go --package v3_0 --license MIT --export-structs True --struct-suffix "" --include-view-pointers True
	//
	doc, err := spdx.FromJSON[*spdx.SpdxDocument](spdx.LD_CONTEXT, strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}

	// doc.SetProfileConformances(spdx.ProfileIdentifierType_Software)

	t.Logf("got: %+v", doc)
}

const sample = `
{
    "@context": "https://spdx.org/rdf/3.0.0/spdx-context.jsonld",
    "@graph": [
        {
            "type": "CreationInfo",
            "@id": "_:creationinfo",
            "createdBy": [
                "http://spdx.example.com/Agent/JoshuaWatt"
            ],
            "specVersion": "3.0.0",
            "created": "2024-03-06T00:00:00Z"
        },
        {
            "type": "Person",
            "spdxId": "http://spdx.example.com/Agent/JoshuaWatt",
            "name": "Joshua Watt",
            "creationInfo": "_:creationinfo",
            "externalIdentifier": [
                {
                    "type": "ExternalIdentifier",
                    "externalIdentifierType": "email",
                    "identifier": "JPEWhacker@gmail.com"
                }
            ]
        },
        {
            "type": "SpdxDocument",
            "spdxId": "http://spdx.example.com/Document1",
            "creationInfo": "_:creationinfo",
            "rootElement": [
                "http://spdx.example.com/BOM1"
            ],
            "profileConformance": [
                "core",
                "software"
            ]
        },
        {
            "type": "software_Sbom",
            "spdxId": "http://spdx.example.com/BOM1",
            "creationInfo": "_:creationinfo",
            "rootElement": [
                "http://spdx.example.com/Package1"
            ],
            "software_sbomType": [
                "build"
            ]
        },
        {
            "type": "software_Package",
            "spdxId": "http://spdx.example.com/Package1",
            "creationInfo": "_:creationinfo",
            "name": "my-package",
            "software_packageVersion": "1.0",
            "software_downloadLocation": "http://dl.example.com/my-package_1.0.0.tar",
            "builtTime": "2024-03-06T00:00:00Z",
            "originatedBy": [
                "http://spdx.example.com/Agent/JoshuaWatt"
            ]
        },
        {
            "type": "software_File",
            "spdxId": "http://spdx.example.com/Package1/myprogram",
            "creationInfo": "_:creationinfo",
            "name": "myprogram",
            "software_primaryPurpose": "executable",
            "software_additionalPurpose": [
                "application"
            ],
            "software_copyrightText": "Copyright 2024, Joshua Watt",
            "builtTime": "2024-03-06T00:00:00Z",
            "originatedBy": [
                "http://spdx.example.com/Agent/JoshuaWatt"
            ]
        },
        {
            "type": "Relationship",
            "spdxId": "http://spdx.example.com/Relationship/1",
            "creationInfo": "_:creationinfo",
            "from": "http://spdx.example.com/Package1",
            "relationshipType": "contains",
            "to": [
                "http://spdx.example.com/Package1/myprogram"
            ],
            "completeness": "complete"
        }
    ]
}
`
