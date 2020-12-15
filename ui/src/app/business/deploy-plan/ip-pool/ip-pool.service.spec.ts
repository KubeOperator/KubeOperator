import { TestBed } from '@angular/core/testing';

import { IpPoolService } from './ip-pool.service';

describe('IpPoolService', () => {
  let service: IpPoolService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(IpPoolService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
