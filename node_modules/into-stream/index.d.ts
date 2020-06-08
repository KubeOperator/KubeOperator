/// <reference types="node"/>
import {Readable as ReadableStream} from 'stream';

declare namespace intoStream {
	type Input =
		| Buffer
		| NodeJS.TypedArray
		| ArrayBuffer
		| string
		| Iterable<Buffer | string>;

	type InputObject =
		| {[key: string]: any}
		| Iterable<{[key: string]: any}>;
}

declare const intoStream: {
	/**
	Convert `input` into a stream. Adheres to the requested chunk size, except for `array` where each element will be a chunk.

	@param input - The input to convert to a stream.
	@returns A [readable stream](https://nodejs.org/api/stream.html#stream_class_stream_readable).

	@example
	```
	import intoStream = require('into-stream');

	intoStream('unicorn').pipe(process.stdout);
	//=> 'unicorn'
	```
	*/
	(input: intoStream.Input | Promise<intoStream.Input>): ReadableStream;

	/**
	Convert object `input` into a stream.

	@param input - The object input to convert to a stream.
	@returns A [readable object stream](https://nodejs.org/api/stream.html#stream_object_mode).
	*/
	object(
		input: intoStream.InputObject | Promise<intoStream.InputObject>
	): ReadableStream;

	// TODO: Remove this for the next major release
	default: typeof intoStream;
};

export = intoStream;
