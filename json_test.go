package obfuscate

import (
	"encoding/json"
	"fmt"
	"testing"
)

const inputJson = `{
  "string": "string\"int",
  "int": 123456,
  "float": 1234.56,
  "booleanTrue": true,
  "booleanFalse": false,
  "null": null,
  "object": {
    "string": "string\"int",
    "int": 123456,
    "float": 1234.56,
    "booleanTrue": true,
    "booleanFalse": false,
    "null": null,
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ]
  },
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": 123456
    }
  ],
  "notMatchedString": "123456",
  "notMatchedInt": 123456,
  "notMatchedFloat": 1234.56,
  "notMatchedBooleanTrue": true,
  "notMatchedBooleanFalse": false,
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "notMatchedString": "123456",
    "notMatchedInt": 123456,
    "notMatchedFloat": 1234.56,
    "notMatchedBooleanTrue": true,
    "notMatchedBooleanFalse": false,
    "nonMatchedNull": null
  },
  "nested": [
    {
      "string": "string\"int",
      "int": 123456,
      "float": 1234.56,
      "booleanTrue": true,
      "booleanFalse": false,
      "null": null,
      "object": {
        "string": "string\"int",
        "int": 123456,
        "float": 1234.56,
        "booleanTrue": true,
        "booleanFalse": false,
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ]
      },
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": 123456
        }
      ],
      "notMatchedString": "123456",
      "notMatchedInt": 123456,
      "notMatchedFloat": 1234.56,
      "notMatchedBooleanTrue": true,
      "notMatchedBooleanFalse": false,
      "nonMatchedNull": null
    }
  ]
}`

func TestJSONDefaultErrorStrategy(t *testing.T) {
	obfuscator := JSON().Build()

	assertEqual(t, OnErrorLog, obfuscator.onError)
}

