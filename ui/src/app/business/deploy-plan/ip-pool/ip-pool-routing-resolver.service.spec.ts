import { TestBed } from '@angular/core/testing';

import { IpPoolRoutingResolverService } from './ip-pool-routing-resolver.service';

describe('IpPoolRoutingResolverService', () => {
  let service: IpPoolRoutingResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(IpPoolRoutingResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
