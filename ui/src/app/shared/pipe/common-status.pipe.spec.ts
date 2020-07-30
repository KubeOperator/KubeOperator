import { CommonStatusPipe } from './common-status.pipe';

describe('CommonStatusPipe', () => {
  it('create an instance', () => {
    const pipe = new CommonStatusPipe();
    expect(pipe).toBeTruthy();
  });
});
