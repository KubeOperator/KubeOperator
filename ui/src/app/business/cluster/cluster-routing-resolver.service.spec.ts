import { TestBed } from '@angular/core/testing';

import { ClusterRoutingResolverService } from './cluster-routing-resolver.service';

describe('ClusterRoutingResolverService', () => {
  let service: ClusterRoutingResolverService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ClusterRoutingResolverService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
