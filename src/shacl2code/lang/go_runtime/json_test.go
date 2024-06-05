package v3_0_test

import (
	"bytes"
	"encoding/json"
	spdx "example.com/shacl2code/src/shacl2code/lang/go_runtime"
	"strings"
	"testing"
)

func Test_example1(t *testing.T) {
	// Generate using:
	// .venv/bin/python -m shacl2code generate -i https://spdx.org/rdf/3.0.0/spdx-model.ttl -i https://spdx.org/rdf/3.0.0/spdx-json-serialize-annotations.ttl -x https://spdx.org/rdf/3.0.0/spdx-context.jsonld golang --output src/shacl2code/lang/go_runtime/code.go --package v3_0 --license MIT --export-structs True --struct-suffix "" --include-view-pointers True
	//
	//doc, err := spdx.FromJSON[*spdx.SpdxDocument](spdx.LD_CONTEXT, strings.NewReader(sample))
	graph, err := spdx.LD_CONTEXT.FromJSON(strings.NewReader(example1))
	if err != nil {
		t.Fatal(err)
	}
	ctx := spdx.SpdxContext(graph)

	//ctx.GetRelationships().Where().FromEq(pkg1).And().TypeEq(spdx.RelationshipType_Contains)

	docs := ctx.GetSpdxDocuments()
	if len(docs) == 0 {
		t.Fatal("no documents")
	}
	_ = ctx.GetPackages()
	doc := docs[0]

	// doc.SetProfileConformances(spdx.ProfileIdentifierType_Software)

	t.Logf("got: %+v", doc)

	relationships := ctx.GetRelationships()

	for _, r := range relationships {
		t.Logf("from: %+v, to: %+v", stringify(r.GetFrom()), stringify(r.GetTo()))
	}
}

func Test_example2(t *testing.T) {
	// Generate using:
	// .venv/bin/python -m shacl2code generate -i https://spdx.org/rdf/3.0.0/spdx-model.ttl -i https://spdx.org/rdf/3.0.0/spdx-json-serialize-annotations.ttl -x https://spdx.org/rdf/3.0.0/spdx-context.jsonld golang --output src/shacl2code/lang/exported_flat/code.go --package v3_0 --license MIT --export-structs True --struct-suffix "" --include-view-pointers True
	//
	//doc, err := spdx.FromJSON[*spdx.SpdxDocument](spdx.LD_CONTEXT, strings.NewReader(sample))
	graph, err := spdx.LD_CONTEXT.FromJSON(strings.NewReader(example2))
	if err != nil {
		t.Fatal(err)
	}
	ctx := spdx.SpdxContext(graph)

	docs := ctx.GetSpdxDocuments()
	if len(docs) == 0 {
		t.Fatal("no documents")
	}
	doc := docs[0]

	// doc.SetProfileConformances(spdx.ProfileIdentifierType_Software)

	t.Logf("got: %+v", doc)

	relationships := ctx.GetRelationships()

	for _, r := range relationships {
		t.Logf("from: %+v, to: %+v", stringify(r.GetFrom()), stringify(r.GetTo()))
	}
}

func stringify(o any) string {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	_ = enc.Encode(o)
	return buf.String()
}

