import { Kubernetes } from './kubernetes';

describe('Kubernetes', () => {
  it('should create an instance', () => {
    expect(new Kubernetes()).toBeTruthy();
  });
});
