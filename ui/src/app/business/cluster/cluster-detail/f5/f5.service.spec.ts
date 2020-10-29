import { TestBed } from '@angular/core/testing';

import { F5Service } from './f5.service';

describe('F5Service', () => {
  let service: F5Service;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(F5Service);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