const example1 = `
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

const example2 = `
{
	 "@context": "https://spdx.org/rdf/3.0.0/spdx-context.jsonld",
	 "@graph": [
		 {
			 "type": "Person",
			 "spdxId": "urn:jane-doe-1@acme.com-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "Application Owner Jane Doe",
			 "externalIdentifiers": [
				 {
					 "type": "ExternalIdentifier",
					 "externalIdentifierType": "email",
					 "identifier": "jane-doe-1@acme.com"
				 }
			 ]
		 },
		 {
			 "type": "Organization",
			 "spdxId": "urn:acme.com-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "Acme Company"
		 },
		 {
			 "type": "Person",
			 "spdxId": "urn:github.com-indutny-c4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "Fedor Indutny",
			 "externalIdentifiers": [
				 {
					 "type": "ExternalIdentifier",
					 "externalIdentifierType": "other",
					 "identifierLocator": "https://github.com/indutny"
				 }
			 ]
		 },
		 {
			 "type": "Organization",
			 "spdxId": "urn:github.com-alpinelinux-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "Alpine Linux"

		 },
		 {
			 "type": "CreationInfo",
			 "@id": "_:creationinfo",
			 "specVersion": "3.0.0",
			 "createdBy": [
				 "urn:jane-doe-1@acme.com-4fe40e24-20e3-11ee-be56-0242ac120002",
				 "urn:acme.com-4fe40e24-20e3-11ee-be56-0242ac120002"
			 ],
			 "created": "2024-05-02T00:00:00Z"
		 },
		 {
			 "type": "SpdxDocument",
			 "spdxId": "http://spdx.example.com/Document1",
			 "creationInfo": "_:creationinfo",
			 "profileConformance": [
				 "core",
				 "software"
			 ],
			 "rootElement": [
				 "urn:example13-sbom.com-4fe40e24-20e3-11ee-be56-0242ac120002"
			 ]
		 },
		 {
			 "type": "software_Sbom",
			 "spdxId": "urn:example13-sbom.com-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "profileConformance": [
                                 "core",
                                 "software"
                         ],
			 "rootElement": [
				 "urn:product-acme-application-1.3-4fe40e24-20e3-11ee-be56-0242ac120002"
			 ]
		 },
		 {
			 "type": "software_Package",
			 "spdxId": "urn:product-acme-application-1.3-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "Acme Application",
			 "software_packageVersion": "1.3",
			 "suppliedBy": "urn:acme.com-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "software_primaryPurpose": "application"
		 },
		 {
			 "type": "software_Package",
			 "spdxId": "urn:npm-elliptic-6.5.2-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "npm-elliptic",
			 "software_packageVersion": "6.5.2",
			 "suppliedBy": "urn:github.com-indutny-c4fe40e24-20e3-11ee-be56-0242ac120002",
			 "software_primaryPurpose": "library",
			 "externalIdentifiers": [
				 {
					 "type": "ExternalIdentifier",
					 "externalIdentifierType": "other",
					 "identifierLocator": "https://github.com/indutny/elliptic/releases/tag/v6.5.2"
                                 }
                         ]

		 },
		 {
			 "type": "software_Package",
			 "spdxId": "urn:container-alpine-latest-sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "alpine:latest",
			 "software_packageVersion": "69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a",
			 "suppliedBy": "urn:github.com-alpinelinux-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "software_primaryPurpose": "container"
		 },
		 {
			 "type": "software_Package",
			 "spdxId": "urn:openssl-3.0.4-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "name": "openssl",
			 "software_packageVersion": "3.0.4",
			 "software_primaryPurpose": "library"
		 },
		 {
			 "type": "Relationship",
			 "spdxId": "urn:acme-relationship-1-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "from": "urn:product-acme-application-1.3-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "to": "urn:jane-doe-1@acme.com-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "relationshipType": "availableFrom"
		 },
		 {
			 "type": "Relationship",
			 "spdxId": "urn:acme-relationship-2-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "from": "urn:product-acme-application-1.3-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "to": "urn:npm-elliptic-6.5.2-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "relationshipType": "contains"
		 },
		 {
			 "type": "Relationship",
			 "spdxId": "urn:acme-relationship-3-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "from": "urn:product-acme-application-1.3-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "to": "urn:container-alpine-latest-sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "relationshipType": "depends_on"
		 },
		 {
			 "type": "Relationship",
			 "spdxId": "urn:acme-relationship-4-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "creationInfo": "_:creationinfo",
			 "from": "urn:container-alpine-latest-sha256:69665d02cb32192e52e07644d76bc6f25abeb5410edc1c7a81a10ba3f0efb90a-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "to": "urn:openssl-3.0.4-4fe40e24-20e3-11ee-be56-0242ac120002",
			 "relationshipType": "contains"
		 }
	 ]
}
`