func TestJSONObfuscatorWithDefaultSettings(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithExcludeForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), &JSONPropertyObfuscationOptions{ForObjects: Exclude}).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithExcludeAllForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), &JSONPropertyObfuscationOptions{ForObjects: ExcludeAll}).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": false,
        "booleanTrue": true,
        "float": 1234.56,
        "int": 123456,
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "string\"int"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": false,
    "booleanTrue": true,
    "float": 1234.56,
    "int": 123456,
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": null,
    "string": "string\"int"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithInheritForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), &JSONPropertyObfuscationOptions{ForObjects: Inherit}).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "ooo",
        "booleanTrue": "ooo",
        "float": "ooo",
        "int": "ooo",
        "nested": [
          {
            "prop1": "ooo",
            "prop2": "ooo"
          }
        ],
        "string": "ooo"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "ooo",
    "booleanTrue": "ooo",
    "float": "ooo",
    "int": "ooo",
    "nested": [
      {
        "prop1": "ooo",
        "prop2": "ooo"
      }
    ],
    "null": "ooo",
    "string": "ooo"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithInheritOverridableForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), &JSONPropertyObfuscationOptions{ForObjects: InheritOverridable}).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "ooo",
            "prop2": "ooo"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "ooo",
        "prop2": "ooo"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithExcludeForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), &JSONPropertyObfuscationOptions{ForArrays: Exclude}).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithExcludeAllForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), &JSONPropertyObfuscationOptions{ForArrays: ExcludeAll}).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": 123456
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": 123456
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithInheritForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), &JSONPropertyObfuscationOptions{ForArrays: Inherit}).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "aaa",
      "aaa"
    ],
    {},
    {
      "int": "aaa"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "aaa",
          "aaa"
        ],
        {},
        {
          "int": "aaa"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithInheritOverridableForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), &JSONPropertyObfuscationOptions{ForArrays: InheritOverridable}).
		WithProperty("null", obfuscator, nil).
		Build()

	expected := `{
  "array": [
    [
      "aaa",
      "aaa"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "aaa",
          "aaa"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultExcludeForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForObjects(Exclude).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultExcludeAllForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForObjects(ExcludeAll).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": false,
        "booleanTrue": true,
        "float": 1234.56,
        "int": 123456,
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "string\"int"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": false,
    "booleanTrue": true,
    "float": 1234.56,
    "int": 123456,
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": null,
    "string": "string\"int"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultInheritForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForObjects(Inherit).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "ooo",
        "booleanTrue": "ooo",
        "float": "ooo",
        "int": "ooo",
        "nested": [
          {
            "prop1": "ooo",
            "prop2": "ooo"
          }
        ],
        "string": "ooo"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "ooo",
    "booleanTrue": "ooo",
    "float": "ooo",
    "int": "ooo",
    "nested": [
      {
        "prop1": "ooo",
        "prop2": "ooo"
      }
    ],
    "null": "ooo",
    "string": "ooo"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultInheritOverridableForObjects(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForObjects(InheritOverridable).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "ooo",
            "prop2": "ooo"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "ooo",
        "prop2": "ooo"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultExcludeForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForArrays(Exclude).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultExcludeAllForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForArrays(ExcludeAll).
		Build()

	expected := `{
  "array": [
    [
      "1",
      "2"
    ],
    {},
    {
      "int": 123456
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "1",
          "2"
        ],
        {},
        {
          "int": 123456
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultInheritForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForArrays(Inherit).
		Build()

	expected := `{
  "array": [
    [
      "aaa",
      "aaa"
    ],
    {},
    {
      "int": "aaa"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "aaa",
          "aaa"
        ],
        {},
        {
          "int": "aaa"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func TestJSONObfuscatorWithDefaultInheritOverridableForArrays(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		ForArrays(InheritOverridable).
		Build()

	expected := `{
  "array": [
    [
      "aaa",
      "aaa"
    ],
    {},
    {
      "int": "***"
    }
  ],
  "booleanFalse": "***",
  "booleanTrue": "***",
  "float": "***",
  "int": "***",
  "nested": [
    {
      "array": [
        [
          "aaa",
          "aaa"
        ],
        {},
        {
          "int": "***"
        }
      ],
      "booleanFalse": "***",
      "booleanTrue": "***",
      "float": "***",
      "int": "***",
      "nonMatchedNull": null,
      "notMatchedBooleanFalse": false,
      "notMatchedBooleanTrue": true,
      "notMatchedFloat": 1234.56,
      "notMatchedInt": 123456,
      "notMatchedString": "123456",
      "null": "***",
      "object": {
        "booleanFalse": "***",
        "booleanTrue": "***",
        "float": "***",
        "int": "***",
        "nested": [
          {
            "prop1": "1",
            "prop2": "2"
          }
        ],
        "string": "***"
      },
      "string": "***"
    }
  ],
  "nonMatchedNull": null,
  "nonMatchedObject": {
    "nonMatchedNull": null,
    "notMatchedBooleanFalse": false,
    "notMatchedBooleanTrue": true,
    "notMatchedFloat": 1234.56,
    "notMatchedInt": 123456,
    "notMatchedString": "123456"
  },
  "notMatchedBooleanFalse": false,
  "notMatchedBooleanTrue": true,
  "notMatchedFloat": 1234.56,
  "notMatchedInt": 123456,
  "notMatchedString": "123456",
  "null": "***",
  "object": {
    "booleanFalse": "***",
    "booleanTrue": "***",
    "float": "***",
    "int": "***",
    "nested": [
      {
        "prop1": "1",
        "prop2": "2"
      }
    ],
    "null": "***",
    "string": "***"
  },
  "string": "***"
}`

	testJSONObfuscation(t, jsonObfuscator, expected)
}

func testJSONObfuscation(t *testing.T, obfuscator Obfuscator, expectedObfuscated string) {
	t.Helper()

	testParseAndObfuscateString(t, obfuscator, inputJson, expectedObfuscated)
}

func TestObfuscateJSONStringOnErrorLog(t *testing.T) {
	logger := newCapturingLogger()
	builder := JSON().OnErrorLog(logger.Logger)

	input := inputJson + "x"

	var v any
	err := json.Unmarshal([]byte(input), &v)
	expectedLogged := fmt.Sprintf("ObfuscateString error: %v\n", err)

	testObfuscateJSONStringWithErrors(t, builder, logger, input, "", expectedLogged)
}

func TestObfuscateJSONStringOnErrorInclude(t *testing.T) {
	builder := JSON().OnErrorInclude()

	input := inputJson + "x"

	var v any
	err := json.Unmarshal([]byte(input), &v)

	testObfuscateJSONStringWithErrors(t, builder, nil, input, fmt.Sprintf("<error: %v>", err), "")
}

func TestObfuscateJSONStringOnErrorDiscard(t *testing.T) {
	builder := JSON().OnErrorDiscard()

	input := inputJson + "x"

	testObfuscateJSONStringWithErrors(t, builder, nil, input, "", "")
}

func testObfuscateJSONStringWithErrors(t *testing.T, builder *JSONObfuscatorBuilder, logger *CapturingLogger, input, expectedObfuscated, expectedLogged string) {
	t.Helper()

	obfuscator := WithFixedLength(3)
	jsonObfuscator := builder.
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		WithProperty("float", obfuscator, nil).
		WithProperty("booleanTrue", obfuscator, nil).
		WithProperty("booleanFalse", obfuscator, nil).
		WithProperty("object", WithFixedLengthWithMask(3, "o"), nil).
		WithProperty("array", WithFixedLengthWithMask(3, "a"), nil).
		WithProperty("null", obfuscator, nil).
		Build()

	testObfuscateStringWithErrors(t, jsonObfuscator, logger, input, expectedObfuscated, expectedLogged)
}

func TestJSONObfuscatorChaining(t *testing.T) {
	obfuscator := WithFixedLength(3)
	jsonObfuscator := JSON().
		WithProperty("string", obfuscator, nil).
		WithProperty("int", obfuscator, nil).
		Build()

	input := `{
  "string": "string\"int",
  "int": 123456,
  "float": 1234.56
}postfix`

	expected := `{
  "float": 1234.56,
  "int": "***",
  "string": "***"
}***`

	obfuscator = jsonObfuscator.UntilLength(len(input) - 7).Then(obfuscator)

	testParseAndObfuscateString(t, obfuscator, input, expected)
}
