import { UserTypePipe } from './user-type.pipe';

describe('UserTypePipe', () => {
  it('create an instance', () => {
    const pipe = new UserTypePipe();
    expect(pipe).toBeTruthy();
  });
});
