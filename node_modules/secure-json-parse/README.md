# secure-json-parse

[![Build Status](https://dev.azure.com/fastify/fastify/_apis/build/status/fastify.secure-json-parse?branchName=master)](https://dev.azure.com/fastify/fastify/_build/latest?definitionId=8&branchName=master) [![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](http://standardjs.com/)

`JSON.parse()` drop-in replacement with prototype poisoning protection.

## Introduction

Consider this:

```js
> const a = '{"__proto__":{ "b":5}}';
'{"__proto__":{ "b":5}}'

> const b = JSON.parse(a);
{ __proto__: { b: 5 } }

> b.b;
undefined

> const c = Object.assign({}, b);
{}

> c.b
5
```

The problem is that `JSON.parse()` retains the `__proto__` property as a plain object key. By
itself, this is not a security issue. However, as soon as that object is assigned to another or
iterated on and values copied, the `__proto__` property leaks and becomes the object's prototype.

## Install
```
npm install secure-json-parse
```

## Usage

Pass the option object as a second (or third) parameter for configuring the action to take in case of a bad JSON, if nothing is configured, the default is to throw a `SyntaxError`.<br/>
You can choose which action to perform in case `__proto__` is present, and in case `constructor` is present.

```js
const sjson = require('secure-json-parse')

const goodJson = '{ "a": 5, "b": 6 }'
const badJson = '{ "a": 5, "b": 6, "__proto__": { "x": 7 }, "constructor": {"prototype": {"bar": "baz"} } }'

console.log(JSON.parse(goodJson), sjson.parse(goodJson, { protoAction: 'remove', constructorAction: 'remove' }))
console.log(JSON.parse(badJson), sjson.parse(badJson, { protoAction: 'remove', constructorAction: 'remove' }))
```

## API

### `sjson.parse(text, [reviver], [options])`

Parses a given JSON-formatted text into an object where:
- `text` - the JSON text string.
- `reviver` - the `JSON.parse()` optional `reviver` argument.
- `options` - optional configuration object where:
    - `protoAction` - optional string with one of:
        - `'error'` - throw a `SyntaxError` when a `__proto__` key is found. This is the default value.
        - `'remove'` - deletes any `__proto__` keys from the result object.
        - `'ignore'` - skips all validation (same as calling `JSON.parse()` directly).
    - `constructorAction` - optional string with one of:
        - `'error'` - throw a `SyntaxError` when a `constructor` key is found. This is the default value.
        - `'remove'` - deletes any `constructor` keys from the result object.
        - `'ignore'` - skips all validation (same as calling `JSON.parse()` directly).

### `sjson.scan(obj, [options])`

Scans a given object for prototype properties where:
- `obj` - the object being scanned.
- `options` - optional configuration object where:
    - `protoAction` - optional string with one of:
        - `'error'` - throw a `SyntaxError` when a `__proto__` key is found. This is the default value.
        - `'remove'` - deletes any `__proto__` keys from the input `obj`.
    - `constructorAction` - optional string with one of:
        - `'error'` - throw a `SyntaxError` when a `constructor` key is found. This is the default value.
        - `'remove'` - deletes any `constructor` keys from the input `obj`.

# Acknowledgements
This project has been forked from [hapijs/bourne](https://github.com/hapijs/bourne).
All the credits before the commit [4690682](https://github.com/hapijs/bourne/commit/4690682c6cdaa06590da7b2485d5df91c09da889) goes to the hapijs/bourne project contributors.
After, the project will be maintained by the Fastify team.

# License
Licensed under BSD-3-Clause.
