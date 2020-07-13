import { Webkubectl } from './webkubectl';

describe('Webkubectl', () => {
  it('should create an instance', () => {
    expect(new Webkubectl()).toBeTruthy();
  });
});
