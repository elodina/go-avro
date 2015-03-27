/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package main

import (
	//    "github.com/linkedin/goavro"
	"fmt"
	//    "bytes"
	//    "bufio"
	//    "encoding/json"
	"github.com/stealthly/go-avro"
)

var schema0 = `{
  "type": "record",
  "name": "Employee",
  "fields": [
      {"name": "name", "type": "string"},
      {"name": "age", "type": "int"},
      {"name": "emails", "type": {"type": "array", "items": "string"}},
      {"name": "boss", "type": ["Employee","null"]}
  ]
}`

var schema1 = `{"namespace": "example.avro",
 "type": "record",
 "name": "User",
 "fields": [
     {"name": "name", "type": "string"},
     {"name": "favorite_number",  "type": ["int", "null"]},
     {"name": "favorite_color", "type": ["string", "null"]}
 ]
}`

var schema2 = `{
     "type": "record",
     "namespace": "com.example",
     "name": "FullName",
     "fields": [
       { "name": "first", "type": "string" },
       { "name": "last", "type": "string" }
     ]
}`

var schema3 = `{
    "type" : "record",
    "name" : "userInfo",
    "namespace" : "my.example",
    "fields" : [{"name" : "age", "type" : "int"}]
} `

var schema4 = `{
    "type" : "record",
    "name" : "userInfo",
    "namespace" : "my.example",
    "fields" : [{"name" : "age", "type" : "int", "default" : -1}]
}`

var schema5 = `{
    "type" : "record",
    "name" : "userInfo",
    "namespace" : "my.example",
    "fields" : [{"name" : "username",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "age",
                "type" : "int",
                "default" : -1},

                {"name" : "phone",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "housenum",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "street",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "city",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "state_province",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "country",
                "type" : "string",
                "default" : "NONE"},

                {"name" : "zip",
                "type" : "string",
                "default" : "NONE"}]
} `

var schema6 = `{
    "type" : "record",
    "name" : "userInfo",
    "namespace" : "my.example",
    "fields" : [{"name" : "username",
                 "type" : "string",
                 "default" : "NONE"},

                {"name" : "age",
                 "type" : "int",
                 "default" : -1},

                 {"name" : "phone",
                  "type" : "string",
                  "default" : "NONE"},

                 {"name" : "housenum",
                  "type" : "string",
                  "default" : "NONE"},

                  {"name" : "address",
                     "type" : {
                         "type" : "record",
                         "name" : "mailing_address",
                         "fields" : [
                            {"name" : "street",
                             "type" : "string",
                             "default" : "NONE"},

                            {"name" : "city",
                             "type" : "string",
                             "default" : "NONE"},

                            {"name" : "state_prov",
                             "type" : "string",
                             "default" : "NONE"},

                            {"name" : "country",
                             "type" : "string",
                             "default" : "NONE"},

                            {"name" : "zip",
                             "type" : "string",
                             "default" : "NONE"}
                            ]}
                }
    ]
} `

var schema7 = `{ "type" : "enum",
  "name" : "Colors",
  "namespace" : "palette",
  "doc" : "Colors supported by the palette.",
  "symbols" : ["WHITE", "BLUE", "GREEN", "RED", "BLACK"]}`

var schema8 = `{"type" : "array", "items" : "string"}`

var schema9 = `{"type" : "map", "values" : "int"}`

var schema10 = `["string", "null"]`

var schema11 = `{
     "type": "record",
     "namespace": "com.example",
     "name": "FullName",
     "fields": [
       { "name": "first", "type": ["string", "null"] },
       { "name": "last", "type": "string", "default" : "Doe" }
     ]
} `

var schema12 = `{"type" : "fixed" , "name" : "bdata", "size" : 1048576}`

var schema13 = `{
  "type" : "record",
  "name" : "twitter_schema",
  "namespace" : "com.miguno.avro",
  "fields" : [ {
    "name" : "username",
    "type" : "string",
    "doc"  : "Name of the user account on Twitter.com"
  }, {
    "name" : "tweet",
    "type" : "string",
    "doc"  : "The content of the user's Twitter message"
  }, {
    "name" : "timestamp",
    "type" : "long",
    "doc"  : "Unix epoch time in seconds"
  } ],
  "doc:" : "A basic schema for storing Twitter messages"
}`

var schema14 = `{
    "namespace": "com.rishav.avro",
    "type": "record",
    "name": "StudentActivity",
    "fields": [
        {
            "name": "id",
            "type": "string"
        },
        {
            "name": "student_id",
            "type": "int"
        },
        {
            "name": "university_id",
            "type": "int"
        },
        {
            "name": "course_details",
            "type": {
                "name": "Activity",
                "type": "record",
                "fields": [
                    {
                        "name": "course_id",
                        "type": "int"
                    },
                    {
                        "name": "enroll_date",
                        "type": "string"
                    },
                    {
                        "name": "verb",
                        "type": "string"
                    },
                    {
                        "name": "result_score",
                        "type": "double"
                    }
                ]
            }
        }
    ]
}`

