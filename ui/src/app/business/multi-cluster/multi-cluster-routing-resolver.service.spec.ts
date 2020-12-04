import { TestBed } from '@angular/core/testing';

import { MultiClusterRoutingResolverService } from './multi-cluster-routing-resolver.service';

describe('MultiClusterRoutingResolverService', () => {
  let service: MultiClusterRoutingResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(MultiClusterRoutingResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
