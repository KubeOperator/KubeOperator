'use strict';
const from = require('from2');
const pIsPromise = require('p-is-promise');

const intoStream = input => {
	if (Array.isArray(input)) {
		input = input.slice();
	}

	let promise;
	let iterator;

	prepare(input);

	function prepare(value) {
		input = value;

		if (
			input instanceof ArrayBuffer ||
			(ArrayBuffer.isView(input) && !Buffer.isBuffer(input))
		) {
			input = Buffer.from(input);
		}

		promise = pIsPromise(input) ? input : null;

		// We don't iterate on strings and buffers since slicing them is ~7x faster
		const shouldIterate = !promise && input[Symbol.iterator] && typeof input !== 'string' && !Buffer.isBuffer(input);
		iterator = shouldIterate ? input[Symbol.iterator]() : null;
	}

	return from(function reader(size, callback) {
		if (promise) {
			(async () => {
				try {
					await prepare(await promise);
					reader.call(this, size, callback);
				} catch (error) {
					callback(error);
				}
			})();

			return;
		}

		if (iterator) {
			const object = iterator.next();
			setImmediate(callback, null, object.done ? null : object.value);
			return;
		}

		if (input.length === 0) {
			setImmediate(callback, null, null);
			return;
		}

		const chunk = input.slice(0, size);
		input = input.slice(size);

		setImmediate(callback, null, chunk);
	});
};

module.exports = intoStream;
// TODO: Remove this for the next major release
module.exports.default = intoStream;

module.exports.object = input => {
	if (Array.isArray(input)) {
		input = input.slice();
	}

	let promise;
	let iterator;

	prepare(input);

	function prepare(value) {
		input = value;
		promise = pIsPromise(input) ? input : null;
		iterator = !promise && input[Symbol.iterator] ? input[Symbol.iterator]() : null;
	}

	return from.obj(function reader(size, callback) {
		if (promise) {
			(async () => {
				try {
					await prepare(await promise);
					reader.call(this, size, callback);
				} catch (error) {
					callback(error);
				}
			})();

			return;
		}

		if (iterator) {
			const object = iterator.next();
			setImmediate(callback, null, object.done ? null : object.value);
			return;
		}

		this.push(input);

		setImmediate(callback, null, null);
	});
};