var schema15 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "ForceAuthn",
    "type" : [ "boolean", "null" ],
    "source" : "attribute ForceAuthn"
  }, {
    "name" : "IsPassive",
    "type" : [ "boolean", "null" ],
    "source" : "attribute IsPassive"
  }, {
    "name" : "ProtocolBinding",
    "type" : [ "string", "null" ],
    "source" : "attribute ProtocolBinding"
  }, {
    "name" : "AssertionConsumerServiceIndex",
    "type" : [ "string", "null" ],
    "source" : "attribute AssertionConsumerServiceIndex"
  }, {
    "name" : "AssertionConsumerServiceURL",
    "type" : [ "string", "null" ],
    "source" : "attribute AssertionConsumerServiceURL"
  }, {
    "name" : "AttributeConsumingServiceIndex",
    "type" : [ "string", "null" ],
    "source" : "attribute AttributeConsumingServiceIndex"
  }, {
    "name" : "ProviderName",
    "type" : [ "string", "null" ],
    "source" : "attribute ProviderName"
  }, {
    "name" : "ID",
    "type" : "string",
    "source" : "attribute ID"
  }, {
    "name" : "Version",
    "type" : "string",
    "source" : "attribute Version"
  }, {
    "name" : "IssueInstant",
    "type" : "string",
    "source" : "attribute IssueInstant"
  }, {
    "name" : "Destination",
    "type" : [ "string", "null" ],
    "source" : "attribute Destination"
  }, {
    "name" : "Consent",
    "type" : [ "string", "null" ],
    "source" : "attribute Consent"
  }, {
    "name" : "Issuer",
    "type" : [ {
      "type" : "record",
      "name" : "NameIDType",
      "fields" : [ {
        "name" : "NameQualifier",
        "type" : [ "string", "null" ],
        "source" : "attribute NameQualifier"
      }, {
        "name" : "SPNameQualifier",
        "type" : [ "string", "null" ],
        "source" : "attribute SPNameQualifier"
      }, {
        "name" : "Format",
        "type" : [ "string", "null" ],
        "source" : "attribute Format"
      }, {
        "name" : "SPProvidedID",
        "type" : [ "string", "null" ],
        "source" : "attribute SPProvidedID"
      } ]
    }, "null" ],
    "source" : "element Issuer"
  }, {
    "name" : "Signature",
    "type" : [ {
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "Id",
        "type" : [ "string", "null" ],
        "source" : "attribute Id"
      }, {
        "name" : "SignedInfo",
        "type" : {
          "type" : "record",
          "name" : "SignedInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "CanonicalizationMethod",
            "type" : {
              "type" : "record",
              "name" : "CanonicalizationMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element CanonicalizationMethod"
          }, {
            "name" : "SignatureMethod",
            "type" : {
              "type" : "record",
              "name" : "SignatureMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "HMACOutputLength",
                "type" : [ "string", "null" ],
                "source" : "element HMACOutputLength"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element SignatureMethod"
          }, {
            "name" : "Reference",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "ReferenceType",
                "fields" : [ {
                  "name" : "Id",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Id"
                }, {
                  "name" : "URI",
                  "type" : [ "string", "null" ],
                  "source" : "attribute URI"
                }, {
                  "name" : "Type",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Type"
                }, {
                  "name" : "Transforms",
                  "type" : [ {
                    "type" : "record",
                    "name" : "TransformsType",
                    "fields" : [ {
                      "name" : "Transform",
                      "type" : {
                        "type" : "array",
                        "items" : {
                          "type" : "record",
                          "name" : "TransformType",
                          "fields" : [ {
                            "name" : "Algorithm",
                            "type" : "string",
                            "source" : "attribute Algorithm"
                          }, {
                            "name" : "others",
                            "type" : {
                              "type" : "map",
                              "values" : "string"
                            }
                          }, {
                            "name" : "XPath",
                            "type" : [ "string", "null" ],
                            "source" : "element XPath"
                          } ]
                        }
                      },
                      "source" : "element Transform"
                    } ]
                  }, "null" ],
                  "source" : "element Transforms"
                }, {
                  "name" : "DigestMethod",
                  "type" : {
                    "type" : "record",
                    "name" : "DigestMethodType",
                    "fields" : [ {
                      "name" : "Algorithm",
                      "type" : "string",
                      "source" : "attribute Algorithm"
                    }, {
                      "name" : "others",
                      "type" : {
                        "type" : "map",
                        "values" : "string"
                      }
                    } ]
                  },
                  "source" : "element DigestMethod"
                }, {
                  "name" : "DigestValue",
                  "type" : "string",
                  "source" : "element DigestValue"
                } ]
              }
            },
            "source" : "element Reference"
          } ]
        },
        "source" : "element SignedInfo"
      }, {
        "name" : "SignatureValue",
        "type" : {
          "type" : "record",
          "name" : "SignatureValueType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        },
        "source" : "element SignatureValue"
      }, {
        "name" : "KeyInfo",
        "type" : [ {
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          }, {
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }, {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          }, {
            "name" : "X509Data",
            "type" : [ {
              "type" : "record",
              "name" : "X509DataType",
              "fields" : [ {
                "name" : "X509IssuerSerial",
                "type" : [ {
                  "type" : "record",
                  "name" : "X509IssuerSerialType",
                  "fields" : [ {
                    "name" : "X509IssuerName",
                    "type" : "string",
                    "source" : "element X509IssuerName"
                  }, {
                    "name" : "X509SerialNumber",
                    "type" : "string",
                    "source" : "element X509SerialNumber"
                  } ]
                }, "null" ],
                "source" : "element X509IssuerSerial"
              }, {
                "name" : "X509SKI",
                "type" : [ "string", "null" ],
                "source" : "element X509SKI"
              }, {
                "name" : "X509SubjectName",
                "type" : [ "string", "null" ],
                "source" : "element X509SubjectName"
              }, {
                "name" : "X509Certificate",
                "type" : [ "string", "null" ],
                "source" : "element X509Certificate"
              }, {
                "name" : "X509CRL",
                "type" : [ "string", "null" ],
                "source" : "element X509CRL"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element X509Data"
          }, {
            "name" : "PGPData",
            "type" : [ {
              "type" : "record",
              "name" : "PGPDataType",
              "fields" : [ {
                "name" : "PGPKeyID",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyID"
              }, {
                "name" : "PGPKeyPacket0",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyPacket"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element PGPData"
          }, {
            "name" : "SPKIData",
            "type" : [ {
              "type" : "record",
              "name" : "SPKIDataType",
              "fields" : [ {
                "name" : "SPKISexp",
                "type" : "string",
                "source" : "element SPKISexp"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element SPKIData"
          }, {
            "name" : "MgmtData",
            "type" : [ "string", "null" ],
            "source" : "element MgmtData"
          }, {
            "name" : "others",
            "type" : {
              "type" : "map",
              "values" : "string"
            }
          } ]
        }, "null" ],
        "source" : "element KeyInfo"
      }, {
        "name" : "Object",
        "type" : {
          "type" : "array",
          "items" : {
            "type" : "record",
            "name" : "ObjectType",
            "fields" : [ {
              "name" : "Id",
              "type" : [ "string", "null" ],
              "source" : "attribute Id"
            }, {
              "name" : "MimeType",
              "type" : [ "string", "null" ],
              "source" : "attribute MimeType"
            }, {
              "name" : "Encoding",
              "type" : [ "string", "null" ],
              "source" : "attribute Encoding"
            }, {
              "name" : "others",
              "type" : {
                "type" : "map",
                "values" : "string"
              }
            } ]
          }
        },
        "source" : "element Object"
      } ]
    }, "null" ],
    "source" : "element Signature"
  }, {
    "name" : "Extensions",
    "type" : [ {
      "type" : "record",
      "name" : "ExtensionsType",
      "fields" : [ {
        "name" : "others",
        "type" : {
          "type" : "map",
          "values" : "string"
        }
      } ]
    }, "null" ],
    "source" : "element Extensions"
  }, {
    "name" : "Subject",
    "type" : [ {
      "type" : "record",
      "name" : "SubjectType",
      "fields" : [ {
        "name" : "BaseID",
        "type" : [ {
          "type" : "record",
          "name" : "BaseIDAbstractType",
          "fields" : [ {
            "name" : "NameQualifier",
            "type" : [ "string", "null" ],
            "source" : "attribute NameQualifier"
          }, {
            "name" : "SPNameQualifier",
            "type" : [ "string", "null" ],
            "source" : "attribute SPNameQualifier"
          } ]
        }, "null" ],
        "source" : "element BaseID"
      }, {
        "name" : "NameID",
        "type" : [ "NameIDType", "null" ],
        "source" : "element NameID"
      }, {
        "name" : "EncryptedID",
        "type" : [ {
          "type" : "record",
          "name" : "EncryptedElementType",
          "fields" : [ {
            "name" : "EncryptedData",
            "type" : {
              "type" : "record",
              "name" : "EncryptedDataType",
              "fields" : [ {
                "name" : "Id",
                "type" : [ "string", "null" ],
                "source" : "attribute Id"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "MimeType",
                "type" : [ "string", "null" ],
                "source" : "attribute MimeType"
              }, {
                "name" : "Encoding",
                "type" : [ "string", "null" ],
                "source" : "attribute Encoding"
              }, {
                "name" : "EncryptionMethod",
                "type" : [ {
                  "type" : "record",
                  "name" : "EncryptionMethodType",
                  "fields" : [ {
                    "name" : "Algorithm",
                    "type" : "string",
                    "source" : "attribute Algorithm"
                  }, {
                    "name" : "KeySize",
                    "type" : [ "string", "null" ],
                    "source" : "element KeySize"
                  }, {
                    "name" : "OAEPparams",
                    "type" : [ "string", "null" ],
                    "source" : "element OAEPparams"
                  }, {
                    "name" : "others",
                    "type" : {
                      "type" : "map",
                      "values" : "string"
                    }
                  } ]
                }, "null" ],
                "source" : "element EncryptionMethod"
              }, {
                "name" : "KeyInfo",
                "type" : [ "KeyInfoType", "null" ],
                "source" : "element KeyInfo"
              }, {
                "name" : "CipherData",
                "type" : {
                  "type" : "record",
                  "name" : "CipherDataType",
                  "fields" : [ {
                    "name" : "CipherValue",
                    "type" : [ "string", "null" ],
                    "source" : "element CipherValue"
                  }, {
                    "name" : "CipherReference",
                    "type" : [ {
                      "type" : "record",
                      "name" : "CipherReferenceType",
                      "fields" : [ {
                        "name" : "URI",
                        "type" : "string",
                        "source" : "attribute URI"
                      }, {
                        "name" : "Transforms",
                        "type" : [ "TransformsType", "null" ],
                        "source" : "element Transforms"
                      } ]
                    }, "null" ],
                    "source" : "element CipherReference"
                  } ]
                },
                "source" : "element CipherData"
              }, {
                "name" : "EncryptionProperties",
                "type" : [ {
                  "type" : "record",
                  "name" : "EncryptionPropertiesType",
                  "fields" : [ {
                    "name" : "Id",
                    "type" : [ "string", "null" ],
                    "source" : "attribute Id"
                  }, {
                    "name" : "EncryptionProperty",
                    "type" : {
                      "type" : "array",
                      "items" : {
                        "type" : "record",
                        "name" : "EncryptionPropertyType",
                        "fields" : [ {
                          "name" : "Target",
                          "type" : [ "string", "null" ],
                          "source" : "attribute Target"
                        }, {
                          "name" : "Id",
                          "type" : [ "string", "null" ],
                          "source" : "attribute Id"
                        }, {
                          "name" : "others",
                          "type" : {
                            "type" : "map",
                            "values" : "string"
                          }
                        } ]
                      }
                    },
                    "source" : "element EncryptionProperty"
                  } ]
                }, "null" ],
                "source" : "element EncryptionProperties"
              } ]
            },
            "source" : "element EncryptedData"
          }, {
            "name" : "EncryptedKey",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "EncryptedKeyType",
                "fields" : [ {
                  "name" : "Recipient",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Recipient"
                }, {
                  "name" : "Id",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Id"
                }, {
                  "name" : "Type",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Type"
                }, {
                  "name" : "MimeType",
                  "type" : [ "string", "null" ],
                  "source" : "attribute MimeType"
                }, {
                  "name" : "Encoding",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Encoding"
                }, {
                  "name" : "EncryptionMethod",
                  "type" : [ "EncryptionMethodType", "null" ],
                  "source" : "element EncryptionMethod"
                }, {
                  "name" : "KeyInfo",
                  "type" : [ "KeyInfoType", "null" ],
                  "source" : "element KeyInfo"
                }, {
                  "name" : "CipherData",
                  "type" : "CipherDataType",
                  "source" : "element CipherData"
                }, {
                  "name" : "EncryptionProperties",
                  "type" : [ "EncryptionPropertiesType", "null" ],
                  "source" : "element EncryptionProperties"
                }, {
                  "name" : "ReferenceList",
                  "type" : [ {
                    "type" : "record",
                    "name" : "type0",
                    "fields" : [ {
                      "name" : "DataReference",
                      "type" : [ "ReferenceType", "null" ],
                      "source" : "element DataReference"
                    }, {
                      "name" : "KeyReference",
                      "type" : [ "ReferenceType", "null" ],
                      "source" : "element KeyReference"
                    } ]
                  }, "null" ],
                  "source" : "element ReferenceList"
                }, {
                  "name" : "CarriedKeyName",
                  "type" : [ "string", "null" ],
                  "source" : "element CarriedKeyName"
                } ]
              }
            },
            "source" : "element EncryptedKey"
          } ]
        }, "null" ],
        "source" : "element EncryptedID"
      }, {
        "name" : "SubjectConfirmation0",
        "type" : {
          "type" : "array",
          "items" : {
            "type" : "record",
            "name" : "SubjectConfirmationType",
            "fields" : [ {
              "name" : "Method",
              "type" : "string",
              "source" : "attribute Method"
            }, {
              "name" : "BaseID",
              "type" : [ "BaseIDAbstractType", "null" ],
              "source" : "element BaseID"
            }, {
              "name" : "NameID",
              "type" : [ "NameIDType", "null" ],
              "source" : "element NameID"
            }, {
              "name" : "EncryptedID",
              "type" : [ "EncryptedElementType", "null" ],
              "source" : "element EncryptedID"
            }, {
              "name" : "SubjectConfirmationData",
              "type" : [ {
                "type" : "record",
                "name" : "SubjectConfirmationDataType",
                "fields" : [ {
                  "name" : "NotBefore",
                  "type" : [ "string", "null" ],
                  "source" : "attribute NotBefore"
                }, {
                  "name" : "NotOnOrAfter",
                  "type" : [ "string", "null" ],
                  "source" : "attribute NotOnOrAfter"
                }, {
                  "name" : "Recipient",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Recipient"
                }, {
                  "name" : "InResponseTo",
                  "type" : [ "string", "null" ],
                  "source" : "attribute InResponseTo"
                }, {
                  "name" : "Address",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Address"
                }, {
                  "name" : "others",
                  "type" : {
                    "type" : "map",
                    "values" : "string"
                  }
                } ]
              }, "null" ],
              "source" : "element SubjectConfirmationData"
            } ]
          }
        },
        "source" : "element SubjectConfirmation"
      } ]
    }, "null" ],
    "source" : "element Subject"
  }, {
    "name" : "NameIDPolicy",
    "type" : [ {
      "type" : "record",
      "name" : "NameIDPolicyType",
      "fields" : [ {
        "name" : "Format",
        "type" : [ "string", "null" ],
        "source" : "attribute Format"
      }, {
        "name" : "SPNameQualifier",
        "type" : [ "string", "null" ],
        "source" : "attribute SPNameQualifier"
      }, {
        "name" : "AllowCreate",
        "type" : [ "boolean", "null" ],
        "source" : "attribute AllowCreate"
      } ]
    }, "null" ],
    "source" : "element NameIDPolicy"
  }, {
    "name" : "Conditions",
    "type" : [ {
      "type" : "record",
      "name" : "ConditionsType",
      "fields" : [ {
        "name" : "NotBefore",
        "type" : [ "string", "null" ],
        "source" : "attribute NotBefore"
      }, {
        "name" : "NotOnOrAfter",
        "type" : [ "string", "null" ],
        "source" : "attribute NotOnOrAfter"
      }, {
        "name" : "Condition",
        "type" : [ {
          "type" : "record",
          "name" : "ConditionAbstractType",
          "fields" : [ ]
        }, "null" ],
        "source" : "element Condition"
      }, {
        "name" : "AudienceRestriction",
        "type" : [ {
          "type" : "record",
          "name" : "AudienceRestrictionType",
          "fields" : [ {
            "name" : "Audience",
            "type" : {
              "type" : "array",
              "items" : "string"
            },
            "source" : "element Audience"
          } ]
        }, "null" ],
        "source" : "element AudienceRestriction"
      }, {
        "name" : "OneTimeUse",
        "type" : [ {
          "type" : "record",
          "name" : "OneTimeUseType",
          "fields" : [ ]
        }, "null" ],
        "source" : "element OneTimeUse"
      }, {
        "name" : "ProxyRestriction",
        "type" : [ {
          "type" : "record",
          "name" : "ProxyRestrictionType",
          "fields" : [ {
            "name" : "Count",
            "type" : [ "string", "null" ],
            "source" : "attribute Count"
          }, {
            "name" : "Audience",
            "type" : {
              "type" : "array",
              "items" : "string"
            },
            "source" : "element Audience"
          } ]
        }, "null" ],
        "source" : "element ProxyRestriction"
      } ]
    }, "null" ],
    "source" : "element Conditions"
  }, {
    "name" : "RequestedAuthnContext",
    "type" : [ {
      "type" : "record",
      "name" : "RequestedAuthnContextType",
      "fields" : [ {
        "name" : "Comparison",
        "type" : [ "string", "null" ],
        "source" : "attribute Comparison"
      }, {
        "name" : "AuthnContextClassRef",
        "type" : {
          "type" : "array",
          "items" : "string"
        },
        "source" : "element AuthnContextClassRef"
      }, {
        "name" : "AuthnContextDeclRef",
        "type" : {
          "type" : "array",
          "items" : "string"
        },
        "source" : "element AuthnContextDeclRef"
      } ]
    }, "null" ],
    "source" : "element RequestedAuthnContext"
  }, {
    "name" : "Scoping",
    "type" : [ {
      "type" : "record",
      "name" : "ScopingType",
      "fields" : [ {
        "name" : "ProxyCount",
        "type" : [ "string", "null" ],
        "source" : "attribute ProxyCount"
      }, {
        "name" : "IDPList",
        "type" : [ {
          "type" : "record",
          "name" : "IDPListType",
          "fields" : [ {
            "name" : "IDPEntry",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "IDPEntryType",
                "fields" : [ {
                  "name" : "ProviderID",
                  "type" : "string",
                  "source" : "attribute ProviderID"
                }, {
                  "name" : "Name",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Name"
                }, {
                  "name" : "Loc",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Loc"
                } ]
              }
            },
            "source" : "element IDPEntry"
          }, {
            "name" : "GetComplete",
            "type" : [ "string", "null" ],
            "source" : "element GetComplete"
          } ]
        }, "null" ],
        "source" : "element IDPList"
      }, {
        "name" : "RequesterID",
        "type" : {
          "type" : "array",
          "items" : "string"
        },
        "source" : "element RequesterID"
      } ]
    }, "null" ],
    "source" : "element Scoping"
  } ]
}`

var schema16 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "ForceAuthn",
    "type" : [ "boolean", "null" ],
    "source" : "attribute ForceAuthn"
  } ]
}`

var schema17 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "IsPassive",
    "type" : [ "boolean", "null" ],
    "source" : "attribute IsPassive"
  } ]
}`

var schema18 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "ProtocolBinding",
    "type" : [ "string", "null" ],
    "source" : "attribute ProtocolBinding"
  } ]
}`

var schema19 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "AssertionConsumerServiceIndex",
    "type" : [ "string", "null" ],
    "source" : "attribute AssertionConsumerServiceIndex"
  } ]
}`

var schema20 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "AssertionConsumerServiceURL",
    "type" : [ "string", "null" ],
    "source" : "attribute AssertionConsumerServiceURL"
  } ]
}`

var schema21 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "AttributeConsumingServiceIndex",
    "type" : [ "string", "null" ],
    "source" : "attribute AttributeConsumingServiceIndex"
  } ]
}`

var schema22 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "ProviderName",
    "type" : [ "string", "null" ],
    "source" : "attribute ProviderName"
  } ]
}`

var schema23 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "Version",
    "type" : "string",
    "source" : "attribute Version"
  } ]
}`

var schema24 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "IssueInstant",
    "type" : "string",
    "source" : "attribute IssueInstant"
  } ]
}`

var schema25 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "Destination",
    "type" : [ "string", "null" ],
    "source" : "attribute Destination"
  } ]
}`

var schema26 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "Consent",
    "type" : [ "string", "null" ],
    "source" : "attribute Consent"
  } ]
}`

var schema27 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "Issuer",
    "type" : [ {
      "type" : "record",
      "name" : "NameIDType",
      "fields" : [ {
        "name" : "NameQualifier",
        "type" : [ "string", "null" ],
        "source" : "attribute NameQualifier"
      }, {
        "name" : "SPNameQualifier",
        "type" : [ "string", "null" ],
        "source" : "attribute SPNameQualifier"
      }, {
        "name" : "Format",
        "type" : [ "string", "null" ],
        "source" : "attribute Format"
      }, {
        "name" : "SPProvidedID",
        "type" : [ "string", "null" ],
        "source" : "attribute SPProvidedID"
      } ]
    }, "null" ],
    "source" : "element Issuer"
  } ]
}`

var schema28 = `{
  "type" : "record",
  "name" : "AuthnRequestType",
  "fields" : [ {
    "name" : "Signature",
    "type" : [ {
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "Id",
        "type" : [ "string", "null" ],
        "source" : "attribute Id"
      }, {
        "name" : "SignedInfo",
        "type" : {
          "type" : "record",
          "name" : "SignedInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "CanonicalizationMethod",
            "type" : {
              "type" : "record",
              "name" : "CanonicalizationMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element CanonicalizationMethod"
          }, {
            "name" : "SignatureMethod",
            "type" : {
              "type" : "record",
              "name" : "SignatureMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "HMACOutputLength",
                "type" : [ "string", "null" ],
                "source" : "element HMACOutputLength"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element SignatureMethod"
          }, {
            "name" : "Reference",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "ReferenceType",
                "fields" : [ {
                  "name" : "Id",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Id"
                }, {
                  "name" : "URI",
                  "type" : [ "string", "null" ],
                  "source" : "attribute URI"
                }, {
                  "name" : "Type",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Type"
                }, {
                  "name" : "Transforms",
                  "type" : [ {
                    "type" : "record",
                    "name" : "TransformsType",
                    "fields" : [ {
                      "name" : "Transform",
                      "type" : {
                        "type" : "array",
                        "items" : {
                          "type" : "record",
                          "name" : "TransformType",
                          "fields" : [ {
                            "name" : "Algorithm",
                            "type" : "string",
                            "source" : "attribute Algorithm"
                          }, {
                            "name" : "others",
                            "type" : {
                              "type" : "map",
                              "values" : "string"
                            }
                          }, {
                            "name" : "XPath",
                            "type" : [ "string", "null" ],
                            "source" : "element XPath"
                          } ]
                        }
                      },
                      "source" : "element Transform"
                    } ]
                  }, "null" ],
                  "source" : "element Transforms"
                }, {
                  "name" : "DigestMethod",
                  "type" : {
                    "type" : "record",
                    "name" : "DigestMethodType",
                    "fields" : [ {
                      "name" : "Algorithm",
                      "type" : "string",
                      "source" : "attribute Algorithm"
                    }, {
                      "name" : "others",
                      "type" : {
                        "type" : "map",
                        "values" : "string"
                      }
                    } ]
                  },
                  "source" : "element DigestMethod"
                }, {
                  "name" : "DigestValue",
                  "type" : "string",
                  "source" : "element DigestValue"
                } ]
              }
            },
            "source" : "element Reference"
          } ]
        },
        "source" : "element SignedInfo"
      }, {
        "name" : "SignatureValue",
        "type" : {
          "type" : "record",
          "name" : "SignatureValueType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        },
        "source" : "element SignatureValue"
      }, {
        "name" : "KeyInfo",
        "type" : [ {
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          }, {
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }, {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          }, {
            "name" : "X509Data",
            "type" : [ {
              "type" : "record",
              "name" : "X509DataType",
              "fields" : [ {
                "name" : "X509IssuerSerial",
                "type" : [ {
                  "type" : "record",
                  "name" : "X509IssuerSerialType",
                  "fields" : [ {
                    "name" : "X509IssuerName",
                    "type" : "string",
                    "source" : "element X509IssuerName"
                  }, {
                    "name" : "X509SerialNumber",
                    "type" : "string",
                    "source" : "element X509SerialNumber"
                  } ]
                }, "null" ],
                "source" : "element X509IssuerSerial"
              }, {
                "name" : "X509SKI",
                "type" : [ "string", "null" ],
                "source" : "element X509SKI"
              }, {
                "name" : "X509SubjectName",
                "type" : [ "string", "null" ],
                "source" : "element X509SubjectName"
              }, {
                "name" : "X509Certificate",
                "type" : [ "string", "null" ],
                "source" : "element X509Certificate"
              }, {
                "name" : "X509CRL",
                "type" : [ "string", "null" ],
                "source" : "element X509CRL"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element X509Data"
          }, {
            "name" : "PGPData",
            "type" : [ {
              "type" : "record",
              "name" : "PGPDataType",
              "fields" : [ {
                "name" : "PGPKeyID",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyID"
              }, {
                "name" : "PGPKeyPacket0",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyPacket"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element PGPData"
          }, {
            "name" : "SPKIData",
            "type" : [ {
              "type" : "record",
              "name" : "SPKIDataType",
              "fields" : [ {
                "name" : "SPKISexp",
                "type" : "string",
                "source" : "element SPKISexp"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element SPKIData"
          }, {
            "name" : "MgmtData",
            "type" : [ "string", "null" ],
            "source" : "element MgmtData"
          }, {
            "name" : "others",
            "type" : {
              "type" : "map",
              "values" : "string"
            }
          } ]
        }, "null" ],
        "source" : "element KeyInfo"
      }, {
        "name" : "Object",
        "type" : {
          "type" : "array",
          "items" : {
            "type" : "record",
            "name" : "ObjectType",
            "fields" : [ {
              "name" : "Id",
              "type" : [ "string", "null" ],
              "source" : "attribute Id"
            }, {
              "name" : "MimeType",
              "type" : [ "string", "null" ],
              "source" : "attribute MimeType"
            }, {
              "name" : "Encoding",
              "type" : [ "string", "null" ],
              "source" : "attribute Encoding"
            }, {
              "name" : "others",
              "type" : {
                "type" : "map",
                "values" : "string"
              }
            } ]
          }
        },
        "source" : "element Object"
      } ]
    }, "null" ],
    "source" : "element Signature"
  } ]
}`

var schema29 = `{
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "Id",
        "type" : [ "string", "null" ],
        "source" : "attribute Id"
      }, {
        "name" : "SignedInfo",
        "type" : {
          "type" : "record",
          "name" : "SignedInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "CanonicalizationMethod",
            "type" : {
              "type" : "record",
              "name" : "CanonicalizationMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element CanonicalizationMethod"
          }, {
            "name" : "SignatureMethod",
            "type" : {
              "type" : "record",
              "name" : "SignatureMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "HMACOutputLength",
                "type" : [ "string", "null" ],
                "source" : "element HMACOutputLength"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element SignatureMethod"
          }, {
            "name" : "Reference",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "ReferenceType",
                "fields" : [ {
                  "name" : "Id",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Id"
                }, {
                  "name" : "URI",
                  "type" : [ "string", "null" ],
                  "source" : "attribute URI"
                }, {
                  "name" : "Type",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Type"
                }, {
                  "name" : "Transforms",
                  "type" : [ {
                    "type" : "record",
                    "name" : "TransformsType",
                    "fields" : [ {
                      "name" : "Transform",
                      "type" : {
                        "type" : "array",
                        "items" : {
                          "type" : "record",
                          "name" : "TransformType",
                          "fields" : [ {
                            "name" : "Algorithm",
                            "type" : "string",
                            "source" : "attribute Algorithm"
                          }, {
                            "name" : "others",
                            "type" : {
                              "type" : "map",
                              "values" : "string"
                            }
                          }, {
                            "name" : "XPath",
                            "type" : [ "string", "null" ],
                            "source" : "element XPath"
                          } ]
                        }
                      },
                      "source" : "element Transform"
                    } ]
                  }, "null" ],
                  "source" : "element Transforms"
                }, {
                  "name" : "DigestMethod",
                  "type" : {
                    "type" : "record",
                    "name" : "DigestMethodType",
                    "fields" : [ {
                      "name" : "Algorithm",
                      "type" : "string",
                      "source" : "attribute Algorithm"
                    }, {
                      "name" : "others",
                      "type" : {
                        "type" : "map",
                        "values" : "string"
                      }
                    } ]
                  },
                  "source" : "element DigestMethod"
                }, {
                  "name" : "DigestValue",
                  "type" : "string",
                  "source" : "element DigestValue"
                } ]
              }
            },
            "source" : "element Reference"
          } ]
        },
        "source" : "element SignedInfo"
      }, {
        "name" : "SignatureValue",
        "type" : {
          "type" : "record",
          "name" : "SignatureValueType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        },
        "source" : "element SignatureValue"
      }, {
        "name" : "KeyInfo",
        "type" : [ {
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          }, {
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }, {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          }, {
            "name" : "X509Data",
            "type" : [ {
              "type" : "record",
              "name" : "X509DataType",
              "fields" : [ {
                "name" : "X509IssuerSerial",
                "type" : [ {
                  "type" : "record",
                  "name" : "X509IssuerSerialType",
                  "fields" : [ {
                    "name" : "X509IssuerName",
                    "type" : "string",
                    "source" : "element X509IssuerName"
                  }, {
                    "name" : "X509SerialNumber",
                    "type" : "string",
                    "source" : "element X509SerialNumber"
                  } ]
                }, "null" ],
                "source" : "element X509IssuerSerial"
              }, {
                "name" : "X509SKI",
                "type" : [ "string", "null" ],
                "source" : "element X509SKI"
              }, {
                "name" : "X509SubjectName",
                "type" : [ "string", "null" ],
                "source" : "element X509SubjectName"
              }, {
                "name" : "X509Certificate",
                "type" : [ "string", "null" ],
                "source" : "element X509Certificate"
              }, {
                "name" : "X509CRL",
                "type" : [ "string", "null" ],
                "source" : "element X509CRL"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element X509Data"
          }, {
            "name" : "PGPData",
            "type" : [ {
              "type" : "record",
              "name" : "PGPDataType",
              "fields" : [ {
                "name" : "PGPKeyID",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyID"
              }, {
                "name" : "PGPKeyPacket0",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyPacket"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element PGPData"
          }, {
            "name" : "SPKIData",
            "type" : [ {
              "type" : "record",
              "name" : "SPKIDataType",
              "fields" : [ {
                "name" : "SPKISexp",
                "type" : "string",
                "source" : "element SPKISexp"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element SPKIData"
          }, {
            "name" : "MgmtData",
            "type" : [ "string", "null" ],
            "source" : "element MgmtData"
          }, {
            "name" : "others",
            "type" : {
              "type" : "map",
              "values" : "string"
            }
          } ]
        }, "null" ],
        "source" : "element KeyInfo"
      }, {
        "name" : "Object",
        "type" : {
          "type" : "array",
          "items" : {
            "type" : "record",
            "name" : "ObjectType",
            "fields" : [ {
              "name" : "Id",
              "type" : [ "string", "null" ],
              "source" : "attribute Id"
            }, {
              "name" : "MimeType",
              "type" : [ "string", "null" ],
              "source" : "attribute MimeType"
            }, {
              "name" : "Encoding",
              "type" : [ "string", "null" ],
              "source" : "attribute Encoding"
            }, {
              "name" : "others",
              "type" : {
                "type" : "map",
                "values" : "string"
              }
            } ]
          }
        },
        "source" : "element Object"
      } ]
    }`

var schema30 = `{
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "Id",
        "type" : [ "string", "null" ],
        "source" : "attribute Id"
      } ]
    }`

var schema31 = `{
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "SignedInfo",
        "type" : {
          "type" : "record",
          "name" : "SignedInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "CanonicalizationMethod",
            "type" : {
              "type" : "record",
              "name" : "CanonicalizationMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element CanonicalizationMethod"
          }, {
            "name" : "SignatureMethod",
            "type" : {
              "type" : "record",
              "name" : "SignatureMethodType",
              "fields" : [ {
                "name" : "Algorithm",
                "type" : "string",
                "source" : "attribute Algorithm"
              }, {
                "name" : "HMACOutputLength",
                "type" : [ "string", "null" ],
                "source" : "element HMACOutputLength"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            },
            "source" : "element SignatureMethod"
          }, {
            "name" : "Reference",
            "type" : {
              "type" : "array",
              "items" : {
                "type" : "record",
                "name" : "ReferenceType",
                "fields" : [ {
                  "name" : "Id",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Id"
                }, {
                  "name" : "URI",
                  "type" : [ "string", "null" ],
                  "source" : "attribute URI"
                }, {
                  "name" : "Type",
                  "type" : [ "string", "null" ],
                  "source" : "attribute Type"
                }, {
                  "name" : "Transforms",
                  "type" : [ {
                    "type" : "record",
                    "name" : "TransformsType",
                    "fields" : [ {
                      "name" : "Transform",
                      "type" : {
                        "type" : "array",
                        "items" : {
                          "type" : "record",
                          "name" : "TransformType",
                          "fields" : [ {
                            "name" : "Algorithm",
                            "type" : "string",
                            "source" : "attribute Algorithm"
                          }, {
                            "name" : "others",
                            "type" : {
                              "type" : "map",
                              "values" : "string"
                            }
                          }, {
                            "name" : "XPath",
                            "type" : [ "string", "null" ],
                            "source" : "element XPath"
                          } ]
                        }
                      },
                      "source" : "element Transform"
                    } ]
                  }, "null" ],
                  "source" : "element Transforms"
                }, {
                  "name" : "DigestMethod",
                  "type" : {
                    "type" : "record",
                    "name" : "DigestMethodType",
                    "fields" : [ {
                      "name" : "Algorithm",
                      "type" : "string",
                      "source" : "attribute Algorithm"
                    }, {
                      "name" : "others",
                      "type" : {
                        "type" : "map",
                        "values" : "string"
                      }
                    } ]
                  },
                  "source" : "element DigestMethod"
                }, {
                  "name" : "DigestValue",
                  "type" : "string",
                  "source" : "element DigestValue"
                } ]
              }
            },
            "source" : "element Reference"
          } ]
        },
        "source" : "element SignedInfo"
      } ]
    }`

var schema32 = `{
      "type" : "record",
      "name" : "SignatureType",
      "fields" : [ {
        "name" : "SignatureValue",
        "type" : {
          "type" : "record",
          "name" : "SignatureValueType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        },
        "source" : "element SignatureValue"
      } ]
    }`

var schema33 = `{
        "name" : "KeyInfo",
        "type" : [ {
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          }, {
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }, {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          }, {
            "name" : "X509Data",
            "type" : [ {
              "type" : "record",
              "name" : "X509DataType",
              "fields" : [ {
                "name" : "X509IssuerSerial",
                "type" : [ {
                  "type" : "record",
                  "name" : "X509IssuerSerialType",
                  "fields" : [ {
                    "name" : "X509IssuerName",
                    "type" : "string",
                    "source" : "element X509IssuerName"
                  }, {
                    "name" : "X509SerialNumber",
                    "type" : "string",
                    "source" : "element X509SerialNumber"
                  } ]
                }, "null" ],
                "source" : "element X509IssuerSerial"
              }, {
                "name" : "X509SKI",
                "type" : [ "string", "null" ],
                "source" : "element X509SKI"
              }, {
                "name" : "X509SubjectName",
                "type" : [ "string", "null" ],
                "source" : "element X509SubjectName"
              }, {
                "name" : "X509Certificate",
                "type" : [ "string", "null" ],
                "source" : "element X509Certificate"
              }, {
                "name" : "X509CRL",
                "type" : [ "string", "null" ],
                "source" : "element X509CRL"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element X509Data"
          }, {
            "name" : "PGPData",
            "type" : [ {
              "type" : "record",
              "name" : "PGPDataType",
              "fields" : [ {
                "name" : "PGPKeyID",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyID"
              }, {
                "name" : "PGPKeyPacket0",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyPacket"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element PGPData"
          }, {
            "name" : "SPKIData",
            "type" : [ {
              "type" : "record",
              "name" : "SPKIDataType",
              "fields" : [ {
                "name" : "SPKISexp",
                "type" : "string",
                "source" : "element SPKISexp"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element SPKIData"
          }, {
            "name" : "MgmtData",
            "type" : [ "string", "null" ],
            "source" : "element MgmtData"
          }, {
            "name" : "others",
            "type" : {
              "type" : "map",
              "values" : "string"
            }
          } ]
        }, "null" ],
        "source" : "element KeyInfo"
      }`

var schema34 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          }, {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          }, {
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }, {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          }, {
            "name" : "X509Data",
            "type" : [ {
              "type" : "record",
              "name" : "X509DataType",
              "fields" : [ {
                "name" : "X509IssuerSerial",
                "type" : [ {
                  "type" : "record",
                  "name" : "X509IssuerSerialType",
                  "fields" : [ {
                    "name" : "X509IssuerName",
                    "type" : "string",
                    "source" : "element X509IssuerName"
                  }, {
                    "name" : "X509SerialNumber",
                    "type" : "string",
                    "source" : "element X509SerialNumber"
                  } ]
                }, "null" ],
                "source" : "element X509IssuerSerial"
              }, {
                "name" : "X509SKI",
                "type" : [ "string", "null" ],
                "source" : "element X509SKI"
              }, {
                "name" : "X509SubjectName",
                "type" : [ "string", "null" ],
                "source" : "element X509SubjectName"
              }, {
                "name" : "X509Certificate",
                "type" : [ "string", "null" ],
                "source" : "element X509Certificate"
              }, {
                "name" : "X509CRL",
                "type" : [ "string", "null" ],
                "source" : "element X509CRL"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element X509Data"
          }, {
            "name" : "PGPData",
            "type" : [ {
              "type" : "record",
              "name" : "PGPDataType",
              "fields" : [ {
                "name" : "PGPKeyID",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyID"
              }, {
                "name" : "PGPKeyPacket0",
                "type" : [ "string", "null" ],
                "source" : "element PGPKeyPacket"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element PGPData"
          }, {
            "name" : "SPKIData",
            "type" : [ {
              "type" : "record",
              "name" : "SPKIDataType",
              "fields" : [ {
                "name" : "SPKISexp",
                "type" : "string",
                "source" : "element SPKISexp"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element SPKIData"
          }, {
            "name" : "MgmtData",
            "type" : [ "string", "null" ],
            "source" : "element MgmtData"
          }, {
            "name" : "others",
            "type" : {
              "type" : "map",
              "values" : "string"
            }
          } ]
        }`

var schema35 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        }`

var schema36 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          } ]
        }`

var schema37 = `{
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }`

var schema38 = `{
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }`

var schema39 = `{
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }`

var schema40 = `{
            "name" : "KeyValue",
            "type" : [ "string", "null" ],
            "source" : "element KeyValue"
          }`

var schema41 = `{
            "name" : "KeyValue",
            "type" : "string",
            "source" : "element KeyValue"
          }`

var schema42 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "Id",
            "type" : [ "string", "null" ],
            "source" : "attribute Id"
          } ]
        }`

var schema43 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "KeyName",
            "type" : [ "string", "null" ],
            "source" : "element KeyName"
          } ]
        }`

var schema44 = `{
            "name" : "KeyValue",
            "type" : [ {
              "type" : "record",
              "name" : "KeyValueType",
              "fields" : [ {
                "name" : "DSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "DSAKeyValueType",
                  "fields" : [ {
                    "name" : "P",
                    "type" : "string",
                    "source" : "element P"
                  }, {
                    "name" : "Q",
                    "type" : "string",
                    "source" : "element Q"
                  }, {
                    "name" : "G",
                    "type" : [ "string", "null" ],
                    "source" : "element G"
                  }, {
                    "name" : "Y",
                    "type" : "string",
                    "source" : "element Y"
                  }, {
                    "name" : "J",
                    "type" : [ "string", "null" ],
                    "source" : "element J"
                  }, {
                    "name" : "Seed",
                    "type" : "string",
                    "source" : "element Seed"
                  }, {
                    "name" : "PgenCounter",
                    "type" : "string",
                    "source" : "element PgenCounter"
                  } ]
                }, "null" ],
                "source" : "element DSAKeyValue"
              }, {
                "name" : "RSAKeyValue",
                "type" : [ {
                  "type" : "record",
                  "name" : "RSAKeyValueType",
                  "fields" : [ {
                    "name" : "Modulus",
                    "type" : "string",
                    "source" : "element Modulus"
                  }, {
                    "name" : "Exponent",
                    "type" : "string",
                    "source" : "element Exponent"
                  } ]
                }, "null" ],
                "source" : "element RSAKeyValue"
              }, {
                "name" : "others",
                "type" : {
                  "type" : "map",
                  "values" : "string"
                }
              } ]
            }, "null" ],
            "source" : "element KeyValue"
          }`

var schema45 = `{
          "type" : "record",
          "name" : "KeyInfoType",
          "fields" : [ {
            "name" : "RetrievalMethod",
            "type" : [ {
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }, "null" ],
            "source" : "element RetrievalMethod"
          } ]
        }`

var schema46 = `{
              "type" : "record",
              "name" : "RetrievalMethodType",
              "fields" : [ {
                "name" : "URI",
                "type" : [ "string", "null" ],
                "source" : "attribute URI"
              }, {
                "name" : "Type",
                "type" : [ "string", "null" ],
                "source" : "attribute Type"
              }, {
                "name" : "Transforms",
                "type" : [ "TransformsType", "null" ],
                "source" : "element Transforms"
              } ]
            }`

func main() {
	//    codec, err := goavro.NewCodec(schema15)
	//    if err != nil {
	//        panic(err)
	//    }
	//
	//    fmt.Println(codec)

	//    buffer := &bytes.Buffer{}
	//    writer := bufio.NewWriter(buffer)

	//    type User struct {
	//        Name string
	//        FavoriteNumber int32
	//        FavoriteColor string
	//    }
	//
	//    user := &User{
	//        Name: "asdasd",
	//        FavoriteNumber: 32,
	//        FavoriteColor: "green",
	//    }
	//
	//
	//    err = codec.Encode(writer, user)
	//    if err != nil {
	//        panic(err)
	//    }

	//    sch := make(map[string]interface{})
	//
	//    err = json.Unmarshal([]byte(schema0), &sch)
	//    if err != nil {
	//        panic(err)
	//    }
	//
	//    record, err := goavro.NewRecord(sch)
	//    if err != nil {
	//        panic(err)
	//    }
	//
	//    fmt.Println(record.Fields[0].Datum)

	schema, err := avro.ParseSchema(schema15)
	if err != nil {
		panic(err)
	}

	fmt.Println(schema)
	//    _ = schema
	//    fmt.Println("ok")

}
